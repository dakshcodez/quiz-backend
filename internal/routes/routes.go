package routes

import (
	"github.com/gin-gonic/gin"
	"quiz-backend/internal/handlers"
	"quiz-backend/internal/store"
)

// SetupRoutes registers all HTTP routes. Pass the in-memory store so handlers can read data.
func SetupRoutes(r *gin.Engine, st *store.MemoryStore) {
	studentHandler := handlers.NewStudentHandler(st)
	teacherHandler := handlers.NewTeacherHandler(st)

	student := r.Group("/student")
	{
		student.GET("/view_quiz/:quiz_id", studentHandler.ViewQuiz)
		student.POST("/give_quiz", studentHandler.GiveQuiz)
	}

	teacher := r.Group("/teacher")
	{
		teacher.POST("/create_quiz", teacherHandler.CreateQuiz)
		teacher.POST("/add_question/:quiz_id", teacherHandler.AddQuestion)
		teacher.PUT("/update_question/:question_id", teacherHandler.UpdateQuestion)
		teacher.DELETE("/delete_question/:question_id", teacherHandler.DeleteQuestion)
		teacher.GET("/view_quiz/:quiz_id", teacherHandler.ViewQuiz)
		teacher.GET("/all_quizzes", teacherHandler.AllQuizzes)
	}
}