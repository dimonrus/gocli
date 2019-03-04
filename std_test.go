package main

import (
	"errors"
	"os"
	"testing"
)

const (
	ApplicationTypeWeb      = "web"
	ApplicationTypeScript   = "script"
	ApplicationTypeConsumer = "consumer"
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
		panic("ENV is not defined")
	}

	app := DNApp{}.New(environment, &config)
	app.ParseFlags(&config.Arguments)

	appType, ok := config.Arguments["app"]
	if ok != true {
		app.FatalError(errors.New("app type is not presents"))
	}

	value := appType.Value.(*string)
	var err error

	switch *value {
	case ApplicationTypeWeb:
		err = app.Start(config.Arguments)
	default:
		err = errors.New("app type is undefined")
	}

	if err != nil {
		app.FatalError(err)
	}

	if !config.Project.Debug {
		app.FatalError(errors.New("debug mast be false"))
	}

	if config.Web.Port != 8000 {
		app.FatalError(errors.New("incorrect port"))
	}
}
