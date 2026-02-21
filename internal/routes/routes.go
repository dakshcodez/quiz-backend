package routes

import (
	"github.com/gin-gonic/gin"
	"quiz-backend/internal/handlers"
	"quiz-backend/internal/store"
)

// SetupRoutes registers all HTTP routes. Pass the in-memory store so handlers can read data.
func SetupRoutes(r *gin.Engine, st *store.MemoryStore) {
	studentHandler := handlers.NewStudentHandler(st)

	student := r.Group("/student")
	{
		student.GET("/view_quiz/:quiz_id", studentHandler.ViewQuiz)
		student.POST("/give_quiz", studentHandler.GiveQuiz)
	}
}