package bot

// Defines the various chat types in Telegram
const (
	ChatTypePrivate    = "private"
	ChatTypeGroup      = "group"
	ChatTypeSupergroup = "supergroup"
	ChatTypeChannel    = "channel"
)

// UpdateResponse represents a response from a Telegram getUpdates method call.
type UpdateResponse struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
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

// IsGroup returns true if the chat type is "group"
func (ur *UpdateResponse) IsGroup() bool {
	return ur.Message.Chat.Type == ChatTypeGroup
}

// IsPrivate returns true if the chat type is "private"
func (ur *UpdateResponse) IsPrivate() bool {
	return ur.Message.Chat.Type == ChatTypePrivate
}

// ChatID is an accessor to p.Message.Chat.ID
func (ur *UpdateResponse) ChatID() int {
	return ur.Message.Chat.ID
}

// FromID is an accessor to p.Message.From.ID
func (ur *UpdateResponse) FromID() int {
	return ur.Message.From.ID
}

// IsBotReply will return true if the message received is a reply to a message from the bot.
func (ur *UpdateResponse) IsBotReply(b *Bot) bool {
	return ur.Message != nil && ur.Message.ReplyToMessage != nil && ur.Message.ReplyToMessage.From.Username == b.BotName
}
