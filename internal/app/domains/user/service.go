package user

import (
	"errors"
	"fmt"

	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
	"golang.org/x/crypto/bcrypt"
)

var (
	loggerService   = logger.Setup("domain.user.service")
	errHash         = errors.New("hash error")
	errValidateUser = errors.New("validate user error")
)

// UserService is the interface for the user service.
type Service interface {
	Signup(user User) error
}

// ServiceImpl is the struct that contains the user service.
type ServiceImpl struct {
	repo Repository
}

// NewUserService creates a new user service.
func NewUserService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}

// Signup creates a new user.
func (s *ServiceImpl) Signup(user User) error {
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		loggerService.WithError(err).Error("Error finding the user by email")

		return err
	}

	if existingUser.Email != "" {
		desc := "User already exists"
		loggerService.WithError(errValidateUser).Error(desc)

		return fmt.Errorf(crosscuting.WrapLabelWithoutError, desc, errValidateUser)
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := s.repo.Create(user); err != nil {
		loggerService.WithError(err).Error("Error creating the user")

		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		desc := "Error hashing the password"
		loggerService.WithError(err).Error(desc)

		return "", fmt.Errorf(crosscuting.WrapLabel, desc, errHash, err.Error())
	}

	return string(bytes), nil
}
