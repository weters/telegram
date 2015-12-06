package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayName(t *testing.T) {
	u := &User{
		ID:        1000,
		FirstName: "John",
	}

	assert.Equal(t, "John", u.DisplayName())

	u.LastName = "Doe"

	assert.Equal(t, "John Doe", u.DisplayName())

	u.Username = "jdoe"

	assert.Equal(t, "jdoe", u.DisplayName())
}
