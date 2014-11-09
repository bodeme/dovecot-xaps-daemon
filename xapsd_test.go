package main

import "testing"

func Test_parseCommand_empty(t *testing.T) {
	_, err := parseCommand("")
	if err != ErrMalformedCommand {
		t.Error("Did not get expected ErrMalformedCommand")
	}
}

func Test_parseCommand_missingParameters(t *testing.T) {
	_, err := parseCommand("HELLO")
	if err != ErrMalformedCommand {
		t.Error("Did not get expected ErrMalformedCommand")
	}
}

func Test_parseCommand(t *testing.T) {
	cmd, err := parseCommand("HELLO a=\"1\"\tb=\"2\"")
	if err != nil {
		t.Error("Failed to parse command")
	}
	if cmd.name != "HELLO" {
		t.Error("cmd.name != HELLO")
	}
	if len(cmd.args) != 2 {
		t.Error("len(cmd.args) != 2")
	}
	if cmd.args["a"] != "1" {
		t.Error("cmd.args[a] != 1")
	}
	if cmd.args["b"] != "2" {
		t.Error("cmd.args[b] != 2")
	}
}
