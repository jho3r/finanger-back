package main

import (
	"github.com/jho3r/finanger-back/internal/app/database"
	"github.com/jho3r/finanger-back/internal/app/server"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerMain = logger.Setup("cmd.main")

func main() {
	loggerMain.Info("Initializing everything ... Wait for it ...")

	// Load envs
	settings.LoadEnvs()

	// Seed the database
	if err := database.Seed(); err != nil {
		loggerMain.WithError(err).Fatal("Error seeding the database")
	}

	// Setup the server
	router := server.SetupServer()

	loggerMain.Infof("Server running on port %s", settings.Commons.Port)

	if err := router.Run(":" + settings.Commons.Port); err != nil {
		loggerMain.WithError(err).Fatal("Error running the server")
	}
}
