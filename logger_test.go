package main

import (
	"log"
	"testing"
)

func TestLog_Print(t *testing.T) {
	l := NewLogger(LogLevelErr|LogLevelWarn|LogLevelInfo, "TestLogger", log.Flags())
	l.Println("hi all")
}
