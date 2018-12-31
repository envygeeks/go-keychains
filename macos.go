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
// where `Label` should be something like "This is a label",
// and `Key` should be something like "this-is-a-key.app.tld.domain",
// and `Group` should be like "app.group.tld.domain"
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
func (m *MacOS) Del() error {
	err := macos.DeleteItem(m.item)
	if err != nil {
		return err
	}

	return nil
}

// Set the item'
func (m *MacOS) Set(s string) error {
	var err error

	m.item.SetData([]byte(s))
	err = macos.AddItem(m.item)
	if err == macos.ErrorDuplicateItem {
		err = m.Del()
		if err == nil {
			err = macos.AddItem(m.item)
		}
	}

	return err
}

// Get the item
func (m *MacOS) Get() (string, error) {
	m.item.SetReturnData(true)
	r, err := macos.QueryItem(m.item)
	if err != nil {
		return "", nil
	}

	s := string(r[0].Data)
	return s, nil
}
