package category

import "github.com/jho3r/finanger-back/internal/infrastructure/logger"

var loggerService = logger.Setup("domains.category.service")

// Service is the interface for the category service.
type Service interface {
	CreateCategory(category Category) error
	GetCategories(catType CategoryType, name string) ([]Category, error)
	SeedCategories() error
}

// ServiceImpl is the struct that contains the category service.
type ServiceImpl struct {
	repo Repository
}

// NewCategoryService creates a new category service.
func NewCategoryService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// CreateCategory creates a new category.
func (s *ServiceImpl) CreateCategory(category Category) error {
	if err := s.repo.Create(category); err != nil {
		loggerService.WithError(err).Error("Error creating the category")
		return err
	}

	return nil
}

// GetCategories returns all the categories given filters.
func (s *ServiceImpl) GetCategories(catType CategoryType, name string) ([]Category, error) {
	categories, err := s.repo.Get(catType, name)
	if err != nil {
		loggerService.WithError(err).Error("Error getting the categories from the repo")
		return nil, err
	}

	return categories, nil
}

// SeedCategories seeds the categories into the database.
func (s *ServiceImpl) SeedCategories() error {
	loggerService.Info("Seeding the categories")

	assets, err := s.repo.Get(Asset, "")
	if err != nil {
		loggerService.WithError(err).Error("Error getting the assets from the repo")

		return err
	}

	if len(assets) == 0 {
		if err := s.repo.CreateMultiple(Assets); err != nil {
			loggerService.WithError(err).Error("Error creating the assets")

			return err
		}
	}

	return nil
}
