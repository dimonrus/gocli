package gocli

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	command := ParseCommand([]byte("-app=script name=migration --class one"))
	if len(command.Arguments()) != 6 {
		t.Fatal("wrong command parsing")
	}
	command = ParseCommand([]byte("-consumer stop=all\n"))
	if len(command.Arguments()) != 3 {
		t.Fatal("wrong command parsing")
	}
	command = ParseCommand([]byte("web -repeat=2\n always true"))
	if len(command.Arguments()) != 5 {
		t.Fatal("wrong command parsing")
	}
	if command.Arguments()[2].Type != ArgumentTypeInt {
		t.Fatal("must be int")
	}
}
