package main

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
	accounts map[string]*Account
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
