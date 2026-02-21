package models

// CreateQuizRequest is the body for POST /api/v1/teacher/quizzes.
type CreateQuizRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// AddQuestionRequest is the body for POST /api/v1/teacher/quizzes/:quizId/questions.
type AddQuestionRequest struct {
	QuestionText  string `json:"question_text" binding:"required"`
	OptionA       string `json:"option_a" binding:"required"`
	OptionB       string `json:"option_b" binding:"required"`
	OptionC       string `json:"option_c" binding:"required"`
	OptionD       string `json:"option_d" binding:"required"`
	CorrectAnswer string `json:"correct_answer" binding:"required"`
}

// UpdateQuestionRequest is the body for PUT /api/v1/teacher/questions/:questionId.
type UpdateQuestionRequest struct {
	QuestionText  string `json:"question_text" binding:"required"`
	OptionA       string `json:"option_a" binding:"required"`
	OptionB       string `json:"option_b" binding:"required"`
	OptionC       string `json:"option_c" binding:"required"`
	OptionD       string `json:"option_d" binding:"required"`
	CorrectAnswer string `json:"correct_answer" binding:"required"`
}

// SubmitQuizRequest is the body for POST /api/v1/quizzes/:quizId/submit.
// Keys are question IDs, values are selected option (e.g. "A", "B").
type SubmitQuizRequest struct {
	Answers map[string]string `json:"answers" binding:"required"`
}
