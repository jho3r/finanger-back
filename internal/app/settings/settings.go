package settings

import (
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
	"github.com/kelseyhightower/envconfig"
)

var (
	settingsLogger = logger.Setup("settings")
	// Commons struct to store all the settings of the application.
	Commons commons
	// Database struct to store all the settings of the database.
	Database database
)

type commons struct {
	XApplicationID string `envconfig:"X_APPLICATION_ID" default:"finanger-back"`
	ProjectName    string `envconfig:"PROJECT_NAME" default:"finanger-back"`
	Port           string `envconfig:"PORT" required:"true"`
}

type database struct {
	ConnURL string `envconfig:"DATABASE_CONN_URL" required:"true"`
	MaxIdle int    `envconfig:"DATABASE_MAX_IDLE_CONNS" default:"10"`
	MaxOpen int    `envconfig:"DATABASE_MAX_OPEN_CONNS" default:"10"`
}

// LoadEnvs loads all the envs of the application.
func LoadEnvs() {
	// Load all the envs
	err := envconfig.Process("", &Commons)
	if err != nil {
		settingsLogger.WithError(err).Fatal("Error loading commons envs")
	}

	err = envconfig.Process("", &Database)
	if err != nil {
		settingsLogger.WithError(err).Fatal("Error loading database envs")
	}
}
