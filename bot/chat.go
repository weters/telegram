package bot

import (
	"encoding/json"
	"net/url"
	"strconv"
)

type Status string

const (
	StatusCreator       Status = "creator"
	StatusAdministrator Status = "administrator"
	StatusMember        Status = "member"
	StatusRestricted    Status = "restricted"
	StatusLeft          Status = "left"
	StatusKicked        Status = "kicked"
)

type ChatMember struct {
	User                  *User  `json:"user"`
	Status                Status `json:"status"`
	UntilDate             int    `json:"until_date,omit_empty"`
	CanBeEdited           bool   `json:"can_be_edited,omit_empty"`
	CanChangeInfo         bool   `json:"can_change_info,omit_empty"`
	CanPostMessages       bool   `json:"can_post_messages,omit_empty"`
	CanEditMessages       bool   `json:"can_edit_messages,omit_empty"`
	CanDeleteMessages     bool   `json:"can_delete_messages,omit_empty"`
	CanInviteUsers        bool   `json:"can_invite_users,omit_empty"`
	CanRestrictMembers    bool   `json:"can_restrict_members,omit_empty"`
	CanPinMessages        bool   `json:"can_pin_messages,omit_empty"`
	CanPromoteMembers     bool   `json:"can_promote_members,omit_empty"`
	CanSendMessages       bool   `json:"can_send_messages,omit_empty"`
	CanSendMediaMessages  bool   `json:"can_send_media_messages,omit_empty"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omit_empty"`
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omit_empty"`
}

func (c *ChatMember) IsValidStatus() bool {
	return c.Status != StatusKicked && c.Status != StatusLeft
}

func (b *Bot) GetChatMember(chatID, userID int64) (*ChatMemberResult, error) {
	v := url.Values{}
	v.Set("chat_id", strconv.FormatInt(chatID, 10))
	v.Set("user_id", strconv.FormatInt(userID, 10))
	url := b.URL("getChatMember") + "?" + v.Encode()

	resp, err := b.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var result ChatMemberResult
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
