package gocli

import "github.com/dimonrus/porterr"

// Application class
type Application interface {
	// Get config struct
	GetConfig() interface{}
	// Get full path to config
	GetConfigPath(env string) string
	// Set config struct
	SetConfig(cfg interface{}) Application
	// Start application
	Start(port string, callback func(command Command)) porterr.IError
	// Init app method
	New(env string, cfg interface{}) Application
	// Behaviour for fatal errors
	FatalError(err error)
	// Get Logger
	GetLogger(level int) Logger
	// Parse console flags
	ParseFlags(args *Arguments)
}
