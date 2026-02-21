package routes

import (
	"quiz-backend/internal/handlers"
	"quiz-backend/internal/middleware"
	"quiz-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// NewRouter returns a new Gin engine with recovery and logger middleware.
func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	return r
}

// Setup configures all routes: public auth and versioned API with RBAC.
func Setup(r *gin.Engine, h *Handlers, jwtSecret string) {
	// JWT secret must be available to auth middleware
	middleware.JWTSecret = []byte(jwtSecret)

	// Public auth
	r.POST("/auth/register", h.Auth.Register)
	r.POST("/auth/login", h.Auth.Login)

	// Versioned API
	v1 := r.Group("/api/v1")
	{
		// Teacher: require auth + role teacher
		teacher := v1.Group("/teacher")
		teacher.Use(middleware.AuthMiddleware(), middleware.RequireRole("teacher"))
		{
			teacher.POST("/quizzes", h.Teacher.CreateQuiz)
			teacher.POST("/quizzes/:quizId/questions", h.Teacher.AddQuestion)
			teacher.PUT("/questions/:questionId", h.Teacher.UpdateQuestion)
			teacher.DELETE("/questions/:questionId", h.Teacher.DeleteQuestion)
			teacher.GET("/quizzes", h.Teacher.ListQuizzes)
			teacher.GET("/quizzes/:quizId/submissions", h.Teacher.GetQuizSubmissions)
		}

		// Student: require auth + role student
		student := v1.Group("")
		student.Use(middleware.AuthMiddleware(), middleware.RequireRole("student"))
		{
			student.GET("/quizzes", h.Student.ListQuizzes)
			student.GET("/quizzes/:quizId", h.Student.GetQuiz)
			student.POST("/quizzes/:quizId/submit", h.Student.SubmitQuiz)
			student.GET("/student/submissions", h.Student.GetMySubmissions)
		}
	}

	// Health check for deployment
	r.GET("/health", func(c *gin.Context) {
		utils.JSONSuccess(c, gin.H{"status": "ok"})
	})
}

// Handlers holds all handler instances for dependency injection.
type Handlers struct {
	Auth    *handlers.AuthHandler
	Teacher *handlers.TeacherHandler
	Student *handlers.StudentHandler
}
