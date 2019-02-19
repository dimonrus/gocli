package main

import (
	"os"
	"testing"
)

type Config struct {
	Project struct {
		Name string
		Debug bool
	}
	Web struct {
		Port int
		Host string
	}
}

func TestName(t *testing.T) {
	var arg Arguments
	var config Config
	environment := os.Getenv("ENV");
	if environment == "" {
		panic("ENV is not defined")
	}
	app := DNApp{}.New(environment, &config)
	app.Start(arg)
}
