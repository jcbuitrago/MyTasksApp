r := gin.Default()
r.Use(cors.Default())

api := r.Group("/")

api.POST("/usuarios", usersHandler.Register)
api.POST("/usuarios/iniciar-sesion", usersHandler.Login)

auth := api.Group("/")
auth.Use(middleware.AuthRequired()) // requiere JWT

// Categorías
auth.POST("/categorias", categoriesHandler.Create)
auth.GET("/categorias", categoriesHandler.List)
auth.DELETE("/categorias/:id", categoriesHandler.Delete)

// Tareas
auth.POST("/tareas", tasksHandler.Create)
auth.PUT("/tareas/:id", tasksHandler.Update)
auth.DELETE("/tareas/:id", tasksHandler.Delete)
auth.GET("/tareas/usuario", tasksHandler.ListByUser)         // ?categoria_id=&estado=
auth.GET("/tareas/:id", tasksHandler.GetByID)

r.Run(":8080")
