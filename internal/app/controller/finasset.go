package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/domains/finasset"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerFinAsset = logger.Setup("controller.finasset")

type FinancialAssetReq struct {
	Symbol string `json:"symbol" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Desc   string `json:"desc" binding:"required"`
	Type   string `json:"type" binding:"required,oneof=currency stock crypto"`
}

type GetFinancialAssetsQuery struct {
	Type   string `form:"type" binding:"omitempty,oneof=currency stock crypto"`
	Symbol string `form:"symbol"`
	Name   string `form:"name"`
}

// CreateFinancialAsset creates a new financial asset.
func CreateFinancialAsset(finAssetService finasset.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request FinancialAssetReq
		if err := c.ShouldBindJSON(&request); err != nil {
			loggerFinAsset.WithError(err).Error("Error binding the financial asset")
			c.JSON(http.StatusBadRequest, Error{Message: "Error binding the financial asset", Error: err.Error()})
			return
		}

		finAsset := finasset.FinancialAsset{
			Symbol: request.Symbol,
			Name:   request.Name,
			Desc:   request.Desc,
			Type:   finasset.AssetType(request.Type),
		}

		if err := finAssetService.Create(finAsset); err != nil {
			loggerFinAsset.WithError(err).Error("Error creating the financial asset")
			c.JSON(http.StatusInternalServerError, Error{Message: "Error creating the financial asset", Error: err.Error()})
			return
		}

		c.JSON(200, Success{Message: "Financial asset created successfully"})
	}
}

// GetFinancialAssets returns all the financial assets given filters.
func GetFinancialAssets(finAssetService finasset.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query GetFinancialAssetsQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			loggerFinAsset.WithError(err).Error("Error binding the query")
			c.JSON(http.StatusBadRequest, Error{Message: "Error binding the query", Error: err.Error()})
			return
		}

		finAssets, err := finAssetService.Get(finasset.FinancialAsset{
			Symbol: query.Symbol,
			Name:   query.Name,
			Type:   finasset.AssetType(query.Type),
		})
		if err != nil {
			loggerFinAsset.WithError(err).Error("Error getting the financial assets")
			c.JSON(http.StatusInternalServerError, Error{Message: "Error getting the financial assets", Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"financial_assets": finAssets}})
	}
}
