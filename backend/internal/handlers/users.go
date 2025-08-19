package handlers

import (
	"context"
	"database/sql"
	"mytasks/internal/middleware"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	DB *sql.DB
}

type registerReq struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	PictureURL string `json:"picture_url"`
}

type userResp struct {
	ID         int64  `json:"id"`
	Username   string `json:"username"`
	PictureURL string `json:"picture_url,omitempty"`
}

func (h *UsersHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	if len(req.Username) < 3 || len(req.Password) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "usuario o contraseña muy cortos"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo encriptar"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var id int64
	err = h.DB.QueryRowContext(ctx,
		`INSERT INTO users (username, password_hash, picture_url)
		 VALUES ($1,$2,$3) RETURNING id`,
		req.Username, string(hash), req.PictureURL).Scan(&id)

	if err != nil {
		// Duplicado (unique)
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "usuario ya existe"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusCreated, userResp{ID: id, Username: req.Username, PictureURL: req.PictureURL})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UsersHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var id int64
	var hash string
	err := h.DB.QueryRowContext(ctx,
		`SELECT id, password_hash FROM users WHERE username=$1`, req.Username).
		Scan(&id, &hash)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	tok, err := middleware.GenerateToken(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tok, "user_name": req.Username})
}
