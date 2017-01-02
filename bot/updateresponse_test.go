package bot

import (
	"encoding/json"
	"log"
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

	u.Message.Chat.Type = ChatTypeSupergroup

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

	assert.Equal(t, int64(12345), u.ChatID())
}

func TestFromID(t *testing.T) {
	u := &UpdateResponse{
		Message: &Message{
			From: &User{
				ID: 23456,
			},
		},
	}

	assert.Equal(t, int64(23456), u.FromID())
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

func TestDecode(t *testing.T) {
	jsonStr := `
{
    "update_id": 897497725,
    "message": {
        "message_id": 4104,
        "from": {
            "id": 144365044,
            "first_name": "John",
            "last_name": "Doe",
			"username" : "jdoe"
        },
        "date": 1449369855,
        "chat": {
            "id": 144255044,
            "type": "private",
            "first_name": "John",
            "last_name": "Doe"
        },
        "reply_to_message": {
            "message_id": 3165,
            "from": {
                "id": 141317493,
                "first_name": "Test Bot",
                "username": "Test_Bot"
            },
            "date": 1449210338,
            "chat": {
                "id": 144354044,
                "type": "private",
                "first_name": "John",
                "last_name": "Doe"
            },
            "text": "What do you want to know?"
        },
        "text": "Will this work?"
    }
}`

	var ur UpdateResponse
	if err := json.Unmarshal([]byte(jsonStr), &ur); err != nil {
		log.Fatal(err)
	}

	assert.NotNil(t, ur)
	assert.Equal(t, int64(897497725), ur.UpdateID, "correct UpdateID")

	assert.NotNil(t, ur.Message)
	assert.Equal(t, int64(4104), ur.Message.ID, "correct Message.ID")
	assert.Equal(t, int64(144365044), ur.Message.From.ID, "correct Message.From.ID")
	assert.Equal(t, "John", ur.Message.From.FirstName, "correct Message.From.FirstName")
	assert.Equal(t, "Doe", ur.Message.From.LastName, "correct Message.From.LastName")
	assert.Equal(t, "jdoe", ur.Message.From.Username, "correct Message.From.Username")

	assert.Equal(t, 1449369855, ur.Message.Date, "correct Message.Date")

	assert.Equal(t, int64(144255044), ur.Message.Chat.ID, "correct Message.Chat.ID")
	assert.Equal(t, "private", ur.Message.Chat.Type, "correct Message.Chat.Type")
	assert.Equal(t, "John", ur.Message.Chat.FirstName, "correct Message.Chat.FirstName")
	assert.Equal(t, "Doe", ur.Message.Chat.LastName, "correct Message.Chat.LastName")

	assert.NotNil(t, ur.Message.ReplyToMessage)
	assert.Equal(t, int64(3165), ur.Message.ReplyToMessage.ID, "correct Message.ReplyToMessage.ID")
	assert.Equal(t, int64(141317493), ur.Message.ReplyToMessage.From.ID, "correct Message.ReplyToMessage.From.ID")
	assert.Equal(t, "Test Bot", ur.Message.ReplyToMessage.From.FirstName, "correct Message.ReplyToMessage.From.FirstName")
	assert.Equal(t, "Test_Bot", ur.Message.ReplyToMessage.From.Username, "correct Message.ReplyToMessage.From.Username")
	assert.Equal(t, 1449210338, ur.Message.ReplyToMessage.Date, "correct Message.ReplyToMessage.Date")
	assert.Equal(t, int64(144354044), ur.Message.ReplyToMessage.Chat.ID, "correct Message.ReplyToMessage.Chat.ID")
	assert.Equal(t, "private", ur.Message.ReplyToMessage.Chat.Type, "correct Message.ReplyToMessage.Chat.Type")
	assert.Equal(t, "John", ur.Message.ReplyToMessage.Chat.FirstName, "correct Message.ReplyToMessage.Chat.FirstName")
	assert.Equal(t, "Doe", ur.Message.ReplyToMessage.Chat.LastName, "correct Message.ReplyToMessage.Chat.LastName")
	assert.Equal(t, "What do you want to know?", ur.Message.ReplyToMessage.Text, "correct Message.ReplyToMessage.Text")

	assert.Equal(t, "Will this work?", ur.Message.Text, "correct Message.Text")

}
