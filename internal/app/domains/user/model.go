package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/app/domains/finasset"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenExp  = 30 * time.Minute
	refreshTokenExp = 7 * 24 * time.Hour
	hashCost        = 14
)

var (
	loggerModel    = logger.Setup("domain.user.model")
	errHash        = errors.New("hash error")
	errGenerateJWT = errors.New("generate jwt error")
)

type (
	// User is the struct for the user.
	User struct {
		gorm.Model
		Name             string                  `json:"name" gorm:"not null"`
		Email            string                  `json:"email" gorm:"not null;unique"`
		Password         string                  `json:"-" gorm:"not null"`
		FinancialAssetID uint                    `json:"fin_asset_id" gorm:"not null"`
		FinancialAsset   finasset.FinancialAsset `json:"financial_asset" gorm:"foreignKey:FinancialAssetID"`
	}

	// Claims is the struct that contains the claims of the JWT.
	Claims struct {
		UserID uint   `json:"user_id"`
		Email  string `json:"email"`
		jwt.RegisteredClaims
	}

	RefreshClaims struct {
		UserID uint `json:"user_id"`
		jwt.RegisteredClaims
	}
)

// hashPassword hashes the user password.
func (u *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), hashCost)
	if err != nil {
		desc := "Error hashing the password"
		loggerModel.WithError(err).Error(desc)

		return fmt.Errorf(crosscuting.WrapLabel, desc, errHash, err.Error())
	}

	u.Password = string(bytes)

	return nil
}

// generateToken generates a JWT access token for the user.
func (u User) GenerateToken() (string, error) {
	claims := Claims{
		u.ID,
		u.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(settings.Auth.JWTSecret))
	if err != nil {
		return "", fmt.Errorf(crosscuting.WrapLabel, "Error generating the token", errGenerateJWT, err.Error())
	}

	return tokenString, nil
}

// generateRefreshToken generates a JWT refresh token for the user.
func (u User) GenerateRefreshToken() (string, error) {
	claims := RefreshClaims{
		u.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(settings.Auth.JWTRefreshSecret))
	if err != nil {
		return "", fmt.Errorf(crosscuting.WrapLabel, "Error generating the token", errGenerateJWT, err.Error())
	}

	return tokenString, nil
}
