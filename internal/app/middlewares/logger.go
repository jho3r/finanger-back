package middlewares

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerMidLogger = logger.Setup("middlewares.logger")

func LoggerMiddleware(basePath string) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{basePath + "/health"},
		Formatter: formatter,
	})
}

// formatter format the log from the gin server.
func formatter(param gin.LogFormatterParams) string {
	output := new(bytes.Buffer)
	log := logger.SetID(loggerMidLogger, "GIN")
	log.Logger.SetOutput(output)

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	log.Infof(
		"| %3d | %s | %5v | %s",
		param.StatusCode,
		param.Method,
		param.Latency,
		param.Path,
	)

	return output.String()
}
