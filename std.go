package gocli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/dimonrus/porterr"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"regexp"
	"strings"
	"testing"
)

const (
	CommandSessionHost = "localhost"
	CommandSessionType = "tcp"
)

// Dynamic Name Application
type DNApp struct {
	config     interface{}
	logger     Logger
	ConfigPath string
}

// Get config struct
func (a DNApp) GetConfig() interface{} {
	return a.config
}

// Get full path to config
func (a DNApp) GetConfigPath(env string) string {
	return fmt.Sprintf("%s/%s.yaml", a.ConfigPath, env)
}

// Set config struct
func (a DNApp) SetConfig(cfg interface{}) Application {
	a.config = cfg
	return &a
}

// Start application
func (a DNApp) Start(port string, callback func(command Command)) porterr.IError {
	// Listen localhost socket connection
	l, err := net.Listen(CommandSessionType, CommandSessionHost+":"+port)
	if err != nil {
		return porterr.NewF(porterr.PortErrorIO, "Listen socket error: %s", err.Error())
	}
	defer func() {
		err := l.Close()
		if err != nil {
			a.GetLogger(LogLevelErr).Errorln(err)
		}
	}()
	a.GetLogger(LogLevelInfo).Infof("Start listening %s commands on %s:%s", CommandSessionType, CommandSessionHost, port)
	var e porterr.IError
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			e = porterr.NewF(porterr.PortErrorIO, "Accept socket error: %s", err.Error())
			break
		}
		// Handle command
		go func(c net.Conn) {
			var buf bytes.Buffer
			_, err := io.Copy(&buf, c)
			if err != nil {
				a.GetLogger(LogLevelErr).Errorln(err)
				return
			}
			defer func() {
				if err := recover(); err != nil {
					a.GetLogger(LogLevelErr).Errorln("Command processor error:", err)
				}
				// Always close the connection after process command
				err = c.Close()
				if err != nil {
					a.GetLogger(LogLevelErr).Errorln(err)
					return
				}
			}()
			// Parse command and run
			command := ParseCommand(buf.Bytes())
			command.BindConnection(c)
			callback(command)
		}(conn)
	}
	return e
}

// Fatal error
func (a DNApp) FatalError(err error) {
	panic(err)
}

// Get logger
func (a DNApp) GetLogger(level int) Logger {
	if a.logger == nil || a.logger.GetLevel() != level {
		a.logger = NewLogger(level, "Application: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return a.logger
}

// Init app
func (a DNApp) New(env string, config interface{}) Application {
	data, err := ioutil.ReadFile(a.GetConfigPath(env))
	if err != nil {
		a.FatalError(err)
	}
	// check if config has depends on other configs
	r, _ := regexp.Compile(`depends:(.*)`)
	matches := r.FindStringSubmatch(string(data))
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		// load parent config
		a.New(strings.TrimSpace(matches[1]), config)
	}
	// unmarshal config file in config struct
	err = yaml.Unmarshal(data, config)
	if err != nil {
		a.FatalError(err)
	}
	a.config = config

	return &a
}

// Parse console arguments
func (a DNApp) ParseFlags(args *Arguments) {
	for key, argument := range *args {
		argument.Name = key
		switch argument.Type {
		case ArgumentTypeString:
			var value string
			argument.Value = &value
			(*args)[key] = argument
			flag.StringVar(&value, key, "", argument.Label)
		case ArgumentTypeInt:
			var value int64
			argument.Value = &value
			(*args)[key] = argument
			flag.Int64Var(&value, key, 0, argument.Label)
		case ArgumentTypeUint:
			var value uint64
			argument.Value = &value
			(*args)[key] = argument
			flag.Uint64Var(&value, key, 0, argument.Label)
		case ArgumentTypeBool:
			var value bool
			argument.Value = &value
			(*args)[key] = argument
			flag.BoolVar(&value, key, false, argument.Label)
		default:
			a.FatalError(errors.New("argument type: " + argument.Type + " is not supported. Argument: " + argument.Label))
		}
	}
	testing.Init()
	flag.Parse()
}
