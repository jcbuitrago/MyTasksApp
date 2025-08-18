package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"mytasks/internal/db"
	"mytasks/internal/handlers"
	"mytasks/internal/middleware"
	"mytasks/internal/repository"
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

	// ----- Inyección de dependencias -----
	// Users (como ya lo tenías)
	uh := &handlers.UsersHandler{DB: database}

	// Categories (si tu compañero usa repo, cambia por NewCategoryRepository/NewCategoriesHandler)
	ch := &handlers.CategoriesHandler{DB: database}

	// Tasks (usa el repo que te pasé)
	taskRepo := repository.NewTaskRepository(database)
	th := handlers.NewTaskHandler(taskRepo)

	// ----- Router -----
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	api := r.Group("/")

	// Usuarios (público)
	api.POST("/usuarios", uh.Register)
	api.POST("/usuarios/iniciar-sesion", uh.Login)

	// Healthcheck
	api.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// Rutas autenticadas
	auth := api.Group("/")
	auth.Use(middleware.AuthRequired()) // tu middleware existente

	// Categorías
	auth.POST("/categorias", ch.Create)
	auth.GET("/categorias", ch.List)
	auth.DELETE("/categorias/:id", ch.Delete)

	// Tareas
	auth.POST("/tareas", th.Create)
	auth.PUT("/tareas/:id", th.Update)
	auth.DELETE("/tareas/:id", th.Delete)
	auth.GET("/tareas/usuario", th.ListByUser) // ?categoria_id=&estado=&q=&page=&page_size=
	auth.GET("/tareas/:id", th.GetByID)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
