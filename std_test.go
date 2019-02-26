package main

import (
	"os"
	"testing"
)

type Config struct {
	Project struct {
		Name  string
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
	environment := os.Getenv("ENV")
	if environment == "" {
		panic("ENV is not defined")
	}
	app := DNApp{}.New(environment, &config)
	err := app.Start(arg)
	if err != nil {
		t.Fatal(err)
	}

	if !config.Project.Debug {
		t.Fatal("debug mast be false")
	}

	if config.Web.Port != 8000 {
		t.Fatal("incorrect port")
	}
}
