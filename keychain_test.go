package keychain

import (
	"os/user"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 */
type TestService struct {
	Key   string
	Label string
	Group string
	Lc    string
	Lv    string
}

func (t *TestService) Get() (string, error) { t.Lc = "get"; return "", nil }
func (t *TestService) Set(s string) error   { t.Lc, t.Lv = "set", s; return nil }
func (t *TestService) Del() error           { t.Lc = "del"; return nil }
func NewForTest(key, label, group string) Item {
	return &TestService{
		Group: group,
		Label: label,
		Key:   key,
	}
}

func init() {
	keychains[runtime.GOOS] = NewForTest
	supported = append(supported,
		runtime.GOOS)
}

// --

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
