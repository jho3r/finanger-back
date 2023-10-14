package asset

import (
	"time"

	"github.com/jho3r/finanger-back/internal/app/domains/category"
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
		IsLiquid         bool                    `json:"is_liquid" gorm:"not null"`
		CategoryID       uint                    `json:"category_id" gorm:"not null"`
		Category         category.Category       `json:"category" gorm:"foreignKey:CategoryID"`
		AdquisitionDate  time.Time               `json:"adquisition_date" gorm:"not null"`
	}
)
