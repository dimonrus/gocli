package main

const (
	ArgumentTypeString = "string"
	ArgumentTypeInt    = "int"
)

// Console app arguments
type Arguments map[string]Argument

type Argument struct {
	Type  string
	Value interface{}
	Label string
}
