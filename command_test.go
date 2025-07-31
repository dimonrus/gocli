package gocli

import (
	"fmt"
	"testing"
)

func TestParseCommand(t *testing.T) {
	command := ParseCommand([]byte("-app=script name=migration --class one"))
	if len(command.Arguments()) != 6 {
		t.Fatal("wrong command parsing")
	}
	if *command.arguments[0].Value.(*string) != "app" {
		t.Fatal("wrong parser")
	}
	if *command.arguments[1].Value.(*string) != "script" {
		t.Fatal("wrong parser")
	}
	if *command.arguments[2].Value.(*string) != "name" {
		t.Fatal("wrong parser")
	}
	if *command.arguments[3].Value.(*string) != "migration" {
		t.Fatal("wrong parser")
	}
	if *command.arguments[4].Value.(*string) != "class" {
		t.Fatal("wrong parser")
	}
	if *command.arguments[5].Value.(*string) != "one" {
		t.Fatal("wrong parser")
	}
	command = ParseCommand([]byte("-consumer stop=all\n"))
	if len(command.Arguments()) != 3 {
		t.Fatal("wrong command parsing")
	}
	command = ParseCommand([]byte("web -repeat=2\n always true"))
	if len(command.Arguments()) != 5 {
		t.Fatal("wrong command parsing")
	}
	if command.Arguments()[2].Type != ArgumentTypeUint {
		t.Fatal("must be int")
	}
	command = ParseCommand([]byte("web     repeat=2\n always 	 true\t false"))
	if len(command.Arguments()) != 6 {
		t.Fatal("wrong command parsing")
	}
	command = ParseCommand([]byte("docker restart some-goods-1"))
	if len(command.Arguments()) != 3 {
		t.Fatal("wrong command parsing some-goods-1")
	}
	fmt.Println(command.String())
}

func TestParseCommand2(t *testing.T) {
	command := ParseCommand([]byte("consumer set count 1 report"))
	args := command.Arguments()
	if len(args) != 5 {
		t.Fatal("wrong command parsing")
	}
	if args[1].GetString() != "set" {
		t.Fatal("wrong parser set")
	}
	if args[3].GetInt() != 1 {
		t.Fatal("wrong parser count number")
	}
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gocli
// cpu: Apple M2 Max
// BenchmarkParseCommand
// BenchmarkParseCommand-12    	 2019850	       597.0 ns/op	    1376 B/op	      14 allocs/op
func BenchmarkParseCommand(b *testing.B) {
	com := []byte("-app=script name=migration --class one")
	for i := 0; i < b.N; i++ {
		ParseCommand(com)
	}
	b.ReportAllocs()
}
