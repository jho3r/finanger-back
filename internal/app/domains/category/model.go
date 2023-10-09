package category

import "github.com/jho3r/finanger-back/internal/infrastructure/database/gorm"

const (
	// Asset is the type for the asset.
	Asset CategoryType = "asset"
	// Liability is the type for the liability.
	Liability CategoryType = "liability"
	// Income is the type for the income.
	Income CategoryType = "income"
	// Expense is the type for the expense.
	Expense CategoryType = "expense"
)

type (
	// CategoryType can be asset, liability, income, expense, etc.
	CategoryType string

	// Category represents a category.
	Category struct {
		gorm.Model
		Name        string       `gorm:"not null"`
		Description string       `gorm:"not null"`
		Type        CategoryType `gorm:"not null"`
	}
)
