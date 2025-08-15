package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mytasks/internal/db"
	"mytasks/internal/handlers"
)

func main() {
	// Abrir DB
	database, err := db.Open()
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer database.Close()

	// Pingar para validar
	if err := database.Ping(); err != nil {
		log.Fatalf("DB ping error: %v", err)
	}

	r := gin.Default()

	// CORS abierto para desarrollo
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	uh := &handlers.UsersHandler{DB: database}

	// Endpoints Rol A
	r.POST("/usuarios", uh.Register)             // Crear usuario
	r.POST("/usuarios/iniciar-sesion", uh.Login) // Login -> token

	// Healthcheck opcional
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
