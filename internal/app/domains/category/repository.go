package category

import (
	"strings"

	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerRepo = logger.Setup("domains.category.repository")

// Repository is the interface for the category repository.
type Repository interface {
	Create(category Category) error
	Get(catType CategoryType, name string) ([]Category, error)
	CreateMultiple(categories []Category) error
}

// RepositoryImpl is the struct that contains the category repository.
type RepositoryImpl struct {
	db gorm.Gorm
}

// NewCategoryRepository creates a new category repository.
func NewCategoryRepository(db gorm.Gorm) Repository {
	return &RepositoryImpl{db: db}
}

// CreateCategory creates a new category.
func (r *RepositoryImpl) Create(category Category) error {
	if err := r.db.Create(&category); err != nil {
		loggerRepo.WithError(err).Error("Error creating the category")
		return err
	}

	return nil
}

// GetCategories returns all the categories given filters.
func (r *RepositoryImpl) Get(catType CategoryType, name string) ([]Category, error) {
	var categories []Category

	var queryConditions []string
	var args []interface{}

	if catType != "" {
		queryConditions = append(queryConditions, "type = ?")
		args = append(args, catType)
	}

	if name != "" {
		queryConditions = append(queryConditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	query := strings.Join(queryConditions, " AND ")

	if err := r.db.WhereFind(&categories, false, query, args...); err != nil {
		loggerRepo.WithError(err).Error("Error getting the categories from the database")
		return nil, err
	}

	return categories, nil
}

// CreateMultiple creates multiple categories.
func (r *RepositoryImpl) CreateMultiple(categories []Category) error {
	if err := r.db.Create(&categories); err != nil {
		loggerRepo.WithError(err).Error("Error creating the categories")
		return err
	}

	return nil
}
