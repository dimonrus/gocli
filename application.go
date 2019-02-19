package main

// Application class
type Application interface {
	// Get config struct
	GetConfig() interface{}
	// Get full path to config
	GetConfigPath(env string) string
	// Set config struct
	SetConfig(cfg interface{}) Application
	// Start application
	Start(arguments Arguments) ApplicationErrorInterface
	// Init app method
	New(env string, cfg interface{}) Application
}