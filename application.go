package gocli

import "github.com/dimonrus/porterr"

// Application class
type Application interface {
	// Get config struct
	GetConfig() interface{}
	// Get full path to config
	GetConfigPath(env string) string
	// Get absolute path
	GetAbsolutePath(path string, dir string) (string, porterr.IError)
	// Set config struct
	SetConfig(cfg interface{}) Application
	// Parse config
	ParseConfig(env string) Application
	// Start application
	Start(port string, callback func(command *Command)) porterr.IError
	// Behaviour for fatal errors
	FatalError(err error)
	// Get Logger
	GetLogger(level int) Logger
	// Success log message with command repeat
	SuccessMessage(message string, command ...*Command)
	// Warning log message with command repeat
	AttentionMessage(message string, command ...*Command)
	// Fail log message with command repeat
	FailMessage(message string, command ...*Command)
	// Parse console flags
	ParseFlags(args *Arguments)
}