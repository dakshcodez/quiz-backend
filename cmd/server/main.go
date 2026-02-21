package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"quiz-backend/internal/routes"
	"quiz-backend/internal/store"
)

func main() {
	// In-memory store: data resets on every server restart. No database.
	memStore := store.NewMemoryStore()

	r := gin.Default()
	routes.SetupRoutes(r, memStore)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}