// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.
// +build !windows,!linux

package keychain

import (
	macos "github.com/keybase/go-keychain"
	"github.com/sirupsen/logrus"
)

type MacOS struct {
	item macos.Item

	User  string
	Label string
	Key   string
}

// NewForMacOS creates a new macOS item
func NewForMacOS(key, label, group string) Item {
	i := macos.NewItem()

	i.SetService(key)
	i.SetSynchronizable(macos.SynchronizableNo)
	i.SetAccessible(macos.AccessibleWhenUnlocked)
	i.SetSecClass(macos.SecClassGenericPassword)
	i.SetAccessGroup(group)
	i.SetAccount(User())
	i.SetLabel(label)

	m := "created new keychain item"
	logrus.WithFields(logrus.Fields{
		"label": label,
		"key":   key,
	}).Debug(m)

	return &MacOS{
		item:  i,
		Label: label,
		User:  User(),
		Key:   key,
	}
}

func init() {
	keychains["darwin"] = NewForMacOS
	supported = append(supported,
		"darwin")
}

// Del the item
func (t *MacOS) Del() error {
	err := macos.DeleteItem(t.item)
	if err != nil {
		return err
	}

	return nil
}

// Set the item'
func (t *MacOS) Set(s string) error {
	var err error

	t.item.SetData([]byte(s))
	err = macos.AddItem(t.item)
	if err == macos.ErrorDuplicateItem {
		err = t.Del()
		if err == nil {
			err = macos.AddItem(t.item)
		}
	}

	return err
}

// Get the item
func (t *MacOS) Get() (string, error) {
	t.item.SetReturnData(true)
	r, err := macos.QueryItem(t.item)
	if err != nil {
		return "", nil
	}

	s := string(r[0].Data)
	return s, nil
}
