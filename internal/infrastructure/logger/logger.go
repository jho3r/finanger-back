// package to log Correctly in logdna with version of the project, service and trace the all logs of a process
package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Setup Init the logger with the service name
// All logs will be tagged with the app, service, uuid and block
func Setup(service string) *log.Entry {
	var logger = log.Logger{
		Out:          os.Stdout,
		Level:        log.InfoLevel,
		ReportCaller: true,
		Formatter:    &CustomFormatter{},
	}

	return logger.WithFields(log.Fields{
		"app":     os.Getenv("X_APPLICATION_ID"),
		"service": service, // This is for the service name like controller.health
		"uuid":    "",      // This is for the trace of the process
		"block":   "",      // This is for the block of the process like a function name
	})
}

// SetID Set and Id to trace the log of a process
// Look loggerExample_test.go
func SetID(log *log.Entry, id string) *log.Entry {
	return log.WithField("uuid", id)
}
