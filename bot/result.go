package bot

import "encoding/json"

type GenericResult struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

// Result represents the result of a SendMessage result.
type MessageResult struct {
	GenericResult
	Result *Message `json:"result"`
}

type ChatMemberResult struct {
	GenericResult
	Result *ChatMember `json:"result"`
}

func (m *MessageResult) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(b)
}
