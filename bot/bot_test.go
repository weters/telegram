package bot

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type capture struct {
	helpHandlerCalled     bool
	sessionHandlerCalled  bool
	callbackHandlerCalled bool
	defaultHandlerCalled  bool
}

func (c *capture) helpHandler(b *Bot, u *UpdateResponse, args string) {
}

func (c *capture) sessionHandler(b *Bot, u *UpdateResponse, s SessionRecord) {
}

func (c *capture) callbackHandler(b *Bot, u *UpdateResponse) {
}

func (c *capture) defaultHandler(b *Bot, u *UpdateResponse, args string) {
}

func TestIncomingRequests(t *testing.T) {
	c := &capture{}
	s := newTestSession()

	b := New("Test_Bot", "mysecrettoken")
	assert.NotNil(t, b)

	b.AddCommandHandler("help", c.helpHandler)
	b.AddSessionHandler(100, c.sessionHandler)
	b.SetDefaultHandler(c.defaultHandler)
	b.SetBeforeCommandCallback(c.callbackHandler)
	b.SetSession(s)
}

type testSessionRecord struct {
	authorID int
	chatID   int
	stateID  int
	data     string
}

func (r *testSessionRecord) AuthorID() int {
	return r.authorID
}

func (r *testSessionRecord) ChatID() int {
	return r.chatID
}

func (r *testSessionRecord) StateID() int {
	return r.stateID
}

func (r *testSessionRecord) Data() string {
	return r.data
}

type testSession struct {
	data map[string]*testSessionRecord
}

func newTestSession() *testSession {
	return &testSession{
		data: make(map[string]*testSessionRecord),
	}
}

func (s *testSession) key(authorID, chatID int) string {
	return fmt.Sprintf("%d.%d", authorID, chatID)
}

// SetSession should set a session for a user in a chat.
func (s *testSession) SetSession(authorID, chatID, stateID int, data string) error {
	s.data[s.key(authorID, chatID)] = &testSessionRecord{authorID, chatID, stateID, data}
	return nil
}

// DeleteSessionByAuthorIDAndChatID should delete a session for a user in a chat
func (s *testSession) DeleteSessionByAuthorIDAndChatID(authorID, chatID int) error {
	delete(s.data, s.key(authorID, chatID))
	return nil
}

// SessionByAuthorIDAndChatID should return a session for a user. If there is no session, but otherwise there was no error,
// (nil, nil) should be returned.
func (s *testSession) SessionByAuthorIDAndChatID(authorID, chatID int) (SessionRecord, error) {
	r, ok := s.data[s.key(authorID, chatID)]
	if !ok {
		return nil, errors.New("could not retrieve session record")
	}

	return r, nil
}
