package keychain

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	wrappers          = map[string](func(string, string) Item){}
	AccessGroup       = "tweedy.group.io.envygeeks"
	ServiceNameSuffix = "tweedy.io.envygeeks"
	ServicePrefix     = "twitter-"
	supported         = []string{}
)

// Item is an interface to all of the
// possible keychain wrappers we can create
// in the future, such as CredentialStore,
// Keychain and on Linux: dbus junk.
type Item interface {
	Set(string) error
	Get() (string, error)
	Del() error
}

func Supported() bool {
	for _, s := range supported {
		if runtime.GOOS == s {
			return true
		}
	}

	return false
}

// LabelToKey will convert "This is the label"
// to "this-is-the-label" for use as a key
func LabelToKey(s string) string {
	s = strings.ToLower(s)
	r, _ := regexp.Compile(`\s+`)
	b := r.ReplaceAll([]byte(s), []byte("-"))
	return string(b)
}

// KeyToLabel will convert "this-is-the-label"
// to "This is the label" for label use
func KeyToLabel(s string) string {
	return strings.Title(strings.Replace(s, "-", " ", -1))
}

// New creates a new Item so you can `Get`,
// `Set`, or `Delete` said item. You can optionally
// skip this and just do `keychain.Get`, `Set`
// and `Delete`, this should normally only
// be used for persistent actions
func New(key string, label string) Item {
	if !Supported() {
		logrus.Fatalln("unsupported keychain OS")
	} else {
		if f, ok := wrappers[runtime.GOOS]; ok {
			i := f(key, label)
			return i
		}
	}

	return nil
}

// Get, Set, Del wraps around `New()` → `Get()`, `Set()`, `Del()`
func Set(k, s string) error        { return New(k, KeyToLabel(k)).Set(s) }
func Get(k string) (string, error) { return New(k, KeyToLabel(k)).Get() }
func Del(k string) error           { return New(k, KeyToLabel(k)).Del() }
