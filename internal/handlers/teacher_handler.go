package handlers

import (
	"net/http"

	"quiz-backend/internal/models"
	"quiz-backend/internal/store"

	"github.com/gin-gonic/gin"
)

// TeacherHandler handles teacher-facing quiz and question CRUD endpoints.
type TeacherHandler struct {
	store *store.MemoryStore
}

// NewTeacherHandler creates a teacher handler with the given in-memory store.
func NewTeacherHandler(store *store.MemoryStore) *TeacherHandler {
	return &TeacherHandler{store: store}
}

// CreateQuiz handles POST /teacher/create_quiz
func (h *TeacherHandler) CreateQuiz(c *gin.Context) {
	var req models.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz := models.Quiz{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Questions:   []models.Question{},
	}
	if !h.store.CreateQuiz(quiz) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz id already exists"})
		return
	}
	c.JSON(http.StatusCreated, quiz)
}

// AddQuestion handles POST /teacher/add_question/:quiz_id
func (h *TeacherHandler) AddQuestion(c *gin.Context) {
	quizID := c.Param("quiz_id")
	if quizID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "quiz_id is required"})
		return
	}

	var req models.AddQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	q := models.Question{
		ID:            req.ID,
		QuestionText:  req.QuestionText,
		OptionA:       req.OptionA,
		OptionB:       req.OptionB,
		OptionC:       req.OptionC,
		OptionD:       req.OptionD,
		CorrectAnswer: req.CorrectAnswer,
	}
	if !h.store.AddQuestion(quizID, q) {
		c.JSON(http.StatusNotFound, gin.H{"error": "quiz not found"})
		return
	}
	c.JSON(http.StatusCreated, q)
}

// UpdateQuestion handles PUT /teacher/update_question/:question_id
func (h *TeacherHandler) UpdateQuestion(c *gin.Context) {
	questionID := c.Param("question_id")
	if questionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question_id is required"})
		return
	}

	var req models.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := models.Question{
		ID:            questionID,
		QuestionText:  req.QuestionText,
		OptionA:       req.OptionA,
		OptionB:       req.OptionB,
		OptionC:       req.OptionC,
		OptionD:       req.OptionD,
		CorrectAnswer: req.CorrectAnswer,
	}
	if !h.store.UpdateQuestion(questionID, update) {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}
	c.JSON(http.StatusOK, update)
}

// DeleteQuestion handles DELETE /teacher/delete_question/:question_id
func (h *TeacherHandler) DeleteQuestion(c *gin.Context) {
	questionID := c.Param("question_id")
	if questionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question_id is required"})
		return
	}

	if !h.store.DeleteQuestion(questionID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "question not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "question deleted"})
}

// ViewQuiz handles GET /teacher/view_quiz/:quiz_id (full quiz including correct answers).
func (h *TeacherHandler) ViewQuiz(c *gin.Context) {
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
	c.JSON(http.StatusOK, quiz)
}

// AllQuizzes handles GET /teacher/all_quizzes
func (h *TeacherHandler) AllQuizzes(c *gin.Context) {
	quizzes := h.store.GetAllQuizzes()
	c.JSON(http.StatusOK, quizzes)
}
