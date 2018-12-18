package keychain

import (
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	u, _ := user.Current()
	type TestStruct struct {
		description string
		expected    string
	}

	for _, ts := range []TestStruct{
		TestStruct{
			description: "it gives the user",
			expected:    u.Username,
		},
	} {
		a := User()
		assert.Equal(t, ts.expected, a,
			ts.description)
	}
}
