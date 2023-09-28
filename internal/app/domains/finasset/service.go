package finasset

import (
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerService = logger.Setup("domain.finasset.service")

// Service is the interface for the financial asset service.
type Service interface {
	Create(finAsset FinancialAsset) error
	Get(finAsset FinancialAsset) ([]FinancialAsset, error)
}

// ServiceImpl is the struct that contains the financial asset service.
type ServiceImpl struct {
	repo Repository
}

// NewFinAssetService creates a new financial asset service.
func NewFinAssetService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// Create creates a new financial asset.
func (s *ServiceImpl) Create(finAsset FinancialAsset) error {
	if err := s.repo.Create(finAsset); err != nil {
		loggerService.WithError(err).Error("Error creating the financial asset")
		return err
	}

	return nil
}

// Get returns all the financial assets given filters.
func (s *ServiceImpl) Get(finAsset FinancialAsset) ([]FinancialAsset, error) {
	finAssets, err := s.repo.Get(finAsset)
	if err != nil {
		loggerService.WithError(err).Error("Error getting the financial assets from the repo")
		return nil, err
	}

	return finAssets, nil
}
