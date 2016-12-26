package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageResultString(t *testing.T) {
	m := &MessageResult{OK: true}
	assert.Equal(t, `{"result":null,"ok":true}`, m.String())

	m = &MessageResult{ErrorCode: 100, Description: "it failed"}
	assert.Equal(t, `{"result":null,"ok":false,"error_code":100,"description":"it failed"}`, m.String())
}
