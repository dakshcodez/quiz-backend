package models

import "time"

// Submission represents a row in the submissions table.
type Submission struct {
	ID         string    `json:"id"`
	QuizID     string    `json:"quiz_id"`
	StudentID  string    `json:"student_id"`
	Score      int       `json:"score"`
	Total      int       `json:"total"`
	CreatedAt  time.Time `json:"created_at"`
}

// SubmissionAnswer represents a row in the submission_answers table.
type SubmissionAnswer struct {
	ID             string    `json:"id"`
	SubmissionID   string    `json:"submission_id"`
	QuestionID     string    `json:"question_id"`
	SelectedAnswer string    `json:"selected_answer"`
	IsCorrect      bool      `json:"is_correct"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

// SubmissionDetailItem is one entry in the submit response "details" array.
// Does not expose the correct answer, only whether the choice was correct.
type SubmissionDetailItem struct {
	QuestionID string `json:"question_id"`
	Correct    bool   `json:"correct"`
}
