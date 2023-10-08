package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/app/domains/user"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

const (
	errBindingUser = "Error binding the user"
)

var (
	loggerUser       = logger.Setup("controller.user")
	errGettingUserID = errors.New("Error getting the user ID from the context")
)

// Login is the controller for the login endpoint.
func Signup(userService user.Service) gin.HandlerFunc {

	type signupRequest struct {
		Name       string `json:"name" binding:"required"`
		Email      string `json:"email" binding:"required"`
		FinAssetID uint   `json:"fin_asset_id" binding:"required"`
		Password   string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var request signupRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			loggerUser.WithError(err).Error(errBindingUser)
			c.JSON(http.StatusBadRequest, Error{Message: errBindingUser, Error: err.Error()})
			return
		}

		user := user.User{
			Name:             request.Name,
			Email:            request.Email,
			FinancialAssetID: request.FinAssetID,
			Password:         request.Password,
		}

		if err := userService.Signup(user); err != nil {
			desc := "Error signing up the user"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Success{Message: "User created successfully"})
	}
}

// Login is the controller for the login endpoint.
func Login(userService user.Service) gin.HandlerFunc {

	type loginRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	return func(c *gin.Context) {
		var request loginRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			loggerUser.WithError(err).Error(errBindingUser)
			c.JSON(http.StatusBadRequest, Error{Message: errBindingUser, Error: err.Error()})
			return
		}

		user := user.User{
			Email:    request.Email,
			Password: request.Password,
		}

		token, refreshToken, err := userService.Login(user)
		if err != nil {
			desc := "Error logging in the user"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.SetCookie(settings.Auth.RefreshTokenCookieName, refreshToken, 60*60*24*7, "/", "", false, true)

		c.JSON(http.StatusOK, Data{Data: gin.H{"token": token}})
	}
}

// RefreshToken is the controller for the refresh token endpoint.
func RefreshToken(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerUser.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		refreshToken, err := c.Cookie(settings.Auth.RefreshTokenCookieName)
		if err != nil {
			desc := "Error getting the refresh token from the cookie"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
			return
		}

		token, err := userService.RefreshToken(userID, refreshToken)
		if err != nil {
			desc := "Error refreshing the token"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"token": token}})
	}
}

// GetMe is the controller for the get me endpoint.
func GetMe(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerUser.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		user, err := userService.GetMe(userID)
		if err != nil {
			desc := "Error getting the user"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"user": user}})
	}
}

// Logout is the controller for the logout endpoint.
func Logout(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// refreshToken, err := c.Cookie(settings.Auth.RefreshTokenCookieName)
		// if err != nil {
		// 	desc := "Error getting the refresh token from the cookie"
		// 	loggerUser.WithError(err).Error(desc)
		// 	c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
		// 	return
		// }

		// if err := userService.Logout(refreshToken); err != nil {
		// 	desc := "Error logging out the user"
		// 	loggerUser.WithError(err).Error(desc)
		// 	c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
		// 	return
		// }

		c.SetCookie(settings.Auth.RefreshTokenCookieName, "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, Success{Message: "User logged out successfully"})
	}
}

func getUserIdFromContext(c *gin.Context) (uint, error) {
	desc := "Error getting the user ID from the context"

	userIDAny, exists := c.Get(settings.Auth.UserIDContextKey)
	if !exists {
		return 0, fmt.Errorf(crosscuting.WrapLabelWithoutError, desc, errGettingUserID)
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		return 0, fmt.Errorf(crosscuting.WrapLabelWithoutError, desc, errGettingUserID)
	}

	return userID, nil
}
