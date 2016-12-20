// Package bot contains functioality for interacting with Telegram's Bot API.
package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// Bot represents a Telegram bot.
type Bot struct {
	BotName                string
	Token                  string
	CommandHandlers        map[string]Handler
	CommandPatternHandlers map[*regexp.Regexp]PatternHandler
	SessionHandlers        map[int]SessionHandler
	DefaultHandler         Handler
	BeforeCommandCallback  Callback
	Debug                  bool
	Session                Session

	botDirectMsgRegex *regexp.Regexp

	// allow us to inject a client for testing
	client *http.Client
}

// Handler represents a function that can handle an update from Telegram.
type Handler func(b *Bot, ur *UpdateResponse, args string)

// PatternHandler represents a function that can handle an update from Telegram.
type PatternHandler func(b *Bot, ur *UpdateResponse, matches []string)

// SessionHandler represents a function that can handle an update from Telegram with a session active.
type SessionHandler func(b *Bot, ur *UpdateResponse, s SessionRecord)

// Callback represents a function that can handle a callback.
type Callback func(b *Bot, ur *UpdateResponse)

// New instantiates a new Telegram instance.
func New(botName, token string) *Bot {
	return &Bot{
		BotName:                botName,
		Token:                  token,
		CommandHandlers:        make(map[string]Handler),
		CommandPatternHandlers: make(map[*regexp.Regexp]PatternHandler),
		SessionHandlers:        make(map[int]SessionHandler),
		botDirectMsgRegex:      regexp.MustCompile(fmt.Sprintf("^@%s\\s+", botName)),
		client:                 http.DefaultClient,
	}
}

// AddCommandHandler will register a Handler with a specific command.
//
// Example:
//   b.AddCommandHandler("help", HelpHandler)
//
// When a user types "/help" or "/help@YourBot", the HelpHandler will be called.
func (b *Bot) AddCommandHandler(c string, ch Handler) {
	b.CommandHandlers[c] = ch
}

// AddCommandPatternHandler will register a Handler with a specific pattern.
//
// Example:
//   b.AddCommandPatternHandler(regexp.MustCompile("delete\\d+"), DeleteHandler)
//
// When a user types "/help" or "/help@YourBot", the HelpHandler will be called.
func (b *Bot) AddCommandPatternHandler(r *regexp.Regexp, ph PatternHandler) {
	b.CommandPatternHandlers[r] = ph
}

// AddSessionHandler will register a SessionHandler for a given sID
func (b *Bot) AddSessionHandler(sID int, sh SessionHandler) {
	b.SessionHandlers[sID] = sh
}

// SetDefaultHandler wil register a default handler to be called if a message was received
// and it wasn't a command.
func (b *Bot) SetDefaultHandler(dh Handler) {
	b.DefaultHandler = dh
}

// SetBeforeCommandCallback will set a callback which is executed before a command is executed.
func (b *Bot) SetBeforeCommandCallback(cb Callback) {
	b.BeforeCommandCallback = cb
}

// SetSession sets the session object which is responsible for getting, setting, and deleting sessions.
func (b *Bot) SetSession(s Session) {
	b.Session = s
}

var cmdRegex = regexp.MustCompile("^(?i)/([a-z0-9_]+)(?:@([a-z0-9_]+))?(?:\\s+(.*))?\\z")

// HandleUpdate will call an appropriate Handler depending on the UpdateResponse payload.
// Attempts to find a command handler. If not found, attempts to find a session handler if there
// is an active session. Finally the default handler is called.
func (b *Bot) HandleUpdate(r *http.Request) error {
	d := json.NewDecoder(r.Body)
	var ur UpdateResponse
	if err := d.Decode(&ur); err != nil {
		return err
	}

	if ur.Message == nil {
		copy, _ := json.Marshal(ur)
		log.Printf("Error: null message found: %s\n", copy)
		return errors.New("null message found")
	}

	if b.Debug {
		copy, _ := json.Marshal(ur)
		log.Printf("received in %s: %s\n", b.BotName, copy)
	}

	if match := cmdRegex.FindStringSubmatch(ur.Message.Text); match != nil {
		// It's a command, but it's not intended for our bot
		if match[2] != "" && match[2] != b.BotName {
			return nil
		}

		if cb := b.BeforeCommandCallback; cb != nil {
			cb(b, &ur)
		}

		if h, ok := b.CommandHandlers[match[1]]; ok {
			h(b, &ur, match[3])
		} else {
			for r, h := range b.CommandPatternHandlers {
				if matches := r.FindStringSubmatch(match[1]); matches != nil {
					h(b, &ur, matches)
				}
			}
		}

		return nil
	}

	// if this was a direct message, strip out the bot name callout
	// "@My_Bot Hello" -> "Hello"
	ur.Message.Text = b.botDirectMsgRegex.ReplaceAllLiteralString(ur.Message.Text, "")

	if b.Session != nil {
		s, err := b.Session.SessionByAuthorIDAndChatID(ur.FromID(), ur.ChatID())
		if err != nil {
			return err
		}

		if s != nil {
			b.Session.DeleteSessionByAuthorIDAndChatID(ur.FromID(), ur.ChatID())

			if h, ok := b.SessionHandlers[s.StateID()]; ok {
				h(b, &ur, s)
				return nil
			}
		}
	}

	if b.DefaultHandler != nil {
		b.DefaultHandler(b, &ur, "")
	}

	return nil
}

// PostSendDocument will send a document and return the result from the server.
func (b *Bot) PostSendDocument(document *SendDocument) error {
	if document.Document == "" {
		return errors.New("bot: Document not specified")
	}

	file, err := os.Open(document.Document)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("document", filepath.Base(document.Document))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	writer.WriteField("chat_id", strconv.Itoa(document.ChatID))

	if document.ReplyToMessageID > 0 {
		writer.WriteField("reply_to_message_id", strconv.Itoa(document.ReplyToMessageID))
	}

	if document.ReplyMarkup != nil {
		b, err := json.Marshal(document.ReplyMarkup)
		if err != nil {
			return err
		}

		writer.WriteField("reply_markup", string(b))
	}

	if err := writer.Close(); err != nil {
		return err
	}

	resp, err := b.client.Post(b.URL("sendDocument"), writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (b *Bot) genericPost(endpoint string, msg interface{}) (*MessageResult, error) {
	bts := &bytes.Buffer{}
	j := json.NewEncoder(bts)
	if err := j.Encode(msg); err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", b.URL(endpoint), bts)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")
	resp, err := b.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var result MessageResult
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}
	return &result, err
}

// PostSendMessage will send a message and return the result from the server.
func (b *Bot) PostSendMessage(msg *SendMessage) (*MessageResult, error) {
	return b.genericPost("sendMessage", msg)
}

// PostEditMessageText will send a message and return the result from the server.
func (b *Bot) PostEditMessageText(msg *EditMessageText) (*MessageResult, error) {
	return b.genericPost("editMessageText", msg)
}

// SetWebhook will post a message to Telegram's setWebhook method. This will allow the bot
// to register with Telegram for notifications.
func (b *Bot) SetWebhook(uri string, certFile string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if certFile != "" {
		file, err := os.Open(certFile)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err := writer.CreateFormFile("certificate", filepath.Base(certFile))
		if err != nil {
			return err
		}

		if _, err := io.Copy(part, file); err != nil {
			return err
		}
	}

	writer.WriteField("url", uri)

	if err := writer.Close(); err != nil {
		return err
	}

	resp, err := b.client.Post(b.URL("setWebhook"), writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// URL returns the correct Telegram URL to use.
func (b *Bot) URL(m string) string {
	return "https://api.telegram.org/bot" + b.Token + "/" + m
}
