package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jho3r/finanger-back/internal/app/controller"
	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/app/domains/user"
	"github.com/jho3r/finanger-back/internal/app/settings"
)

var (
	errUnauthorized = errors.New("unauthorized")
	errParseJWT     = errors.New("parse jwt error")
)

// AuthUser is the middleware that validates the JWT.
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controller.Error{Message: "Authorization header is required", Error: errUnauthorized.Error()})

			return
		}

		// Get the token from the header
		token := strings.Replace(authHeader, "Bearer ", "", 1)

		// Validate the token
		claims, err := validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controller.Error{Message: "Invalid token", Error: err.Error()})

			return
		}

		// Set the user ID in the context
		c.Set(settings.Commons.UserIDContextKey, claims.UserID)

		c.Next()
	}
}

// validateToken validates the token and returns the claims if the token is valid or an error otherwise.
func validateToken(tokenStr string) (*user.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &user.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.Auth.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf(crosscuting.WrapLabel, "Error parsing the token", errParseJWT, err.Error())
	}

	claims, ok := token.Claims.(*user.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf(crosscuting.WrapLabelWithoutError, "Error getting the claims", errParseJWT)
	}

	return claims, nil
}
