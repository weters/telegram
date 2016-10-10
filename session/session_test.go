package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemorySession(t *testing.T) {
	s := NewMemorySession()

	aSession, err := s.SessionByAuthorIDAndChatID(100, 200)
	assert.Nil(t, aSession)
	assert.NoError(t, err)

	err = s.SetSession(100, 200, 5, "my data")
	assert.NoError(t, err)

	aSession, err = s.SessionByAuthorIDAndChatID(100, 200)
	assert.NoError(t, err)
	assert.Equal(t, 100, aSession.AuthorID(), "correct AuthorID")
	assert.Equal(t, 200, aSession.ChatID(), "correct ChatID")
	assert.Equal(t, 5, aSession.StateID(), "correct StateID")
	assert.Equal(t, "my data", aSession.Data(), "correct data")

	err = s.SetSession(200, 100, 1, "ignore me")
	assert.NoError(t, err)

	aSession, err = s.SessionByAuthorIDAndChatID(100, 200)
	assert.Equal(t, 5, aSession.StateID(), "old record still returned")

	err = s.SetSession(100, 200, 10, "my data 2")
	assert.NoError(t, err)

	aSession, err = s.SessionByAuthorIDAndChatID(100, 200)
	assert.NoError(t, err)
	assert.Equal(t, 10, aSession.StateID(), "StateID updated")
	assert.Equal(t, "my data 2", aSession.Data(), "data updated")

	err = s.DeleteSessionByAuthorIDAndChatID(100, 200)
	assert.NoError(t, err)

	aSession, err = s.SessionByAuthorIDAndChatID(100, 200)
	assert.Nil(t, aSession, "record successfully deleted")
	assert.NoError(t, err)
}
