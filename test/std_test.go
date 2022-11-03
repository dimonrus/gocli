package test

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/porterr"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	ApplicationTypeWeb      = "web"
	ApplicationTypeScript   = "script"
	ApplicationTypeConsumer = "consumer"
)

func TestServerApp(t *testing.T) {
	var config Config
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "local"
	}

	rootPath, err := filepath.Abs("")
	if err != nil {
		panic(err)
	}
	_ = os.Setenv("WEB_PORT", "8000")
	_ = os.Setenv("WEB_HOST", "0.0.0.0")
	app := gocli.NewApplication(environment, rootPath+"/config/yaml", &config)
	app.ParseFlags(&config.Arguments)

	appType, ok := config.Arguments["app"]
	if ok != true {
		app.FatalError(errors.New("app type is not presents"))
	}

	p, _ := app.GetAbsolutePath("cool", "test")
	fmt.Println(p)

	value := appType.GetString()
	exit := make(chan struct{})

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
	if config.Web.Host != "0.0.0.0" {
		app.FatalError(errors.New("incorrect host"))
	}

	go func() {
		err = app.Start(":3333", func(command *gocli.Command) {
			v := command.Arguments()[0]
			app.SuccessMessage("Receive command: "+command.String(), command)
			if v.Name == "exit" {
				app.AttentionMessage("Exit...", command)
				exit <- struct{}{}
			} else if v.Name == "show" {
				app.AttentionMessage(gohelp.AnsiYellow+"The show is began"+gohelp.AnsiReset, command)
			} else {
				app.AttentionMessage(gohelp.AnsiRed+"Unknown command: "+command.String()+gohelp.AnsiReset, command)
			}
		})
	}()
	<-exit
	app.GetLogger().Infoln("Server shutdown.")
}

func Dial() (conn net.Conn, e porterr.IError) {
	var err error
	conn, err = net.Dial("tcp", "0.0.0.0:3333")
	if err != nil {
		e = porterr.New(porterr.PortErrorIO, err.Error())
	}
	return
}

func TestRunCommand(t *testing.T) {
	conn, e := Dial()
	if e != nil {
		return
	}
	defer func() {
		conn.Close()
	}()
	go func() {
		r := bufio.NewReader(conn)
		for {
			l, _, err := r.ReadLine()
			if err == io.EOF {
				t.Log("Конец чтения")
				break
			}
			if err != nil {
				t.Log(err.Error())
			}
			_ = l
			t.Log(string(l))
			time.Sleep(time.Millisecond * 300)
		}
	}()
	_, err := conn.Write([]byte("restart; puko;\n"))
	t.Log("Команда restart отправлена")
	if err != nil {
		e = porterr.New(porterr.PortErrorRequest, err.Error())
		return
	}
	_, err = conn.Write([]byte("show;\n"))
	t.Log("Команда show отправлена")
	if err != nil {
		e = porterr.New(porterr.PortErrorRequest, err.Error())
		return
	}
	_, err = conn.Write([]byte("exito;\n"))
	t.Log("Команда exit отправлена")
	if err != nil {
		e = porterr.New(porterr.PortErrorRequest, err.Error())
		return
	}
	time.Sleep(time.Second * 15)
	return
}
