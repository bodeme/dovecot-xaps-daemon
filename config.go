// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

type config struct {
	socket   string
	database string
}

func defaultConfig() config {
	return config{
		socket:   "/tmp/xapsd.sock",
		database: "/tmp/xapsd.json",
	}
}
