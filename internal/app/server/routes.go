package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/controller"
	"github.com/jho3r/finanger-back/internal/app/domains/asset"
	"github.com/jho3r/finanger-back/internal/app/domains/category"
	"github.com/jho3r/finanger-back/internal/app/domains/finasset"
	"github.com/jho3r/finanger-back/internal/app/domains/user"
	"github.com/jho3r/finanger-back/internal/app/middlewares"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerServer = logger.Setup("server")

// SetupServer Init the server with the middlewares and the routes.
func SetupServer() *gin.Engine {
	loggerServer.Info("Initializing server ...")

	basePath := fmt.Sprintf("/api/%s", settings.Commons.ProjectName)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(basePath))
	router.Use(middlewares.CorsMiddleware(middlewares.SplitOrigins(settings.Auth.AllowedOrigins)))

	// Dependencies

	// Infrastructure
	gormDB := gorm.NewGormDB(settings.Database.ConnURL, settings.Database.MaxIdle, settings.Database.MaxOpen)

	// Repos
	userRepo := user.NewUserRepository(gormDB)
	finAssetRepo := finasset.NewCurrencyRepository(gormDB)
	assetRepo := asset.NewAssetRepository(gormDB)
	categoryRepo := category.NewCategoryRepository(gormDB)

	// Services
	userService := user.NewUserService(userRepo)
	finAssetService := finasset.NewFinAssetService(finAssetRepo)
	assetService := asset.NewAssetService(assetRepo)
	categoryService := category.NewCategoryService(categoryRepo)

	// Routes

	base := router.Group(basePath)
	base.GET("/health", controller.HealthCheck)

	finassets := base.Group("/financial-assets")
	finassets.POST("/", middlewares.AuthApiKey(), controller.CreateFinancialAsset(finAssetService))
	finassets.GET("/", controller.GetFinancialAssets(finAssetService))

	users := base.Group("/users")
	users.POST("/signup", controller.Signup(userService))
	users.POST("/login", controller.Login(userService))
	users.POST("/refresh-token", middlewares.AuthUser(), controller.RefreshToken(userService))
	users.GET("/me", middlewares.AuthUser(), controller.GetMe(userService))
	users.POST("/logout", middlewares.AuthUser(), controller.Logout(userService))

	assets := base.Group("/assets")
	assets.POST("/", middlewares.AuthUser(), controller.CreateAsset(assetService))
	assets.GET("/", middlewares.AuthUser(), controller.GetAssets(assetService))
	assets.GET("/:id", middlewares.AuthUser(), controller.GetAsset(assetService))
	assets.PUT("/:id", middlewares.AuthUser(), controller.UpdateAsset(assetService))
	assets.DELETE("/:id", middlewares.AuthUser(), controller.DeleteAsset(assetService))

	categories := base.Group("/categories")
	categories.POST("/", middlewares.AuthApiKey(), controller.CreateCategory(categoryService))
	categories.GET("/", controller.GetCategories(categoryService))

	return router
}
