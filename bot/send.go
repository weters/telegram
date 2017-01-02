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
	ChatID                int64        `json:"chat_id"`
	Text                  string       `json:"text"`
	ParseMode             string       `json:"parse_mode"`
	DisableWebPagePreview bool         `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool         `json:"disable_notification"`
	ReplyToMessageID      int64        `json:"reply_to_message_id,omitempty"`
	ReplyMarkup           *ReplyMarkup `json:"reply_markup,omitempty"`
}

// EditMessageText represents the payload that needs to be sent to Telegram's editMessageText method.
type EditMessageText struct {
	ChatID                int64        `json:"chat_id"`
	MessageID             int64        `json:"message_id"`
	Text                  string       `json:"text"`
	ParseMode             string       `json:"parse_mode"`
	DisableWebPagePreview bool         `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup           *ReplyMarkup `json:"reply_markup,omitempty"`
}

// SendDocument represents the payload that needs to be sent to Telegram's sendDocument method.
type SendDocument struct {
	ChatID           int64        `json:"chat_id"`
	Document         string       `json:"-"`
	ReplyToMessageID int64        `json:"reply_to_message_id,omitempty"`
	ReplyMarkup      *ReplyMarkup `json:"reply_markup"`
}
