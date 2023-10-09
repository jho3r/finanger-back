package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jho3r/finanger-back/internal/app/domains/category"
	"github.com/jho3r/finanger-back/internal/infrastructure/logger"
)

var loggerCategory = logger.Setup("controller.category")

// CreateCategory is the controller for the create category endpoint.
func CreateCategory(categoryService category.Service) gin.HandlerFunc {

	type createCategoryReq struct {
		Name string `json:"name" binding:"required"`
		Desc string `json:"desc" binding:"required"`
		Type string `json:"type" binding:"required,oneof=asset liability income expense"`
	}

	return func(c *gin.Context) {
		var request createCategoryReq
		if err := c.ShouldBindJSON(&request); err != nil {
			desc := "There is something wrong with the data sent"
			loggerCategory.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})

			return
		}

		category := category.Category{
			Name:        request.Name,
			Description: request.Desc,
			Type:        category.CategoryType(request.Type),
		}

		if err := categoryService.CreateCategory(category); err != nil {
			desc := "Error creating the category"
			loggerCategory.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusCreated, Success{Message: "Category created successfully"})
	}
}

// GetCategories is the controller for the get categories endpoint.
func GetCategories(categoryService category.Service) gin.HandlerFunc {

	type getCategoriesQuery struct {
		Type string `form:"type" binding:"omitempty,oneof=asset liability income expense"`
		Name string `form:"name"`
	}

	return func(c *gin.Context) {
		var query getCategoriesQuery
		if err := c.ShouldBindQuery(&query); err != nil {
			desc := "There is something wrong with the data sent"
			loggerCategory.WithError(err).Error(desc)
			c.JSON(http.StatusBadRequest, Error{Message: desc, Error: err.Error()})
			return
		}

		categories, err := categoryService.GetCategories(category.CategoryType(query.Type), query.Name)
		if err != nil {
			desc := "Error getting the categories"
			loggerCategory.WithError(err).Error(desc)
			c.JSON(http.StatusInternalServerError, Error{Message: desc, Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, Data{Data: gin.H{"categories": categories}})
	}
}
