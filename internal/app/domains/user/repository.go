package user

import (
	"github.com/jho3r/finanger-back/internal/infrastructure/database/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerRepo = logger.Setup("domain.user.repository")

// Repository is the interface for the user repository.
type Repository interface {
	FindByEmail(email string) (User, error)
	Create(user User) error
}

// RepositoryImpl is the struct that contains the user repository.
type RepositoryImpl struct {
	db gorm.Gorm
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db gorm.Gorm) Repository {
	return &RepositoryImpl{db: db}
}

// FindByEmail finds a user by email.
func (r *RepositoryImpl) FindByEmail(email string) (User, error) {
	var user User
	if err := r.db.WhereFirst(&user, "email = ?", email); err != nil {
		loggerRepo.WithError(err).Error("Error querying the user by email")

		return User{}, err
	}

	return user, nil
}

// Create creates a new user.
func (r *RepositoryImpl) Create(user User) error {
	if err := r.db.Create(&user); err != nil {
		loggerRepo.WithError(err).Error("Error creating record in the database")

		return err
	}

	return nil
}
