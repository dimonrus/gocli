package test

import (
	"errors"
	"fmt"
	"github.com/dimonrus/gocli"
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

	app := gocli.NewApplication(environment, rootPath+"/config/yaml", &config)
	app.ParseFlags(&config.Arguments)

	appType, ok := config.Arguments["app"]
	if ok != true {
		app.FatalError(errors.New("app type is not presents"))
	}

	p, _ := app.GetAbsolutePath("cool", "test")
	fmt.Println(p)

	value := appType.GetString()
	cos := make(chan bool)

	value = ApplicationTypeWeb

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

	go func() {
		err = app.Start("3333", func(command *gocli.Command) {
			v := command.Arguments()[0]
			app.SuccessMessage("Receive command: " + command.String())
			if v.Name == "exit" {
				app.AttentionMessage("Exit...", command)
				cos <- true
			} else {
				app.AttentionMessage(gohelp.AnsiRed+"Unknown command: "+command.String()+gohelp.AnsiReset, command)
			}
		})
	}()
	<-cos
	app.GetLogger(gocli.LogLevelInfo).Infoln("Server shutdown.")
}
