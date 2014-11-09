package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var ErrUnknownApsSubtopic = errors.New("Unknown aps-subtopic received")
var ErrUnknownCommand = errors.New("Unknown command received")
var ErrMalformedCommand = errors.New("Malformed command")

type command struct {
	name string
	args map[string]string
}

func unescapeParameterValue(v string) string {
	if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
		return v[1 : len(v)-1]
	} else {
		return v
	}
}

func parseCommand(s string) (command, error) {
	cmd := command{args: map[string]string{}}

	parts := strings.SplitN(s, " ", 2)
	if len(parts) != 2 {
		return command{}, ErrMalformedCommand
	}
	cmd.name = parts[0]

	for _, t := range strings.Split(parts[1], "\t") {
		pair := strings.SplitN(t, "=", 2)
		if len(pair) != 2 {
			return command{}, ErrMalformedCommand
		}
		cmd.args[pair[0]] = unescapeParameterValue(pair[1])
	}

	return cmd, nil
}

type app struct {
	db *database
}

func (a *app) handleRegister(cmd command) error {
	if cmd.args["aps-subtopic"] != "com.apple.mobilemail" {
		return ErrUnknownApsSubtopic
	}
	return a.db.addRegistration(cmd.args["dovecot-username"], cmd.args["aps-account-id"],
		cmd.args["aps-device-token"], []string{"INBOX"}) // cmd.args["dovecot-mailboxes"])
}

func (a *app) handleNotify(cmd command) error {
	for _, device := range a.db.findDevices(cmd.args["dovecot-username"], cmd.args["dovecot-mailbox"]) {
		log.Printf("Sending notification for %s/%s to %s", cmd.args["dovecot-username"],
			cmd.args["dovecot-mailbox"], device.AccountId)
	}
	return nil
}

func (a *app) handleUnknownCommand(cmd command) error {
	return ErrUnknownCommand
}

func (a *app) dispatchCommand(cmd command) (string, error) {
	return "", ErrUnknownCommand
}

func (a *app) handleConnection(c net.Conn) {
	log.Print("Processing connection: ", c)

	reader := bufio.NewReader(c)
	writer := bufio.NewWriter(c)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Print("Can't read line: ", err)
			}
			break
		}
		log.Print("Incoming: ", line)

		cmd, err := parseCommand(line)
		if err != nil {
			writer.WriteString("ERROR " + err.Error() + "\n")
			writer.Flush()
		} else {
			result, err := a.dispatchCommand(cmd)
			if err != nil {
				writer.WriteString("ERROR " + err.Error() + "\n")
				writer.Flush()
			} else {
				writer.WriteString("OK " + result)
				writer.Flush()
			}
		}
	}

	log.Print("Done with connection: ", c)
}

func main() {

	a := app{
		db: &database{},
	}

	l, err := net.Listen("unix", "/tmp/xapsd.sock")
	if err != nil {
		log.Fatal("Cannot listen: ", err)
		return
	}

	defer os.Remove("/tmp/xapsd.sock")

	log.Print("Listening on /tmp/xapsd.sock")

	for {
		c, err := l.Accept()
		if err != nil {
			log.Print("Accept error", err)
			return
		}
		go a.handleConnection(c)
	}
}
