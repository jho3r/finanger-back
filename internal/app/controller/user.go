package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/domains/user"
	"github.com/jho3r/finanger-back/internal/app/settings"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

const (
	errBindingUser = "Error binding the user"
)

var loggerUser = logger.Setup("controller.user")

// Login is the controller for the login endpoint.
func Signup(userService user.Service) gin.HandlerFunc {

	type signupRequest struct {
		Name       string `json:"name" binding:"required"`
		Email      string `json:"email" binding:"required"`
		FinAssetID int    `json:"fin_asset_id" binding:"required"`
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

		c.SetCookie("refresh_token", refreshToken, 60*60*24*7, "/", "", false, true)

		c.JSON(http.StatusOK, Data{Data: gin.H{"token": token}})
	}
}

// RefreshToken is the controller for the refresh token endpoint.
func RefreshToken(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			desc := "Error getting the refresh token from the cookie"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
			return
		}

		token, err := userService.RefreshToken(refreshToken)
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
		userID, exist := c.Get(settings.Commons.UserIDContextKey)
		if !exist {
			desc := "Error getting the user ID from the context"
			loggerUser.Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc})
			return
		}

		userID, ok := userID.(uint)
		if !ok {
			desc := "Error getting the user ID from the context"
			loggerUser.Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc})
			return
		}

		user, err := userService.GetMe(userID.(uint))
		if err != nil {
			desc := "Error getting the user"
			loggerUser.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"user": user}})
	}
}
