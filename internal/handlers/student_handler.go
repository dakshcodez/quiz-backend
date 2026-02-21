package handlers

import (
	"context"
	"time"

	"quiz-backend/internal/models"
	"quiz-backend/internal/services"
	"quiz-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// StudentHandler handles student-only endpoints; uses QuizService and RBAC.
type StudentHandler struct {
	quizService *services.QuizService
	timeout     time.Duration
}

// NewStudentHandler returns a new StudentHandler.
func NewStudentHandler(quizService *services.QuizService, timeout time.Duration) *StudentHandler {
	return &StudentHandler{quizService: quizService, timeout: timeout}
}

func (h *StudentHandler) userID(c *gin.Context) string {
	id, _ := c.Get(string(utils.ContextKeyUserID))
	s, _ := id.(string)
	return s
}

func (h *StudentHandler) withTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), h.timeout)
}

// ListQuizzes handles GET /api/v1/quizzes.
func (h *StudentHandler) ListQuizzes(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	list, err := h.quizService.ListQuizzes(ctx)
	if err != nil {
		utils.JSONInternal(c, "failed to list quizzes")
		return
	}
	utils.JSONSuccess(c, list)
}

// GetQuiz handles GET /api/v1/quizzes/:quizId. Returns quiz without correct answers.
func (h *StudentHandler) GetQuiz(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	quizID := c.Param("quizId")
	if quizID == "" {
		utils.JSONBadRequest(c, "quizId required")
		return
	}
	view, err := h.quizService.GetQuizForStudent(ctx, quizID)
	if err != nil {
		if err == services.ErrQuizNotFound {
			utils.JSONNotFound(c, "quiz not found")
			return
		}
		utils.JSONInternal(c, "failed to get quiz")
		return
	}
	utils.JSONSuccess(c, view)
}

// SubmitQuiz handles POST /api/v1/quizzes/:quizId/submit.
func (h *StudentHandler) SubmitQuiz(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	studentID := h.userID(c)
	quizID := c.Param("quizId")
	if quizID == "" {
		utils.JSONBadRequest(c, "quizId required")
		return
	}
	var req models.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, err.Error())
		return
	}
	result, err := h.quizService.SubmitQuiz(ctx, studentID, quizID, req.Answers)
	if err != nil {
		if err == services.ErrQuizNotFound {
			utils.JSONNotFound(c, "quiz not found")
			return
		}
		utils.JSONInternal(c, "failed to submit quiz")
		return
	}
	utils.JSONSuccess(c, result)
}

// GetMySubmissions handles GET /api/v1/student/submissions.
func (h *StudentHandler) GetMySubmissions(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	studentID := h.userID(c)
	list, err := h.quizService.GetStudentSubmissions(ctx, studentID)
	if err != nil {
		utils.JSONInternal(c, "failed to get submissions")
		return
	}
	utils.JSONSuccess(c, list)
}
