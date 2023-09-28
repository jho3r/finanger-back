package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/domains/user"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerUser = logger.Setup("controller.user")

type User struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login is the controller for the login endpoint
func Signup(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request User
		if err := c.ShouldBindJSON(&request); err != nil {
			loggerUser.WithError(err).Error("Error binding the user")
			c.JSON(http.StatusBadRequest, Error{Message: "Error binding the user", Error: err.Error()})
			return
		}

		user := user.User{
			Name:     request.Name,
			Email:    request.Email,
			Currency: request.Currency,
			Password: request.Password,
		}

		if err := userService.Signup(user); err != nil {
			loggerUser.WithError(err).Error("Error signing up the user")
			c.JSON(http.StatusInternalServerError, Error{Message: "Error signing up the user", Error: err.Error()})
			return
		}

		c.JSON(200, Success{Message: "User created successfully"})
	}
}
