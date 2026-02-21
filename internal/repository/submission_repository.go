package repository

import (
	"context"

	"quiz-backend/internal/database"
	"quiz-backend/internal/models"

	"github.com/google/uuid"
)

// SubmissionRepository handles submissions and submission_answers tables.
type SubmissionRepository struct {
	db *database.DB
}

// NewSubmissionRepository returns a new SubmissionRepository.
func NewSubmissionRepository(db *database.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

// CreateSubmission inserts a submission and returns it.
func (r *SubmissionRepository) CreateSubmission(ctx context.Context, quizID, studentID string, score, total int) (*models.Submission, error) {
	id := uuid.New().String()
	query := `INSERT INTO submissions (id, quiz_id, student_id, score, total) VALUES ($1, $2, $3, $4, $5)
	          RETURNING id, quiz_id, student_id, score, total, created_at`
	var s models.Submission
	err := r.db.Pool.QueryRow(ctx, query, id, quizID, studentID, score, total).Scan(
		&s.ID, &s.QuizID, &s.StudentID, &s.Score, &s.Total, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// CreateSubmissionAnswers inserts multiple submission_answers in one go.
func (r *SubmissionRepository) CreateSubmissionAnswers(ctx context.Context, submissionID string, answers []models.SubmissionAnswer) error {
	for _, a := range answers {
		id := uuid.New().String()
		query := `INSERT INTO submission_answers (id, submission_id, question_id, selected_answer, is_correct) VALUES ($1, $2, $3, $4, $5)`
		_, err := r.db.Pool.Exec(ctx, query, id, submissionID, a.QuestionID, a.SelectedAnswer, a.IsCorrect)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSubmissionsByQuizID returns all submissions for a quiz (for teacher).
func (r *SubmissionRepository) GetSubmissionsByQuizID(ctx context.Context, quizID string) ([]models.Submission, error) {
	query := `SELECT id, quiz_id, student_id, score, total, created_at FROM submissions WHERE quiz_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(&s.ID, &s.QuizID, &s.StudentID, &s.Score, &s.Total, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// GetSubmissionsByStudentID returns all submissions by a student.
func (r *SubmissionRepository) GetSubmissionsByStudentID(ctx context.Context, studentID string) ([]models.Submission, error) {
	query := `SELECT id, quiz_id, student_id, score, total, created_at FROM submissions WHERE student_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Submission
	for rows.Next() {
		var s models.Submission
		if err := rows.Scan(&s.ID, &s.QuizID, &s.StudentID, &s.Score, &s.Total, &s.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

// GetSubmissionByID returns a submission by ID. Returns pgx.ErrNoRows if not found.
func (r *SubmissionRepository) GetSubmissionByID(ctx context.Context, id string) (*models.Submission, error) {
	query := `SELECT id, quiz_id, student_id, score, total, created_at FROM submissions WHERE id = $1`
	var s models.Submission
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.QuizID, &s.StudentID, &s.Score, &s.Total, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
