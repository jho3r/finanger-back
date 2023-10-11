package asset

import (
	"github.com/jho3r/finanger-back/internal/app/domains/finasset"
	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
)

type (
	// Asset is the struct that contains the asset data.
	Asset struct {
		gorm.Model
		Name             string                  `json:"name" gorm:"not null"`
		Value            float64                 `json:"value" gorm:"not null"`
		Description      string                  `json:"description" gorm:"not null"`
		FinancialAssetID uint                    `json:"fin_asset_id" gorm:"not null"`
		FinancialAsset   finasset.FinancialAsset `json:"financial_asset" gorm:"foreignKey:FinancialAssetID"`
		UserID           uint                    `json:"user_id" gorm:"not null"`
	}
)
