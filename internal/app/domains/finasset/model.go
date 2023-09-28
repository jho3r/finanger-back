package finasset

import "github.com/jho3r/finanger-back/internal/infrastructure/database/gorm"

const (
	// Currency is the type for the currency.
	Currency AssetType = "currency"
	// Stock is the type for the stock.
	Stock AssetType = "stock"
	// Crypto is the type for the crypto.
	Crypto AssetType = "crypto"
)

type (
	// AssetType can be currency, stock, crypto, etc.
	AssetType string

	// Currency is the struct for the currency.
	FinancialAsset struct {
		gorm.Model
		Symbol string    `json:"symbol" gorm:"not null;unique"`
		Name   string    `json:"name" gorm:"not null"`
		Desc   string    `json:"desc" gorm:"not null;column:description"`
		Type   AssetType `json:"type" gorm:"not null"`
	}
)
