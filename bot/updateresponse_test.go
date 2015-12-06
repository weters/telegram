package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	u := &UpdateResponse{
		Message: &Message{
			Chat: &Chat{
				Type: ChatTypePrivate,
			},
		},
	}

	assert.False(t, u.IsGroup())
	assert.True(t, u.IsPrivate())

	u.Message.Chat.Type = ChatTypeGroup

	assert.True(t, u.IsGroup())
	assert.False(t, u.IsPrivate())
}

func TestChatID(t *testing.T) {
	u := &UpdateResponse{
		Message: &Message{
			Chat: &Chat{
				ID: 12345,
			},
		},
	}

	assert.Equal(t, 12345, u.ChatID())
}

func TestFromID(t *testing.T) {
	u := &UpdateResponse{
		Message: &Message{
			From: &User{
				ID: 23456,
			},
		},
	}

	assert.Equal(t, 23456, u.FromID())
}

func TestIsBotReply(t *testing.T) {
	b := &Bot{
		BotName: "Test_Bot",
	}

	ur := &UpdateResponse{}

	assert.False(t, ur.IsBotReply(b))

	ur = &UpdateResponse{
		Message: &Message{
			ReplyToMessage: &Message{
				From: &User{
					Username: "Test_Bot",
				},
			},
		},
	}

	assert.True(t, ur.IsBotReply(b))

	b.BotName = "Test2_Bot"

	assert.False(t, ur.IsBotReply(b))
}
