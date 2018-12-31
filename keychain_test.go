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

func TestKtl(t *testing.T) {
	type TestStruct struct {
		expected    string
		description string
		actual      string
	}

	for _, ts := range []TestStruct{
		TestStruct{
			actual:      "Hello World",
			description: "it works for Hello World",
			expected:    "Hello World",
		},
		TestStruct{
			actual:      "hello world",
			description: "it works for hello world",
			expected:    "Hello World",
		},
		TestStruct{
			actual:      "hello-world",
			description: "it works for hello-world",
			expected:    "Hello World",
		},
	} {
		a := ktl(ts.actual)
		assert.Equal(t, ts.expected, a,
			ts.description)
	}
}

func TestTol(t *testing.T) {
	type TestStruct struct {
		expected    string
		description string
		actual      string
		error       bool
	}

	for _, ts := range []TestStruct{
		TestStruct{
			actual:      "Hello World",
			description: "it should error if it has a space",
			expected:    "",
			error:       true,
		},
		TestStruct{
			description: "it should work",
			expected:    "hello-world.example.test",
			actual:      "hello-world",
		},
	} {
		a, err := tol(ts.actual, "example.test")
		assert.Equal(t, ts.expected, a)
		if ts.error {
			assert.Error(t, err, ts.description)
		}
	}
}
