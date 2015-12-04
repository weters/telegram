package bot

// SessionRecord represents an individual session.
type SessionRecord interface {
	// AuthorID should be value found in ur.Message.From.ID
	AuthorID() int
	// ChatID should be the chat ID found in ur.Message.Chat.ID
	ChatID() int
	// StateID should be an ID specified by the bot.
	StateID() int
	// Data is optional data that should be stored with the session.
	Data() string
}

// Session is an interface that has some capabilities for setting, deleting, and getting sessions.
type Session interface {
	// SetSession should set a session for a user in a chat.
	SetSession(authorID, chatID, stateID int, data string) error

	// DeleteSessionByAuthorIDAndChatID should delete a session for a user in a chat
	DeleteSessionByAuthorIDAndChatID(authorID, chatID int) error

	// SessionByAuthorIDAndChatID should return a session for a user. If there is no session, but otherwise there was no error,
	// (nil, nil) should be returned.
	SessionByAuthorIDAndChatID(authorID, chatID int) (SessionRecord, error)
}
