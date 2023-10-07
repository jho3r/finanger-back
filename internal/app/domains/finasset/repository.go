package finasset

import (
	"strings"

	"github.com/jho3r/finanger-back/internal/infrastructure/database/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerRepo = logger.Setup("domain.finasset.repository")

// Repository is the interface for the financial asset repository.
type Repository interface {
	Create(finAsset FinancialAsset) error
	Get(finAsset FinancialAsset) ([]FinancialAsset, error)
}

// RepositoryImpl is the struct that contains the financial asset repository.
type RepositoryImpl struct {
	db gorm.Gorm
}

// NewCurrencyRepository creates a new financial asset repository.
func NewCurrencyRepository(db gorm.Gorm) Repository {
	return &RepositoryImpl{db: db}
}

// Create creates a new financial asset.
func (r *RepositoryImpl) Create(finAsset FinancialAsset) error {
	if err := r.db.Create(&finAsset); err != nil {
		loggerRepo.WithError(err).Error("Error creating record in the database")

		return err
	}

	return nil
}

// Get returns all the financial assets given filters.
// Filters as equal type, and like symbol and name.
func (r *RepositoryImpl) Get(finAsset FinancialAsset) ([]FinancialAsset, error) {
	var finAssets []FinancialAsset

	var queryConditions []string
	var args []interface{}

	if finAsset.Type != "" {
		queryConditions = append(queryConditions, "type = ?")
		args = append(args, finAsset.Type)
	}

	if finAsset.Symbol != "" {
		queryConditions = append(queryConditions, "symbol LIKE ?")
		args = append(args, "%"+finAsset.Symbol+"%")
	}

	if finAsset.Name != "" {
		queryConditions = append(queryConditions, "name LIKE ?")
		args = append(args, "%"+finAsset.Name+"%")
	}

	query := strings.Join(queryConditions, " AND ")

	if err := r.db.WhereFind(&finAssets, false, query, args...); err != nil {
		loggerRepo.WithError(err).Error("Error getting records from the database")

		return nil, err
	}

	return finAssets, nil
}
