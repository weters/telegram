package bot

// User represents a Telegram user or bot.
type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
}

// DisplayName will return the display name of the user or bot.
func (u *User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	}

	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}

	return u.FirstName
}
