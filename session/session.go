// Package session provides thread-safe in-memory session management for github.com/weters/telegram/bot.
package session

import (
	"fmt"
	"sync"

	"github.com/weters/telegram/bot"
)

// Record represents an individual session.
type Record struct {
	authorID int64
	chatID   int64
	stateID  int
	data     string
}

// AuthorID returns the value found in u.Message.From.ID
func (s *Record) AuthorID() int64 {
	return s.authorID
}

// ChatID returns the chat ID found in u.Message.Chat.ID.
func (s *Record) ChatID() int64 {
	return s.chatID
}

// StateID returns the state ID.
func (s *Record) StateID() int {
	return s.stateID
}

// Data returns the data value.
func (s *Record) Data() string {
	return s.data
}

// MemorySession provides capabilities for setting, getting, and deleting sessions in-memory.
type MemorySession struct {
	sessions map[string]*Record
	mutex    sync.RWMutex
}

// NewMemorySession returns a new MemorySession object.
func NewMemorySession() *MemorySession {
	return &MemorySession{
		sessions: make(map[string]*Record),
		mutex:    sync.RWMutex{},
	}
}

// SetSession sets a session for a user in a chat.
func (m *MemorySession) SetSession(authorID, chatID int64, stateID int, data string) error {
	s := &Record{
		authorID: authorID,
		chatID:   chatID,
		stateID:  stateID,
		data:     data,
	}

	key := m.key(authorID, chatID)

	m.mutex.Lock()
	m.sessions[key] = s
	m.mutex.Unlock()

	return nil
}

// DeleteSessionByAuthorIDAndChatID deletes a session for a user in a chat
func (m *MemorySession) DeleteSessionByAuthorIDAndChatID(authorID, chatID int64) error {
	key := m.key(authorID, chatID)

	m.mutex.Lock()
	delete(m.sessions, key)
	m.mutex.Unlock()

	return nil
}

// SessionByAuthorIDAndChatID returns a session for a user. If there is no session, but otherwise there was no error,
// (nil, nil) will be returned.
func (m *MemorySession) SessionByAuthorIDAndChatID(authorID, chatID int64) (bot.SessionRecord, error) {
	key := m.key(authorID, chatID)

	m.mutex.RLock()
	s := m.sessions[key]
	m.mutex.RUnlock()

	return s, nil
}

func (m *MemorySession) key(authorID, chatID int64) string {
	return fmt.Sprintf("%d:%d", authorID, chatID)
}
