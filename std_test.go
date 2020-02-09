package gocli

import (
	"errors"
	"github.com/dimonrus/gohelp"
	"os"
	"path/filepath"
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

	rootPath, err := filepath.Abs("")
	if err != nil {
		panic(err)
	}

	app := DNApp{ConfigPath: rootPath + "/config/yaml"}.New(environment, &config)
	app.ParseFlags(&config.Arguments)

	appType, ok := config.Arguments["app"]
	if ok != true {
		app.FatalError(errors.New("app type is not presents"))
	}

	value := appType.GetString()
	cos := make(chan bool)

	value = ApplicationTypeWeb

	go func() {
		err = app.Start("3333", func(command Command) {
			v := command.Arguments()[0]
			app.GetLogger(LogLevelInfo).Infof("Receive command: %s", command.GetOrigin())
			if v.Name == "exit" {
				cos <- true
			} else {
				app.GetLogger(LogLevelWarn).Warnf(gohelp.AnsiRed+"Unknown command: %s"+gohelp.AnsiReset, command.GetOrigin())
				e := command.Result([]byte("Unknown command\n"))
				if e != nil {
					app.GetLogger(LogLevelInfo).Errorln(e)
				}
			}
		})
	}()

	switch value {
	case ApplicationTypeWeb:
		//start web
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

	<-cos
	app.GetLogger(LogLevelInfo).Warnln("Exit...")
}
