package user

import "github.com/jho3r/finanger-back/internal/infrastructure/database/gorm"

type (
	// User is the struct for the user.
	User struct {
		gorm.Model
		Name     string `json:"name" gorm:"not null"`
		Email    string `json:"email" gorm:"not null;unique"`
		Currency string `json:"currency" gorm:"not null"`
		Password string `json:"password" gorm:"not null"`
	}
)
