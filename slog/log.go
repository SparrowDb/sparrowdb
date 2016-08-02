package slog

import (
	"flag"
	"fmt"
	stdLog "log"

	"github.com/golang/glog"
)

// Logger logging interface
type Logger interface {
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

var (
	log Logger = stdLogger{}
)

// SetLogger sets logger
func SetLogger(lo Logger) {
	log = lo
}

// Infof prints INFO message
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf prints WARN message
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf prints ERROR message
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf prints FATAL message
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Holds glog logger
type glogLogger struct{}

func (glogLogger) Infof(format string, args ...interface{}) {
	glog.InfoDepth(3, fmt.Sprintf(format, args...))
}

func (glogLogger) Warnf(format string, args ...interface{}) {
	glog.WarningDepth(3, fmt.Sprintf(format, args...))
}

func (glogLogger) Errorf(format string, args ...interface{}) {
	glog.ErrorDepth(3, fmt.Sprintf(format, args...))
}

func (glogLogger) Fatalf(format string, args ...interface{}) {
	glog.FatalDepth(3, fmt.Sprintf(format, args...))
}

// NewGlog returns new glog logger
func NewGlog() glogLogger {
	flag.Lookup("logtostderr").Value.Set("true")
	return glogLogger{}
}

// Holds default logger
type stdLogger struct{}

func (stdLogger) Infof(format string, args ...interface{}) {
	stdLog.Printf("INFO: "+format, args...)
}

func (stdLogger) Warnf(format string, args ...interface{}) {
	stdLog.Printf("WARN: "+format, args...)
}

func (stdLogger) Errorf(format string, args ...interface{}) {
	stdLog.Printf("ERROR: "+format, args...)
}

func (stdLogger) Fatalf(format string, args ...interface{}) {
	stdLog.Fatalf("FATAL: "+format, args...)
}
