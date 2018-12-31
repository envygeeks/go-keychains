package keychain

import (
	"errors"
	"fmt"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Func Provides an interface to service functions
// where (string, string, string) should be the available
// key, label, and group to be passed upstream
type Func func(string, string, string) Item

/**
 */
var (
	keychains = map[string]Func{}
	supported = []string{}
)

// User gets the current user
func User() string {
	u, err := user.Current()
	if err != nil {
		logrus.Fatalln(err)
	}

	s := u.Username
	return s
}

// Supported : valid service?
func Supported() bool {
	for _, s := range supported {
		if runtime.GOOS == s {
			return true
		}
	}

	return false
}

// Item is an interface to a wrapper
// Exp MacOSKeychainItem, WindowsCredentialItem
// LinuxKeychainItem, KWalletItem
type Item interface {
	Set(string) error
	Get() (string, error)
	Del() error
}

var (
	ErrUnsupportedOS = errors.New("unsupported OS, no keychain")
	ErrNoKeychainFnd = errors.New("no keychain found")
)

// NewItem creates a new Item
func (s *Service) NewItem(k string) (Item, error) {
	var (
		err error
		l   string
	)

	if !Supported() {
		err = ErrUnsupportedOS
		goto fail
	} else {
		if f, ok := keychains[runtime.GOOS]; ok {
			l, err = tol(k, s.domain)
			if err == nil {
				return f(k, l, s.group), nil
			}
		}
	}
fail:
	return nil, err
}

// Service is a service wrapper
// `domain` example: `app.tld.domain`
// `group`: app.group.tld.domain, or otherwise
//   ↳ can also be a proper system group
type Service struct {
	domain, group string
}

// New is a new Keychain
func New(domain, group string) *Service {
	return &Service{
		domain, group,
	}
}

// ktl will convert "this-is-the-label" to
// "This is the label" for label use
func ktl(s string) string {
	return strings.Title(strings.Replace(s, "-", " ", -1))
}

// tol converts a key to a label
func tol(k, domain string) (string, error) {
	if strings.Contains(k, " ") {
		err := errors.New("key cannot be a label")
		return "", err
	}

	l := fmt.Sprintf("%s.%s", k, domain)
	return l, nil
}

// Set → `NewItem()` → `Set()`
func (s *Service) Set(k, v string) error {
	i, err := s.NewItem(k)
	if err != nil {
		return err
	}

	return i.Set(v)
}

// Get → `NewItem()`  → `Get()`
func (s *Service) Get(k string) (string, error) {
	i, err := s.NewItem(k)
	if err != nil {
		return "", err
	}

	return i.Get()
}

// String → `Get()`
func (s *Service) String(k string) (string, error) {
	return s.Get(k)
}

// Int → `Get()`
func (s *Service) Int(k string) (int, error) {
	var i int

	ss, err := s.String(k)
	if err != nil {
		goto fail
	}

	i, err = strconv.Atoi(ss)
	if err == nil {
		return i, nil
	}
fail:
	return 0, err
}

// Bool → `Get()`
func (s *Service) Bool(k string) (bool, error) {
	var b bool

	ss, err := s.String(k)
	if err != nil {
		goto fail
	}

	b, err = strconv.ParseBool(ss)
	if err == nil {
		return b, nil
	}
fail:
	return false, err
}

// Del → `NewItem()` → `Del()`
func (s *Service) Del(k string) error {
	i, err := s.NewItem(k)
	if err != nil {
		return err
	}

	return i.Del()
}
