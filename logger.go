package gocli

import (
	"context"
	"fmt"
	"log"
	"os"
)

const (
	LogLevelDebug = 1 << iota
	LogLevelInfo
	LogLevelWarn
	LogLevelErr

	DefaultCallDepth = 3

	DefaultKeyLoggerPrefix = "prefix"
)

// Logger Common logger interface
type Logger interface {
	GetConfig() LoggerConfig
	Output(callDepth int, message string) error

	Print(v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})

	Info(v ...interface{})
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})

	Warn(v ...interface{})
	Warnln(v ...interface{})
	Warnf(format string, v ...interface{})

	Error(v ...interface{})
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{})
}

// LoggerFormat logger prefix
type LoggerFormat map[string]string

// FromContext create log prefix from context
func (lf LoggerFormat) FromContext(ctx context.Context) string {
	var prefix string
	for key, value := range lf {
		if v, ok := ctx.Value(key).(string); ok {
			prefix += fmt.Sprintf(value, v) + " "
		}
	}
	return prefix
}

// String serialize format to string
func (lf LoggerFormat) String() string {
	if v, ok := lf[DefaultKeyLoggerPrefix]; ok {
		return v + " "
	}
	return ""
}

// LoggerConfig configuration of logger
type LoggerConfig struct {
	// Level of log message
	Level int
	// Default is 2
	Depth int
	// Flags
	Flags int
	// Format for multiple arguments
	Format LoggerFormat
}

// NewLogger Init logger struct
func NewLogger(config LoggerConfig) Logger {
	config.Depth = config.Depth | DefaultCallDepth
	config.Level = config.Level | LogLevelDebug
	config.Flags = config.Flags | log.Ldate | log.Ltime | log.Lshortfile
	return &logger{config: config, stdLogger: log.New(os.Stdout, config.Format.String(), config.Flags)}
}

// logger struct
type logger struct {
	config    LoggerConfig
	stdLogger *log.Logger
}

// GetConfig return logger config
func (l logger) GetConfig() LoggerConfig {
	return l.config
}

// Output printing message
func (l logger) Output(callDepth int, message string) error {
	return l.stdLogger.Output(callDepth, message)
}

// Print printing message
func (l logger) Print(v ...interface{}) {
	if len(v) > 0 {
		if c, ok := v[0].(context.Context); ok {
			l.Output(l.config.Depth, l.config.Format.FromContext(c)+fmt.Sprint(v[1:]...))
		} else {
			l.Output(l.config.Depth, fmt.Sprint(v...))
		}
	}
}

// Info printing message at info level
func (l logger) Info(v ...interface{}) {
	if (l.config.Level & (LogLevelInfo | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c)+fmt.Sprint(v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprint(v...))
			}
		}
	}
}

// Warn printing message at warn level
func (l logger) Warn(v ...interface{}) {
	if (l.config.Level & (LogLevelWarn | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c)+fmt.Sprint(v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprint(v...))
			}
		}
	}
}

// Error printing message at error level
func (l logger) Error(v ...interface{}) {
	if (l.config.Level & (LogLevelErr | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c)+fmt.Sprint(v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprint(v...))
			}
		}
	}
}

// Println printing message with new line symbol
func (l logger) Println(v ...interface{}) {
	if len(v) > 0 {
		if c, ok := v[0].(context.Context); ok {
			l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintln(v[1:]...))
		} else {
			l.Output(l.config.Depth, fmt.Sprintln(v...))
		}
	}
}

// Infoln printing message with new line symbol at info level
func (l logger) Infoln(v ...interface{}) {
	if (l.config.Level & (LogLevelInfo | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintln(v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprintln(v...))
			}
		}
	}
}

// Warnln printing message with new line symbol at warn level
func (l logger) Warnln(v ...interface{}) {
	if (l.config.Level & (LogLevelWarn | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintln(v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprintln(v...))
			}
		}
	}
}

// Errorln printing message with new line symbol at error level
func (l logger) Errorln(v ...interface{}) {
	if (l.config.Level & (LogLevelErr | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintln(v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprintln(v...))
			}
		}
	}
}

// Printf printing message in custom format
func (l logger) Printf(format string, v ...interface{}) {
	if len(v) > 0 {
		if c, ok := v[0].(context.Context); ok {
			l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintf(format, v[1:]...))
		} else {
			l.Output(l.config.Depth, fmt.Sprintf(format, v...))
		}
	}
}

// Infof printing message in custom format at info level
func (l logger) Infof(format string, v ...interface{}) {
	if (l.config.Level & (LogLevelInfo | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintf(format, v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprintf(format, v...))
			}
		}
	}
}

// Warnf printing message in custom format at warn level
func (l logger) Warnf(format string, v ...interface{}) {
	if (l.config.Level & (LogLevelWarn | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintf(format, v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprintf(format, v...))
			}
		}
	}
}

// Errorf printing message in custom format at error level
func (l logger) Errorf(format string, v ...interface{}) {
	if (l.config.Level & (LogLevelErr | LogLevelDebug)) != 0 {
		if len(v) > 0 {
			if c, ok := v[0].(context.Context); ok {
				l.Output(l.config.Depth, l.config.Format.FromContext(c) + fmt.Sprintf(format, v[1:]...))
			} else {
				l.Output(l.config.Depth, fmt.Sprintf(format, v...))
			}
		}
	}
}
