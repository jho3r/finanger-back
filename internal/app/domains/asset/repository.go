package asset

import (
	"strings"

	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerRepo = logger.Setup("domain.assets.repository")

// Repository is the interface for the asset repository.
type Repository interface {
	Create(asset Asset) error
	GetAssets(userID uint, finAssetID uint, name string) ([]Asset, error)
	GetAssetByID(userID uint, assetID uint) (Asset, error)
	Update(asset Asset) error
	Delete(userID uint, assetID uint) error
}

// RepositoryImpl is the struct that contains the asset repository.
type RepositoryImpl struct {
	db gorm.Gorm
}

// NewAssetRepository creates a new asset repository.
func NewAssetRepository(db gorm.Gorm) Repository {
	return &RepositoryImpl{db: db}
}

// Create creates a new asset.
func (r *RepositoryImpl) Create(asset Asset) error {
	if err := r.db.Create(&asset); err != nil {
		loggerRepo.WithError(err).Error("Error creating record in the database")

		return err
	}

	return nil
}

// GetAssets gets all the assets.
func (r *RepositoryImpl) GetAssets(userID uint, finAssetID uint, name string) ([]Asset, error) {
	var assets []Asset

	queryConditions := []string{"user_id = ?"}
	args := []interface{}{userID}

	if finAssetID != 0 {
		queryConditions = append(queryConditions, "financial_asset_id = ?")
		args = append(args, finAssetID)
	}

	if name != "" {
		queryConditions = append(queryConditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	query := strings.Join(queryConditions, " AND ")

	if err := r.db.WhereFind(&assets, true, query, args...); err != nil {
		loggerRepo.WithError(err).Error("Error getting records from the database")

		return nil, err
	}

	return assets, nil
}

// GetAssetByID gets an asset by its ID.
func (r *RepositoryImpl) GetAssetByID(userID uint, assetID uint) (Asset, error) {
	var asset Asset

	query := "user_id = ? AND id = ?"

	if err := r.db.WhereFirst(&asset, query, userID, assetID); err != nil {
		loggerRepo.WithError(err).Error("Error getting record from the database")

		return Asset{}, err
	}

	return asset, nil
}

// Update updates an asset.
func (r *RepositoryImpl) Update(asset Asset) error {
	if err := r.db.Save(&asset); err != nil {
		loggerRepo.WithError(err).Error("Error updating record in the database")

		return err
	}

	return nil
}

// Delete deletes an asset.
func (r *RepositoryImpl) Delete(userID uint, assetID uint) error {
	asset, err := r.GetAssetByID(userID, assetID)
	if err != nil {
		return err
	}

	if err := r.db.Delete(&asset); err != nil {
		loggerRepo.WithError(err).Error("Error deleting record in the database")

		return err
	}

	return nil
}
