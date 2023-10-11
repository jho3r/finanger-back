package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/crosscuting"
	"github.com/jho3r/finanger-back/internal/app/domains/asset"
	"github.com/jho3r/finanger-back/internal/infrastructure/gorm"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var (
	loggerAssets      = logger.Setup("controller.assets")
	errGettingAssetID = errors.New("Error getting the asset ID from the context")
)

// CreateAsset is the controller for the create asset endpoint.
func CreateAsset(assetService asset.Service) gin.HandlerFunc {
	type createRequest struct {
		Name        string  `json:"name" binding:"required"`
		Value       float64 `json:"value" binding:"required"`
		FinAssetID  uint    `json:"fin_asset_id" binding:"required"`
		Description string  `json:"description" binding:"required"`
	}

	return func(c *gin.Context) {
		var request createRequest

		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			desc := "There is something wrong with the data sent"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
		}

		asset := asset.Asset{
			Name:             request.Name,
			Value:            request.Value,
			Description:      request.Description,
			FinancialAssetID: request.FinAssetID,
			UserID:           userID,
		}

		err = assetService.Create(asset)
		if err != nil {
			desc := "Error creating the asset"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Success{Message: "Asset created successfully"})
	}
}

// GetAssets is the controller for the get assets endpoint.
func GetAssets(assetService asset.Service) gin.HandlerFunc {
	type getAssetsQuery struct {
		FinAssetID uint   `form:"fin_asset_id"`
		Name       string `form:"name"`
	}

	return func(c *gin.Context) {
		var query getAssetsQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			desc := "Error binding the query"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
			return
		}

		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		assets, err := assetService.GetAssets(userID, query.FinAssetID, query.Name)
		if err != nil {
			desc := "Error getting the assets"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"assets": assets}})
	}
}

// GetAsset is the controller for the get asset endpoint.
func GetAsset(assetService asset.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetID, err := getIdFromParams(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		asset, err := assetService.GetAssetByID(userID, assetID)
		if err != nil {
			desc := "Error getting the asset"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"asset": asset}})
	}
}

// UpdateAsset is the controller for the update asset endpoint.
func UpdateAsset(assetService asset.Service) gin.HandlerFunc {
	type updateRequest struct {
		Name        string  `json:"name"`
		Value       float64 `json:"value"`
		FinAssetID  uint    `json:"fin_asset_id"`
		Description string  `json:"description"`
	}

	return func(c *gin.Context) {
		var request updateRequest

		assetID, err := getIdFromParams(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			desc := "There is something wrong with the data sent"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
		}

		asset := asset.Asset{
			Model:            gorm.Model{ID: assetID},
			Name:             request.Name,
			Value:            request.Value,
			Description:      request.Description,
			FinancialAssetID: request.FinAssetID,
			UserID:           userID,
		}

		err = assetService.Update(asset)
		if err != nil {
			desc := "Error updating the asset"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Success{Message: "Asset updated successfully"})
	}
}

// DeleteAsset is the controller for the delete asset endpoint.
func DeleteAsset(assetService asset.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		assetID, err := getIdFromParams(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		userID, err := getUserIdFromContext(c)
		if err != nil {
			loggerAssets.Error(err.Error())
			c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
			return
		}

		err = assetService.Delete(userID, assetID)
		if err != nil {
			desc := "Error deleting the asset"
			loggerAssets.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Success{Message: "Asset deleted successfully"})
	}
}

func getIdFromParams(c *gin.Context) (uint, error) {
	assetIDStr, found := c.Params.Get("id")
	if !found {
		return 0, errGettingAssetID
	}

	assetID, err := strconv.ParseUint(assetIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf(crosscuting.WrapLabel, "Error parsing the asset ID", errGettingAssetID, err.Error())
	}

	return uint(assetID), nil
}
