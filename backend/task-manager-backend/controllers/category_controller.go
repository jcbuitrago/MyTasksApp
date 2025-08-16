package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"task-manager-backend/models"
	"task-manager-backend/services"
)

type CategoryController struct {
	service services.CategoryService
}

func NewCategoryController(service services.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := c.service.CreateCategory(&category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	ctx.JSON(http.StatusCreated, category)
}

func (c *CategoryController) GetCategories(ctx *gin.Context) {
	categories, err := c.service.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	category, err := c.service.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	var category models.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := c.service.UpdateCategory(id, &category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.service.DeleteCategory(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}