package bot

// Result represents the result of a SendMessage result.
type MessageResult struct {
	Result      *Message `json:"result"`
	OK          bool     `json:"ok"`
	ErrorCode   int      `json:"error_code"`
	Description string   `json:"description"`
}
