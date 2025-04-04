package main

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/ouzrourextra/official-ouzrour/config"
	"github.com/ouzrourextra/official-ouzrour/handlers"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Middleware : fichiers statiques
	r.Static("/static", "./static")

	// HTMLRender : chargé avec Gin-compatible multitemplate
	r.HTMLRender = config.LoadTemplates()

	// Session sécurisée
	store := cookie.NewStore([]byte("Ime8hm9HPmihiMTRYwFv7xJIaAi47oGTw79Q+mNjtXBvNVrkCFOkpq4aI5EZ2l3zvMi4gMqC0yPxBzWwLfjvlQ=="))
	store.Options(sessions.Options{
		HttpOnly: true,
		Secure:   false, // true en production
		Path:     "/",
		SameSite: 2, // Strict
	})
	r.Use(sessions.Sessions("session", store))

	// Routes publiques
	r.GET("/", handlers.Home)

	// Auth
	r.GET("/login", handlers.LoginPage)
	r.POST("/login", handlers.Login)
	r.GET("/logout", handlers.Logout)

	// Zone protégée
	protected := r.Group("/tasks")
	protected.Use(handlers.RequireAuth())
	{
		protected.GET("/", handlers.TasksIndex)
		protected.POST("/add", handlers.TaskAdd)
		protected.POST("/update/:id", handlers.TaskUpdate)
		protected.POST("/delete/:id", handlers.TaskDelete)
	}

	log.Println("✅ Serveur lancé sur : http://localhost:8080")
	r.Run("localhost:8080")
}
