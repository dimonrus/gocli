package gocli

import "github.com/dimonrus/porterr"

// Application interface
type Application interface {
	// GetConfig Get config struct
	GetConfig() interface{}
	// GetConfigPath Get full path to config
	GetConfigPath(env string) string
	// GetAbsolutePath Get absolute path
	GetAbsolutePath(path string, dir string) (string, porterr.IError)
	// SetConfig Set config struct
	SetConfig(cfg interface{}) Application
	// ParseConfig Parse config
	ParseConfig(env string) Application
	// Start run application
	Start(port string, callback func(command *Command)) porterr.IError
	// FatalError Behaviour for fatal errors
	FatalError(err error)
	// GetLogger Get Logger
	GetLogger() Logger
	// SetLogger set custom logger
	SetLogger(logger Logger)
	// SuccessMessage Success log message with command repeat
	SuccessMessage(message string, command ...*Command)
	// AttentionMessage Warning log message with command repeat
	AttentionMessage(message string, command ...*Command)
	// FailMessage Fail log message with command repeat
	FailMessage(message string, command ...*Command)
	// ParseFlags Parse console flags
	ParseFlags(args *Arguments)
}