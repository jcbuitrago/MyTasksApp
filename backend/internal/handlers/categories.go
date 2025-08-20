// backend/internal/handlers/categories.go
package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoriesHandler struct {
	DB *sql.DB
}

type categoryReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type categoryResp struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// POST /categorias
func (h *CategoriesHandler) Create(c *gin.Context) {
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name muy corto (min 3)"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	// Verificación simple de duplicado (no atómica, pero suficiente para el proyecto)
	var existingID int64
	err := h.DB.QueryRowContext(ctx,
		`SELECT id FROM categories WHERE name=$1 LIMIT 1`, req.Name).
		Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "la categoría ya existe"})
		return
	} else if err != sql.ErrNoRows && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	var id int64
	err = h.DB.QueryRowContext(ctx,
		`INSERT INTO categories (name, description) VALUES ($1,$2) RETURNING id`,
		req.Name, req.Description).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusCreated, categoryResp{ID: id, Name: req.Name, Description: req.Description})
}

// GET /categorias
func (h *CategoriesHandler) List(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	rows, err := h.DB.QueryContext(ctx,
		`SELECT id, name, description FROM categories ORDER BY id ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	resp := make([]categoryResp, 0)
	for rows.Next() {
		var r categoryResp
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			continue
		}
		resp = append(resp, r)
	}
	c.JSON(http.StatusOK, resp)
}

// DELETE /categorias/:id
func (h *CategoriesHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id inválido"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	res, err := h.DB.ExecContext(ctx, `DELETE FROM categories WHERE id=$1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no existe"})
		return
	}
	c.Status(http.StatusNoContent)
}
