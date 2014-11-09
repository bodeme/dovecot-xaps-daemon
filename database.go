package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

/*
accounts: [
  { username: "stefan@arentz.ca",
    devices: [ { account_id:"a1", token:"t1", mailboxes:["INBOX"] },
               { account_id:"a2", token:"t2", mailboxes:["INBOX","Ham"] } ] }
  ]
]
*/

type Device struct {
	AccountId   string
	DeviceToken string
	Mailboxes   []string
}

type Account struct {
	Username string
	Devices  []Device
}

type database struct {
	path     string
	accounts map[string]*Account
}

func newDatabase() *database {
	return &database{
		path:     "",
		accounts: map[string]*Account{},
	}
}

func loadDatabase(path string) (*database, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return &database{
				path:     "",
				accounts: map[string]*Account{},
			}, nil
		} else {
			return nil, err
		}
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, err
	}

	db := &database{
		path:     "",
		accounts: map[string]*Account{},
	}

	for i := range accounts {
		db.accounts[accounts[i].Username] = &accounts[i]
	}

	return db, nil
}

func (db *database) addRegistration(username, accountId, deviceToken string, mailboxes []string) error {
	if account, ok := db.accounts[username]; ok {
		found := false
		for i, _ := range account.Devices {
			if account.Devices[i].AccountId == accountId && account.Devices[i].DeviceToken == deviceToken {
				account.Devices[i].Mailboxes = mailboxes
				found = true
			}
		}
		if !found {
			device := Device{
				AccountId:   accountId,
				DeviceToken: deviceToken,
				Mailboxes:   mailboxes,
			}
			account.Devices = append(account.Devices, device)
		}
	} else {
		db.accounts[username] = &Account{
			Username: username,
			Devices: []Device{
				Device{
					AccountId:   accountId,
					DeviceToken: deviceToken,
					Mailboxes:   mailboxes,
				},
			},
		}
	}
	return nil
}

func (db *database) findDevices(username, mailbox string) []Device {
	if account, ok := db.accounts[username]; ok {
		return account.Devices // TODO: Filter those that have the named mailbox
	}
	return nil
}
