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
	Arguments Arguments
}

func TestName(t *testing.T) {
	var config Config
	environment := os.Getenv("ENV")
	if environment == "" {
		t.Fatal("ENV is not defined")
	}

	app := DNApp{}.New(environment, &config)
	app.ParseFlags(&config.Arguments)

	if appType, ok := config.Arguments["app"]; ok == true {
		switch appType.Value {
		case ApplicationTypeWeb:
			err := app.Start(config.Arguments)
			if err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatal("wrong type")
		}
	}

	if !config.Project.Debug {
		t.Fatal("debug mast be false")
	}

	if config.Web.Port != 8000 {
		t.Fatal("incorrect port")
	}
}
