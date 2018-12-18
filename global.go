// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package keychain

import (
	"os/user"

	"github.com/sirupsen/logrus"
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
