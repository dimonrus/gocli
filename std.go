package gocli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/dimonrus/gohelp"
	"github.com/dimonrus/porterr"
	"gopkg.in/yaml.v3"
	"io"
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

var (
	// Depends config files
	RegExpDepends, _ = regexp.Compile(`depends:(.*)`)
	// ENV variables in config
	RegExpENV, _ = regexp.Compile(`\$\{(.*)\}`)
)

// DNApp Dynamic Name Application
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

// NewApplication Create new Application
func NewApplication(env string, configPath string, values interface{}) Application {
	app := &DNApp{
		config: config{
			values: values,
			path:   configPath,
		},
	}
	return app.ParseConfig(env)
}

// GetConfig Get config struct
func (a *DNApp) GetConfig() interface{} {
	return a.config.values
}

// SetConfig Set config struct
func (a *DNApp) SetConfig(cfg interface{}) Application {
	a.config.values = cfg
	return a
}

// GetConfigPath Get full path to config
func (a *DNApp) GetConfigPath(env string) string {
	return fmt.Sprintf("%s/%s.yaml", a.config.path, env)
}

// GetAbsolutePath Get absolute path to application
func (a *DNApp) GetAbsolutePath(path string, dir string) (string, porterr.IError) {
	rootPath, err := filepath.Abs("")
	if err != nil {
		return "", porterr.New(porterr.PortErrorArgument, "root path is incorrect")
	}
	if rootPath[len(rootPath)-1:] != string(os.PathSeparator) {
		rootPath = rootPath + string(os.PathSeparator)
	}
	return gohelp.BeforeString(rootPath, dir) + dir + string(os.PathSeparator) + path, nil
}

// FatalError Fatal error
func (a *DNApp) FatalError(err error) {
	panic(err)
}

// GetLogger Get logger
func (a *DNApp) GetLogger() Logger {
	if a.logger == nil {
		a.logger = NewLogger(LoggerConfig{})
	}
	return a.logger
}

// SetLogger Set logger
func (a *DNApp) SetLogger(logger Logger) {
	a.logger = logger
	return
}

// SuccessMessage printing success message
func (a *DNApp) SuccessMessage(message string, command ...*Command) {
	message = gohelp.AnsiGreen + message + gohelp.AnsiReset
	a.GetLogger().Output(DefaultCallDepth, message)
	for _, c := range command {
		e := c.Result([]byte(message + "\n"))
		if e != nil {
			a.GetLogger().Errorln(e)
		}
	}
	return
}

// AttentionMessage printing attention message
func (a *DNApp) AttentionMessage(message string, command ...*Command) {
	message = gohelp.AnsiCyan + message + gohelp.AnsiReset
	a.GetLogger().Output(DefaultCallDepth, message)
	for _, c := range command {
		e := c.Result([]byte(message + "\n"))
		if e != nil {
			a.GetLogger().Errorln(e)
		}
	}
	return
}

// FailMessage printing fail message
func (a *DNApp) FailMessage(message string, command ...*Command) {
	message = gohelp.AnsiRed + message + gohelp.AnsiReset
	a.GetLogger().Output(DefaultCallDepth, message)
	for _, c := range command {
		e := c.Result([]byte(message + "\n"))
		if e != nil {
			a.GetLogger().Errorln(e)
		}
	}
	return
}

// ParseConfig parse config depends on env
func (a *DNApp) ParseConfig(env string) Application {
	data, err := os.ReadFile(a.GetConfigPath(env))
	if err != nil {
		a.FatalError(err)
	}
	content := string(data)
	// check if config has depends on other configs
	matches := RegExpDepends.FindStringSubmatch(content)
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		// load parent config
		a.ParseConfig(strings.TrimSpace(matches[1]))
	}
	envMatches := RegExpENV.FindAllStringSubmatch(content, -1)
	for _, m := range envMatches {
		v, ok := os.LookupEnv(m[1])
		if ok {
			content = strings.ReplaceAll(content, m[0], v)
		} else {
			a.FailMessage("Environment: " + m[1] + " is not defined")
		}
	}
	// unmarshal config file in config struct
	err = yaml.Unmarshal([]byte(content), a.config.values)
	if err != nil {
		a.FatalError(err)
	}
	return a
}

// ParseFlags parse console arguments
func (a *DNApp) ParseFlags(args *Arguments) {
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

// Start run application
func (a *DNApp) Start(address string, callback func(command *Command)) porterr.IError {
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
			a.GetLogger().Errorln(err)
		}
	}()
	a.GetLogger().Infof("Start listening %s commands on %s:%s", CommandSessionType, CommandSessionHost, port)
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
					a.GetLogger().Errorln("Command processor error:", err)
				}
				// Always close the connection after process command
				err = c.Close()
				if err != nil {
					a.GetLogger().Errorln(err)
					return
				}
			}()
			r := bufio.NewReader(c)
			for {
				com, _, err := r.ReadLine()
				if err != nil {
					if err == io.EOF {
						a.GetLogger().Errorln(gohelp.AnsiYellow + "Client connection closed" + gohelp.AnsiReset)
					} else {
						a.GetLogger().Errorln(gohelp.AnsiRed + err.Error() + gohelp.AnsiReset)
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
