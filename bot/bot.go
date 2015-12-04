// Package bot contains functioality for interacting with Telegram's Bot API.
package bot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

// Defines the various chat types in Telegram
const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeSupergroup = "supergroup"
	ChatTypeChannel    = "channel"
)

// Bot represents a Telegram bot.
type Bot struct {
	BotName               string
	Token                 string
	CommandHandlers       map[string]Handler
	DefaultHandler        Handler
	BeforeCommandCallback Callback
	Debug                 bool
}

// Handler represents a function that can handle an update from Telegram.
type Handler func(b *Bot, ur *UpdateResponse, args string)

// Callback represents a function that can handle a callback.
type Callback func(b *Bot, ur *UpdateResponse)

// New instantiates a new Telegram instance.
func New(botName, token string) *Bot {
	return &Bot{
		BotName:         botName,
		Token:           token,
		CommandHandlers: make(map[string]Handler),
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

// SetDefaultHandler wil register a default handler to be called if a message was received
// and it wasn't a command.
func (b *Bot) SetDefaultHandler(dh Handler) {
	b.DefaultHandler = dh
}

// SetBeforeCommandCallback will set a callback which is executed before a command is executed.
func (b *Bot) SetBeforeCommandCallback(cb Callback) {
	b.BeforeCommandCallback = cb
}

var cmdRegex = regexp.MustCompile("^(?i)/([a-z0-9_]+)(?:@([a-z0-9_]+))?(?:\\s+(.*))?\\z")

// HandleUpdate will call an appropriate Handler depending on the UpdateResponse payload.
func (b *Bot) HandleUpdate(r *http.Request) error {
	d := json.NewDecoder(r.Body)
	var ur UpdateResponse
	if err := d.Decode(&ur); err != nil {
		return err
	}

	if b.Debug {
		copy, _ := json.Marshal(ur)
		log.Printf("%s\n", copy)
	}

	if match := cmdRegex.FindStringSubmatch(ur.Message.Text); match != nil {
		// It's a command, but it's not intended for our bot
		if match[2] != "" && match[2] != b.BotName {
			return nil
		}

		if cb := b.BeforeCommandCallback; b != nil {
			cb(b, &ur)
		}

		if h, ok := b.CommandHandlers[match[1]]; ok {
			h(b, &ur, match[3])
		}
	} else if b.DefaultHandler != nil {
		b.DefaultHandler(b, &ur, "")
	}

	return nil
}

// PostSendMessage will post a message to Telegram's sendMessage method.
func (b *Bot) PostSendMessage(msg *SendMessage) error {
	b := &bytes.Buffer{}
	j := json.NewEncoder(b)
	if err := j.Encode(msg); err != nil {
		return err
	}

	r, err := http.NewRequest("POST", b.URL("sendMessage"), b)
	if err != nil {
		return err
	}

	r.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	resp, err := c.Do(r)
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
