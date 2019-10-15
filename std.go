package gocli

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dimonrus/porterr"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"testing"
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
func (a DNApp) Start(arguments Arguments) porterr.IError {
	return nil
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
	err = yaml.Unmarshal([]byte(data), config)
	if err != nil {
		a.FatalError(err)
	}
	a.config = config

	return &a
}

// Parse console arguments
func (a DNApp) ParseFlags(args *Arguments) {
	for key, argument := range *args {
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
		default:
			a.FatalError(errors.New("not implemented argument type"))
		}
	}
	testing.Init()
	flag.Parse()
}
