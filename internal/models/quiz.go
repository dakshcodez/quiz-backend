package models

import "time"

// Quiz represents a row in the quizzes table.
type Quiz struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	TeacherID   string     `json:"teacher_id"`
	CreatedAt   time.Time  `json:"created_at"`
	Questions   []Question `json:"questions,omitempty"`
}

// Question represents a row in the questions table.
// CorrectAnswer is server-side only; do not expose to student responses.
type Question struct {
	ID            string    `json:"id"`
	QuizID        string    `json:"quiz_id"`
	QuestionText  string    `json:"question_text"`
	OptionA       string    `json:"option_a"`
	OptionB       string    `json:"option_b"`
	OptionC       string    `json:"option_c"`
	OptionD       string    `json:"option_d"`
	CorrectAnswer string    `json:"correct_answer,omitempty"` // omit in student API
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

// QuestionView is the student-facing question (no correct_answer).
type QuestionView struct {
	ID           string `json:"id"`
	QuestionText string `json:"question_text"`
	OptionA      string `json:"option_a"`
	OptionB      string `json:"option_b"`
	OptionC      string `json:"option_c"`
	OptionD      string `json:"option_d"`
}

// QuizView is the student-facing quiz (questions without correct answers).
type QuizView struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Questions   []QuestionView `json:"questions"`
}
