package gocli

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestLog_Print(t *testing.T) {
	cfg := LoggerConfig{
		Level: LogLevelErr,
		Format: map[string]string{
			DefaultKeyLoggerPrefix: "Tester:",
			"x-trace-id":           "xid: %s",
			"macro":                "%s",
		},
	}

	ctx1 := context.WithValue(context.Background(), "x-trace-id", "24234234234")
	ctx := context.WithValue(ctx1, "macro", "rightId: 1")
	l := NewLogger(cfg)
	l.Print(ctx, "hi all")
	l.Print(ctx, "hi all")
	l.Println(ctx, "hi all")
	l.Printf("someone %s", ctx1, "hi all")
	l.Warnf("someone %s", "hi all")
	l.Warnf("show the line")
	l.Errorln(ctx, "show error")
}

// BenchmarkDefaultLogger-4   	   99524	     10322 ns/op
func BenchmarkDefaultLogger(b *testing.B) {
	l := log.New(os.Stdout, "Tester: ", log.Ldate|log.Ltime|log.Lshortfile)
	for i := 0; i < b.N; i++ {
		l.Println("xid: 24234234234", "rightId: 1")
	}
}

// BenchmarkLogger-4   	   94099	     11406 ns/op
func BenchmarkLogger(b *testing.B) {
	l := NewLogger(LoggerConfig{Format: LoggerFormat{
		"prefix": "Tester:",
	}})
	for i := 0; i < b.N; i++ {
		l.Println("xid: 24234234234", "rightId: 1")
	}
	b.ReportAllocs()
}

// BenchmarkContextLogger-4   	  102255	     10820 ns/op
func BenchmarkContextLogger(b *testing.B) {
	l := NewLogger(LoggerConfig{Format: LoggerFormat{
		DefaultKeyLoggerPrefix: "Tester:",
		"x-trace-id":           "xid: %s",
		"macro":                "%s",
	}})

	ctx1 := context.WithValue(context.Background(), "x-trace-id", "24234234234")
	ctx := context.WithValue(ctx1, "macro", "rightId: 1")

	for i := 0; i < b.N; i++ {
		l.Println(ctx)
	}
}
