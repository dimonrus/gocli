package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// Dynamic Name Application
type DNApp struct {
	config interface{}
	logger *logger
}

// Get config struct
func (a DNApp) GetConfig() interface{} {
	return a.config
}

// Get full path to config
func (a DNApp) GetConfigPath(env string) string {
	rootPath, err := filepath.Abs("")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s%s", rootPath, "/config/yaml/"+env+".yaml")
}

// Set config struct
func (a DNApp) SetConfig(cfg interface{}) Application {
	a.config = cfg
	return &a
}

// Start application
func (a DNApp) Start(arguments Arguments) IError {
	return nil
}

// Fatal error
func (a DNApp) FatalError(err error) {
	panic(err)
}

// Get logger
func (a DNApp) GetLogger(level int) *logger {
	if a.logger == nil || a.logger.level != level {
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
	flag.Parse()

}
