package handlers

import (
	"context"
	"time"

	"quiz-backend/internal/models"
	"quiz-backend/internal/services"
	"quiz-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// TeacherHandler handles teacher-only endpoints; uses QuizService and RBAC.
type TeacherHandler struct {
	quizService *services.QuizService
	timeout     time.Duration
}

// NewTeacherHandler returns a new TeacherHandler.
func NewTeacherHandler(quizService *services.QuizService, timeout time.Duration) *TeacherHandler {
	return &TeacherHandler{quizService: quizService, timeout: timeout}
}

func (h *TeacherHandler) userID(c *gin.Context) string {
	id, _ := c.Get(string(utils.ContextKeyUserID))
	s, _ := id.(string)
	return s
}

func (h *TeacherHandler) withTimeout(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), h.timeout)
}

// CreateQuiz handles POST /api/v1/teacher/quizzes.
func (h *TeacherHandler) CreateQuiz(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	teacherID := h.userID(c)
	var req models.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, err.Error())
		return
	}
	quiz, err := h.quizService.CreateQuiz(ctx, teacherID, req.Title, req.Description)
	if err != nil {
		utils.JSONInternal(c, "failed to create quiz")
		return
	}
	utils.JSONCreated(c, quiz)
}

// AddQuestion handles POST /api/v1/teacher/quizzes/:quizId/questions.
func (h *TeacherHandler) AddQuestion(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	teacherID := h.userID(c)
	quizID := c.Param("quizId")
	if quizID == "" {
		utils.JSONBadRequest(c, "quizId required")
		return
	}
	var req models.AddQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, err.Error())
		return
	}
	q, err := h.quizService.AddQuestion(ctx, teacherID, quizID, &req)
	if err != nil {
		switch err {
		case services.ErrQuizNotFound:
			utils.JSONNotFound(c, "quiz not found")
			return
		case services.ErrForbidden:
			utils.JSONForbidden(c, "forbidden")
			return
		default:
			utils.JSONInternal(c, "failed to add question")
			return
		}
	}
	utils.JSONCreated(c, q)
}

// UpdateQuestion handles PUT /api/v1/teacher/questions/:questionId.
func (h *TeacherHandler) UpdateQuestion(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	teacherID := h.userID(c)
	questionID := c.Param("questionId")
	if questionID == "" {
		utils.JSONBadRequest(c, "questionId required")
		return
	}
	var req models.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, err.Error())
		return
	}
	err := h.quizService.UpdateQuestion(ctx, teacherID, questionID, &req)
	if err != nil {
		switch err {
		case services.ErrQuestionNotFound:
			utils.JSONNotFound(c, "question not found")
			return
		case services.ErrForbidden:
			utils.JSONForbidden(c, "forbidden")
			return
		default:
			utils.JSONInternal(c, "failed to update question")
			return
		}
	}
	utils.JSONSuccess(c, gin.H{"message": "updated"})
}

// DeleteQuestion handles DELETE /api/v1/teacher/questions/:questionId.
func (h *TeacherHandler) DeleteQuestion(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	teacherID := h.userID(c)
	questionID := c.Param("questionId")
	if questionID == "" {
		utils.JSONBadRequest(c, "questionId required")
		return
	}
	err := h.quizService.DeleteQuestion(ctx, teacherID, questionID)
	if err != nil {
		switch err {
		case services.ErrQuestionNotFound:
			utils.JSONNotFound(c, "question not found")
			return
		case services.ErrForbidden:
			utils.JSONForbidden(c, "forbidden")
			return
		default:
			utils.JSONInternal(c, "failed to delete question")
			return
		}
	}
	utils.JSONSuccess(c, gin.H{"message": "deleted"})
}

// ListQuizzes handles GET /api/v1/teacher/quizzes.
func (h *TeacherHandler) ListQuizzes(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	teacherID := h.userID(c)
	list, err := h.quizService.ListQuizzesForTeacher(ctx, teacherID)
	if err != nil {
		utils.JSONInternal(c, "failed to list quizzes")
		return
	}
	utils.JSONSuccess(c, list)
}

// GetQuizSubmissions handles GET /api/v1/teacher/quizzes/:quizId/submissions.
func (h *TeacherHandler) GetQuizSubmissions(c *gin.Context) {
	ctx, cancel := h.withTimeout(c)
	defer cancel()
	teacherID := h.userID(c)
	quizID := c.Param("quizId")
	if quizID == "" {
		utils.JSONBadRequest(c, "quizId required")
		return
	}
	list, err := h.quizService.GetQuizSubmissions(ctx, teacherID, quizID)
	if err != nil {
		switch err {
		case services.ErrQuizNotFound:
			utils.JSONNotFound(c, "quiz not found")
			return
		case services.ErrForbidden:
			utils.JSONForbidden(c, "forbidden")
			return
		default:
			utils.JSONInternal(c, "failed to get submissions")
			return
		}
	}
	utils.JSONSuccess(c, list)
}
