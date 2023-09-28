package logger

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	// Add any custom fields or settings you need here.
}

// Define ANSI escape codes for colors.
const (
	InfoColor    = "\033[1;34m" // Blue
	WarningColor = "\033[1;33m" // Yellow
	ErrorColor   = "\033[1;31m" // Red
	ResetColor   = "\033[0m"    // Reset color
)

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Define a color variable based on the log level.
	var color string

	switch entry.Level {
	case logrus.InfoLevel:
		color = InfoColor
	case logrus.WarnLevel:
		color = WarningColor
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = ErrorColor
	default:
		color = ResetColor
	}

	level := strings.ToUpper(entry.Level.String())

	// Get the name of the function that called the logger.
	// entry.Caller.Function = github.com/jho3r/finanger-back/internal/app/settings.LoadEnvs just get the last part "LoadEnvs"
	funcName := strings.Split(entry.Caller.Function, ".")[len(strings.Split(entry.Caller.Function, "."))-1]

	// Create a formatted log message with color.
	formattedLog := fmt.Sprintf("%s[%s][%s][%s.%s][%s] -- msg: %s%s",
		color,
		level,
		entry.Data["app"],
		entry.Data["service"],
		funcName,
		entry.Data["uuid"],
		ResetColor,
		entry.Message,
	)

	// If there's an "error" field in the log entry, include it in the formatted message.
	if err, ok := entry.Data["error"]; ok {
		formattedLog = fmt.Sprintf("%s -- %serror: %s%s", formattedLog, color, ResetColor, err)
	}

	// Append a newline character to the log message.
	formattedLog = fmt.Sprintf("%s\n", formattedLog)

	return []byte(formattedLog), nil
}
