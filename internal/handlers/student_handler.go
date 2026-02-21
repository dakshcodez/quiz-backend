package handlers

import (
	"net/http"

	"quiz-backend/internal/models"
	"quiz-backend/internal/store"

	"github.com/gin-gonic/gin"
)

// StudentHandler handles student-facing quiz endpoints.
type StudentHandler struct {
	store *store.MemoryStore
}

// NewStudentHandler creates a student handler with the given in-memory store.
func NewStudentHandler(store *store.MemoryStore) *StudentHandler {
	return &StudentHandler{store: store}
}

// ViewQuiz handles GET /student/view_quiz/:quiz_id
// Returns the quiz with questions, but never exposes correct_answer.
func (h *StudentHandler) ViewQuiz(c *gin.Context) {
	quizID := c.Param("quiz_id")
	if quizID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	quiz, ok := h.store.GetQuiz(quizID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}

	// Build response without exposing correct_answer (teaching: never send answers to client).
	questions := make([]models.ViewQuestion, 0, len(quiz.Questions))
	for _, q := range quiz.Questions {
		questions = append(questions, models.ViewQuestion{
			ID:           q.ID,
			QuestionText: q.QuestionText,
			OptionA:      q.OptionA,
			OptionB:      q.OptionB,
			OptionC:      q.OptionC,
			OptionD:      q.OptionD,
		})
	}

	resp := models.ViewQuizResponse{
		ID:          quiz.ID,
		Title:       quiz.Title,
		Description: quiz.Description,
		Questions:   questions,
	}
	c.JSON(http.StatusOK, resp)
}

// GiveQuiz handles POST /student/give_quiz
// Phase 1: only validates JSON and returns a mock score. Real scoring in Phase 3.
func (h *StudentHandler) GiveQuiz(c *gin.Context) {
	var req models.GiveQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.QuizID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	// Phase 1: mock response. No real scoring yet.
	c.JSON(http.StatusOK, gin.H{
		"quiz_id": req.QuizID,
		"score":   5,
		"message": "Phase 1 mock scoring",
	})
}
