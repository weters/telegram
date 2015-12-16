package bot

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

// SendDocument represents the payload that needs to be sent to Telegram's sendDocument method.
type SendDocument struct {
	ChatID           int          `json:"chat_id"`
	Document         string       `json:"-"`
	ReplyToMessageID int          `json:"reply_to_message_id,omitempty"`
	ReplyMarkup      *ReplyMarkup `json:"reply_markup"`
}
