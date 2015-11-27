// bot contains functioality for interacting with Telegram's Bot API.
package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
)

// Telegram represents a Bot.
type Bot struct {
	BotName        string
	Token          string
	Handlers       map[string]Handler
	DefaultHandler Handler
}

// User represents a Telegram user or bot.
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// DisplayName will return the display name of the user or bot.
func (u *User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	}

	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}

	return u.FirstName
}

// Chat represents a Telegram chat.
type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type,omitempty"` // private group supergroup channel
	Title     string `json:"title,omitempty"`
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
}

// Message represents a Telegram message.
type Message struct {
	ID             int      `json:"message_id"`
	From           *User    `json:"from,omitempty"`
	Date           int      `json:"date"`
	Chat           *Chat    `json:"chat"`
	ForwardFrom    *User    `json:"forward_from,omitempty"`
	ForwardDate    int      `json:"forward_date,omitempty"`
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
	Text           string   `json:"text,omitempty"`
	// Audio
	// Document
	// Photo
	// Sticker
	// Video
	// Voice
	// Caption
	// Contact
	// Location
	NewChatParticipant  *User  `json:"new_chat_participant,omitempty"`
	LeftChatParticipant *User  `json:"left_chat_participant,omitempty"`
	NewChatTitle        string `json:"new_chat_title,omitempty"`
	// NewChatPhoto
	DeleteChatPhoto       bool `json:"delete_chat_photo,omitempty"`
	GroupChatCreated      bool `json:"group_chat_created,omitempty"`
	SupergroupChatCreated bool `json:"supergroup_chat_created,omitempty"`
	// ChannelChatCreated
	MigrateToChatID   int `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID int `json:"migrate_from_chat_id,omitempty"`
}

// UpdateResponse represents a response from a Telegram getUpdates method call.
type UpdateResponse struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

// IsGroup returns true if the chat type is "group"
func (p *UpdateResponse) IsGroup() bool {
	return p.Message.Chat.Type == "group"
}

// ReplyMarkup actually contains three Telegram objects in one: the ReplyKeyboardMarkup, ReplyKeyboardHide, and ForceReply objects.
type ReplyMarkup struct {
	// ReplyKeyboardMarkup
	Keyboard        [][]string `json:"keyboard,omitempty"`
	ResizeKeyboard  bool       `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool       `json:"one_time_keyboard,omitempty"`

	// ReplyKeyboardHide
	HideKeyboard bool `json:"hide_keyboard,omitempty"`

	// ForceReply
	ForceReply bool `json:"force_reply,omitempty"`

	// All
	Selective bool `json:"selective,omitempty"`
}

// SendMessage represents the payload that needs to be sent to Telegram's sendMessage method.
type SendMessage struct {
	ChatID                int          `json:"chat_id"`
	Text                  string       `json:"text"`
	DisableWebPagePreview bool         `json:"disable_web_page_preview,omitempty"`
	ReplyToMessageID      int          `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           *ReplyMarkup `json:"reply_markup"`
}

// Handler represents a function that can handle an update from Telegram.
type Handler func(t *Bot, ur *UpdateResponse, args string)

// New instantiates a new Telegram instance.
func New(botName, token string) *Bot {
	return &Bot{
		BotName:  botName,
		Token:    token,
		Handlers: make(map[string]Handler),
	}
}

// AddCommandHandler will register a Handler with a specific command.
//
// Example:
//   t.AddCommandHandler("help", HelpHandler)
//
// When a user types "/help" or "/help@YourBot", the HelpHandler will be called.
func (t *Bot) AddCommandHandler(c string, ch Handler) {
	t.Handlers[c] = ch
}

// SetDefaultHandler wil register a default handler to be called if a message was received
// and it wasn't a command.
func (t *Bot) SetDefaultHandler(dh Handler) {
	t.DefaultHandler = dh
}

var cmdRegex = regexp.MustCompile("^(?i)/([a-z0-9_]+)(?:@([a-z0-9_]+))?(?:\\s+(.*))?\\z")

// HandleUpdate will call an appropriate Handler depending on the UpdateResponse payload.
func (t *Bot) HandleUpdate(r *http.Request) error {
	d := json.NewDecoder(r.Body)
	var ur UpdateResponse
	if err := d.Decode(&ur); err != nil {
		return err
	}

	if match := cmdRegex.FindStringSubmatch(ur.Message.Text); match != nil {
		if match[2] != "" && match[2] != t.BotName {
			return nil
		}

		h, ok := t.Handlers[match[1]]
		if ok {
			h(t, &ur, match[3])
		}
	} else {
		if t.DefaultHandler != nil {
			t.DefaultHandler(t, &ur, "")
		}
	}

	return nil
}

// PostSendMessage will post a message to Telegram's sendMessage method.
func (t *Bot) PostSendMessage(msg *SendMessage) error {
	b := &bytes.Buffer{}
	j := json.NewEncoder(b)
	if err := j.Encode(msg); err != nil {
		return err
	}

	r, err := http.NewRequest("POST", t.URL("sendMessage"), b)
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
func (t *Bot) URL(m string) string {
	return "https://api.telegram.org/bot" + t.Token + "/" + m
}
