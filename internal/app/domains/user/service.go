package user

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
	"golang.org/x/crypto/bcrypt"
)

var (
	loggerService   = logger.Setup("domain.user.service")
	errValidateUser = errors.New("validate user error")
	errParseJWT     = errors.New("parse jwt error")
	errValidateJWT  = errors.New("validate jwt error")
)

// UserService is the interface for the user service.
type Service interface {
	Signup(user User) error
	Login(user User) (string, string, error)
	RefreshToken(userID uint, refreshToken string) (string, error)
	GetMe(userID uint) (User, error)
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

	err = user.HashPassword()
	if err != nil {
		return err
	}

	if err := s.repo.Create(user); err != nil {
		loggerService.WithError(err).Error("Error creating the user")

		return err
	}

	return nil
}

// Login logs in a user, returns a token if the user is valid or an error otherwise.
func (s *ServiceImpl) Login(user User) (string, string, error) {
	existingUser, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		loggerService.WithError(err).Error("Error finding the user by email")

		return "", "", err
	}

	if existingUser.Email == "" {
		desc := "User not found"
		loggerService.WithError(errValidateUser).Error(desc)

		return "", "", fmt.Errorf(crosscuting.WrapLabelWithoutError, desc, errValidateUser)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		desc := "Invalid password"
		loggerService.WithError(errValidateUser).Error(desc)

		return "", "", fmt.Errorf(crosscuting.WrapLabelWithoutError, desc, errValidateUser)
	}

	token, err := existingUser.GenerateToken()
	if err != nil {
		return "", "", err
	}

	refreshToken, err := existingUser.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

// RefreshToken refreshes the access token.
func (s *ServiceImpl) RefreshToken(userID uint, refreshToken string) (string, error) {
	claims, err := validateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	if claims.UserID != userID {
		desc := "Invalid user ID"
		loggerService.WithError(errValidateJWT).Error(desc)

		return "", fmt.Errorf(crosscuting.WrapLabelWithoutError, desc, errValidateJWT)
	}

	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		loggerService.WithError(err).Error("Error finding the user by ID")

		return "", err
	}

	token, err := user.GenerateToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetMe returns the user with the given ID.
func (s *ServiceImpl) GetMe(userID uint) (User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		loggerService.WithError(err).Error("Error finding the user by ID")

		return User{}, err
	}

	return user, nil
}

// validateRefreshToken validates the refresh token and returns the claims if the token is valid or an error otherwise.
func validateRefreshToken(refreshToken string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.Auth.JWTRefreshSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf(crosscuting.WrapLabel, "Error parsing the token", errParseJWT, err.Error())
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf(crosscuting.WrapLabelWithoutError, "Error validating the token", errValidateJWT)
	}

	return claims, nil
}
