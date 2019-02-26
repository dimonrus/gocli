package main

const (
	ArgumentTypeString = "string"
	ArgumentTypeInt    = "int"

	ApplicationTypeWeb      = "web"
	ApplicationTypeScript   = "script"
	ApplicationTypeConsumer = "consumer"
)

// Console app arguments
type Arguments map[string]Argument

type Argument struct {
	Type  string
	Value interface{}
	Label string
}
