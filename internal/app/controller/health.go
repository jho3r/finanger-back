package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/settings"
)

type HealthStatus struct {
	Status string    `json:"status"`
	Name   string    `json:"name"`
	AppID  string    `json:"app_id"`
	Date   time.Time `json:"date"`
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthStatus{
		Status: "OK",
		Name:   settings.Commons.ProjectName,
		AppID:  settings.Commons.XApplicationID,
		Date:   time.Now(),
	})
}
