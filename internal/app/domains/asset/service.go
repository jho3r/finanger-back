package asset

import "github.com/jho3r/finanger-back/internal/infrastructure/logger"

var loggerService = logger.Setup("domain.assets.service")

// Service is the interface for the asset service.
type Service interface {
	Create(asset Asset) error
	GetAssets(userID uint, finAssetID uint, name string, categoryID uint) ([]Asset, error)
	GetAssetByID(userID uint, assetID uint) (Asset, error)
	Update(asset Asset) error
	Delete(userID uint, assetID uint) error
}

// ServiceImpl is the struct that contains the asset service.
type ServiceImpl struct {
	repo Repository
}

// NewAssetService creates a new asset service.
func NewAssetService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// Create creates a new asset.
func (s *ServiceImpl) Create(asset Asset) error {
	if err := s.repo.Create(asset); err != nil {
		return err
	}

	return nil
}

// GetAssets gets all the assets.
func (s *ServiceImpl) GetAssets(userID uint, finAssetID uint, name string, categoryID uint) ([]Asset, error) {
	assets, err := s.repo.GetAssets(userID, finAssetID, name, categoryID)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

// GetAssetByID gets an asset by its ID.
func (s *ServiceImpl) GetAssetByID(userID uint, assetID uint) (Asset, error) {
	asset, err := s.repo.GetAssetByID(userID, assetID)
	if err != nil {
		return Asset{}, err
	}

	return asset, nil
}

// Update updates an asset.
func (s *ServiceImpl) Update(asset Asset) error {
	if err := s.repo.Update(asset); err != nil {
		return err
	}

	return nil
}

// Delete deletes an asset.
func (s *ServiceImpl) Delete(userID uint, assetID uint) error {
	if err := s.repo.Delete(userID, assetID); err != nil {
		return err
	}

	return nil
}
