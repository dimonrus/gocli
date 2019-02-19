package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// Dynamic Name Application
type DNApp struct {
	config interface{}
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
func (a DNApp) Start(arguments Arguments) ApplicationErrorInterface {
	return nil
}

// Init app
func (a DNApp) New(env string, config interface{}) Application {
	data, err := ioutil.ReadFile(a.GetConfigPath(env))
	if err != nil {
		panic(err)
	}
	r, _ := regexp.Compile(`depends:(.*)`)
	matches := r.FindStringSubmatch(string(data))
	if len(matches) > 1 && strings.TrimSpace(matches[1]) != "" {
		a.New(strings.TrimSpace(matches[1]), config)
	}
	err = yaml.Unmarshal([]byte(data), config)
	if err != nil {
		panic(err)
	}
	return &a
}
