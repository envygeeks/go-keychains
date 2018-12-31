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
	Label string
	Group string
	Key   string
}

func (t *TestService) Set(s string) error   { return nil }
func (t *TestService) Get() (string, error) { return "", nil }
func (t *TestService) Del() error           { return nil }
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
	s := New("tests.example.test", "tests.group.example.tests")
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
		a := s.ktl(ts.actual)
		assert.Equal(t, ts.expected, a,
			ts.description)
	}
}

func TestTok(t *testing.T) {
	s := New("tests.example.test", "tests.group.example.tests")
	type TestStruct struct {
		expected    string
		description string
		actual      string
	}

	for _, ts := range []TestStruct{
		TestStruct{
			actual:      "Hello World",
			description: "it should error if it has a space",
			expected:    "Hello World.tests.example.test",
		},
		TestStruct{
			description: "it should work",
			expected:    "hello-world.tests.example.test",
			actual:      "hello-world",
		},
	} {
		a := s.tok(ts.actual)
		assert.Equal(t, ts.expected, a)
	}
}
