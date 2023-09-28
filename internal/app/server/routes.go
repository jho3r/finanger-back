package server

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/controller"
	"github.com/jho3r/finanger-back/internal/app/domains/finasset"
	"github.com/jho3r/finanger-back/internal/app/domains/user"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/database/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
	"github.com/sirupsen/logrus"
)

var loggerServer = logger.Setup("server")

// SetupServer Init the server with the middlewares and the routes
func SetupServer() *gin.Engine {
	loggerServer.Info("Initializing server ...")

	basePath := fmt.Sprintf("/api/%s", settings.Commons.ProjectName)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{basePath + "/health"},
		Formatter: formatter,
	}))

	// Dependencies

	// Infrastructure
	gormDB := gorm.NewGormDB(settings.Database.ConnURL, settings.Database.MaxIdle, settings.Database.MaxOpen)

	// Repos
	userRepo := user.NewUserRepository(gormDB)
	finAssetRepo := finasset.NewCurrencyRepository(gormDB)

	// Services
	userService := user.NewUserService(userRepo)
	finAssetService := finasset.NewFinAssetService(finAssetRepo)

	// Routes

	base := router.Group(basePath)
	base.GET("/health", controller.HealthCheck)

	finassets := base.Group("/financial-assets")
	finassets.POST("/", controller.CreateFinancialAsset(finAssetService))
	finassets.GET("/", controller.GetFinancialAssets(finAssetService))

	users := base.Group("/users")
	users.POST("/signup", controller.Signup(userService))

	return router
}

// formatter format the log from the gin server.
func formatter(param gin.LogFormatterParams) string {
	output := new(bytes.Buffer)
	log := logger.SetID(loggerServer, "GIN")
	log.Logger.SetLevel(logrus.InfoLevel)
	log.Logger.SetOutput(output)

	var statusColor, methodColor, resetColor, bold string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
		bold = "\033[1m"
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	log.Infof(
		"| %s %3d %s | %s %s %s | %10v | %s %s %s | CLIENT_IP=%s X_APP_ID=%s",
		statusColor, param.StatusCode, resetColor,
		methodColor, param.Method, resetColor,
		param.Latency,
		bold, param.Path, resetColor,
		param.ClientIP,
		param.Request.Header.Get("X-Application-ID"),
	)

	return output.String()
}
