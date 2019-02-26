package main

import (
	"log"
	"os"
)

const (
	LogLevelDebug = 1 << iota
	LogLevelInfo
	LogLevelWarn
	LogLevelErr
)

// New logger
func NewLogger(level int, prefix string, flags int) *logger {
	return &logger{level: level, stdLogger: log.New(os.Stdout, prefix, flags)}
}

// logger struct
type logger struct {
	level     int
	stdLogger *log.Logger
}

// Print
func (l logger) Print(v ...interface{}) {
	l.stdLogger.Print(v...)
}

// Info
func (l logger) Info(v ...interface{}) {
	if (l.level & (LogLevelInfo | LogLevelDebug)) != 0 {
		l.stdLogger.Print(v...)
	}
}

// Warning
func (l logger) Warn(v ...interface{}) {
	if (l.level & (LogLevelWarn | LogLevelDebug)) != 0 {
		l.stdLogger.Print(v...)
	}
}

// Error logging
func (l logger) Error(v ...interface{}) {
	if (l.level & (LogLevelErr | LogLevelDebug)) != 0 {
		l.stdLogger.Print(v...)
	}
}

// Print ln
func (l logger) Println(v ...interface{}) {
	l.stdLogger.Println(v...)
}

// Info ln
func (l logger) Infoln(v ...interface{}) {
	if (l.level & (LogLevelInfo | LogLevelDebug)) != 0 {
		l.stdLogger.Println(v...)
	}
}

// Warning ln
func (l logger) Warnln(v ...interface{}) {
	if (l.level & (LogLevelWarn | LogLevelDebug)) != 0 {
		l.stdLogger.Println(v...)
	}
}

// Error ln
func (l logger) Errorln(v ...interface{}) {
	if (l.level & (LogLevelErr | LogLevelDebug)) != 0 {
		l.stdLogger.Println(v...)
	}
}

// Print format
func (l logger) Printf(format string, v ...interface{}) {
	l.stdLogger.Printf(format, v...)
}

// Info format
func (l logger) Infof(format string, v ...interface{}) {
	if (l.level & (LogLevelInfo | LogLevelDebug)) != 0 {
		l.stdLogger.Printf(format, v...)
	}
}

// Warning format
func (l logger) Warnf(format string, v ...interface{}) {
	if (l.level & (LogLevelWarn | LogLevelDebug)) != 0 {
		l.stdLogger.Printf(format, v...)
	}
}

// Error format
func (l logger) Errorf(format string, v ...interface{}) {
	if (l.level & (LogLevelErr | LogLevelDebug)) != 0 {
		l.stdLogger.Printf(format, v...)
	}
}
