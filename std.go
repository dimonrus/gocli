package gocli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/porterr"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

const (
	CommandSessionHost = "localhost"
	CommandSessionPort = "8080"
	CommandSessionType = "tcp"
)

// Dynamic Name Application
// Implements Application interface
type DNApp struct {
	// Config
	config config
	// Logger
	logger Logger
}

// Application configuration
type config struct {
	// Values of parsed configs
	values interface{}
	// Path of config
	path string
}

// Create new Application
func NewApplication(env string, configPath string, values interface{}) Application {
	app := &DNApp{
		config: config{
			values: values,
			path:   configPath,
		},
	}
	return app.ParseConfig(env)
}

// Get config struct
func (a DNApp) GetConfig() interface{} {
	return a.config.values
}

// Set config struct
func (a *DNApp) SetConfig(cfg interface{}) Application {
	a.config.values = cfg
	return a
}

// Get full path to config
func (a DNApp) GetConfigPath(env string) string {
	return fmt.Sprintf("%s/%s.yaml", a.config.path, env)
}

// Get absolute path to application
func (a DNApp) GetAbsolutePath(path string, dir string) (string, porterr.IError) {
	rootPath, err := filepath.Abs("")
	if err != nil {
		return "", porterr.New(porterr.PortErrorArgument, "root path is incorrect")
	}
	if rootPath[len(rootPath)-1:] != string(os.PathSeparator) {
		rootPath = rootPath + string(os.PathSeparator)
	}
	return gohelp.BeforeString(rootPath, dir) + dir + string(os.PathSeparator) + path, nil
}

// Fatal error
func (a DNApp) FatalError(err error) {
	panic(err)
}

// Get logger
func (a *DNApp) GetLogger(level int) Logger {
	if a.logger == nil || a.logger.GetLevel() != level {
		a.logger = NewLogger(level, "Application: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	return a.logger
}

// Success message
func (a DNApp) SuccessMessage(message string, command ...*Command) {
	message = gohelp.AnsiGreen + message + gohelp.AnsiReset
	a.GetLogger(LogLevelInfo).Infoln(message)
	for _, c := range command {
		e := c.Result([]byte(message + "\n"))
		if e != nil {
			a.GetLogger(LogLevelWarn).Errorln(e)
		}
	}
	return
}

// Attention message
func (a DNApp) AttentionMessage(message string, command ...*Command) {
	message = gohelp.AnsiCyan + message + gohelp.AnsiReset
	a.GetLogger(LogLevelWarn).Warnln(message)
	for _, c := range command {
		e := c.Result([]byte(message + "\n"))
		if e != nil {
			a.GetLogger(LogLevelWarn).Errorln(e)
		}
	}
	return
}

// Fail message
func (a DNApp) FailMessage(message string, command ...*Command) {
	message = gohelp.AnsiRed + message + gohelp.AnsiReset
	a.GetLogger(LogLevelErr).Errorln(message)
	for _, c := range command {
		e := c.Result([]byte(message + "\n"))
		if e != nil {
			a.GetLogger(LogLevelWarn).Errorln(e)
		}
	}
	return
}

// Config parser
func (a *DNApp) ParseConfig(env string) Application {
	data, err := ioutil.ReadFile(a.GetConfigPath(env))
	if err != nil {
		a.FatalError(err)
	}
	// check if config has depends on other configs
	r, _ := regexp.Compile(`depends:(.*)`)
	matches := r.FindStringSubmatch(string(data))
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		// load parent config
		a.ParseConfig(strings.TrimSpace(matches[1]))
	}
	// unmarshal config file in config struct
	err = yaml.Unmarshal(data, a.config.values)
	if err != nil {
		a.FatalError(err)
	}
	return a
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

// Start application
func (a DNApp) Start(address string, callback func(command *Command)) porterr.IError {
	if address == "" {
		return porterr.NewF(porterr.PortErrorArgument, "address is required")
	}
	if callback == nil {
		return porterr.NewF(porterr.PortErrorArgument, "callback is required")
	}
	// host and port
	var host, port = CommandSessionHost, CommandSessionPort
	addressParts := strings.Split(address, ":")
	if len(addressParts) == 2 {
		if len(addressParts[0]) > 0 {
			host = addressParts[0]
		}
		if len(addressParts[1]) > 0 {
			port = addressParts[1]
		}
	}
	l, err := net.Listen(CommandSessionType, host+":"+port)
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
			r := bufio.NewReader(c)
			for {
				com, _, err := r.ReadLine()
				if err != nil {
					if err == io.EOF {
						a.GetLogger(LogLevelErr).Errorln(gohelp.AnsiYellow+"Client connection closed"+gohelp.AnsiReset)
					} else {
						a.GetLogger(LogLevelErr).Errorln(gohelp.AnsiRed+err.Error()+gohelp.AnsiReset)
					}
					break
				}
				commands := strings.Split(string(com), CommandDelimiter)
				for _, comm := range commands {
					comm = strings.Trim(comm, " 	")
					// Parse command and run
					if comm == "" {
						continue
					}
					command := ParseCommand([]byte(comm))
					command.BindConnection(c)
					callback(command)
				}
			}
		}(conn)
	}
	return e
}
