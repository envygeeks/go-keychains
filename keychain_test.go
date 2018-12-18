package keychain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelToKey(t *testing.T) {
	type TestStruct struct {
		expected    string
		description string
		actual      string
	}

	for _, ts := range []TestStruct{
		TestStruct{
			actual:      "hello-world",
			description: "it works for hello-world",
			expected:    "hello-world",
		},
		TestStruct{
			actual:      "hello world",
			description: "it works for hello world",
			expected:    "hello-world",
		},
		TestStruct{
			actual:      "Hello World",
			description: "it works for Hello World",
			expected:    "hello-world",
		},
	} {
		a := LabelToKey(ts.actual)
		assert.Equal(t, ts.expected, a,
			ts.description)
	}
}

func TestKeyToLabel(t *testing.T) {
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
		a := KeyToLabel(ts.actual)
		assert.Equal(t, ts.expected, a,
			ts.description)
	}
}
