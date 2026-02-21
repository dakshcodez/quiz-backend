package models

// Quiz is the in-memory representation of a quiz (used by the store).
// Questions include CorrectAnswer; do not expose it in API responses.
type Quiz struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

// Question holds full question data including the correct answer (server-side only).
type Question struct {
	ID            string `json:"id"`
	QuestionText  string `json:"question_text"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectAnswer string `json:"correct_answer"` // Used only in store; never returned to client.
}

// ViewQuestion is used in GET /student/view_quiz. Does not expose correct_answer.
type ViewQuestion struct {
	ID           string `json:"id"`
	QuestionText string `json:"question_text"`
	OptionA      string `json:"option_a"`
	OptionB      string `json:"option_b"`
	OptionC      string `json:"option_c"`
	OptionD      string `json:"option_d"`
}

// ViewQuizResponse is the JSON response for GET /student/view_quiz/:quiz_id.
type ViewQuizResponse struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Questions   []ViewQuestion `json:"questions"`
}

// GiveQuizRequest is the request body for POST /student/give_quiz.
type GiveQuizRequest struct {
	QuizID  string            `json:"quiz_id" binding:"required"`
	Answers map[string]string `json:"answers" binding:"required"`
}
