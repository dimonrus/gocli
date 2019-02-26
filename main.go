package main

import (
	"errors"
	"os"
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

func main() {
	var config Config
	environment := os.Getenv("ENV")
	if environment == "" {
		panic("ENV is not defined")
	}

	app := DNApp{}.New(environment, &config)
	app.ParseFlags(&config.Arguments)

	if appType, ok := config.Arguments["app"]; ok == true {
		value := appType.Value.(*string)
		switch *value {
		case ApplicationTypeWeb:
			err := app.Start(config.Arguments)
			if err != nil {
				app.FatalError(err)
			}
		default:
			app.FatalError(errors.New("wrong type"))
		}
	}

	if !config.Project.Debug {
		app.FatalError(errors.New("debug mast be false"))
	}

	if config.Web.Port != 8000 {
		app.FatalError(errors.New("incorrect port"))
	}
}
