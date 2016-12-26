package bot

import "encoding/json"

// Result represents the result of a SendMessage result.
type MessageResult struct {
	Result      *Message `json:"result"`
	OK          bool     `json:"ok"`
	ErrorCode   int      `json:"error_code,omitempty"`
	Description string   `json:"description,omitempty"`
}

func (m *MessageResult) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(b)
}
