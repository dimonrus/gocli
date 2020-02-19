package test

import "github.com/dimonrus/gocli"

type Config struct {
	Project struct {
		Name  string
		Debug bool
	}
	Web struct {
		Port int
		Host string
	}
	Arguments gocli.Arguments
}
