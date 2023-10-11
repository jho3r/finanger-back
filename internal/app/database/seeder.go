package database

import (
	"github.com/jho3r/finanger-back/internal/app/domains/category"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerSeeder = logger.Setup("cmd.seeder")

// Seed seeds the database.
func Seed() error {
	loggerSeeder.Info("Seeding the database")

	// Infrastructure
	gormDB := gorm.NewGormDB(settings.Database.ConnURL, settings.Database.MaxIdle, settings.Database.MaxOpen)

	// Repos
	categoryRepo := category.NewCategoryRepository(gormDB)

	// Services
	categoryService := category.NewCategoryService(categoryRepo)

	// Seed
	if err := categoryService.SeedCategories(); err != nil {
		loggerSeeder.WithError(err).Error("Error seeding the categories")

		return err
	}

	return nil
}
