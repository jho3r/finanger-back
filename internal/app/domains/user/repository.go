package user

import (
	"errors"

	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerRepo = logger.Setup("domain.user.repository")

// Repository is the interface for the user repository.
type Repository interface {
	FindByEmail(email string) (User, error)
	Create(user User) error
	FindByID(id uint) (User, error)
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
	err := r.db.WhereFirst(&user, "email = ?", email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

func (r *RepositoryImpl) FindByID(id uint) (User, error) {
	var user User
	err := r.db.WhereFirst(&user, "id = ?", id)
	if err != nil {
		loggerRepo.WithError(err).Error("Error querying the user by id")

		return User{}, err
	}

	return user, nil
}
