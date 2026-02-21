package repository

import (
	"context"

	"quiz-backend/internal/database"
	"quiz-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// QuizRepository handles quizzes and questions table access.
type QuizRepository struct {
	db *database.DB
}

// NewQuizRepository returns a new QuizRepository.
func NewQuizRepository(db *database.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

// CreateQuiz inserts a quiz. teacherID is the authenticated teacher's user ID.
func (r *QuizRepository) CreateQuiz(ctx context.Context, title, description, teacherID string) (*models.Quiz, error) {
	id := uuid.New().String()
	query := `INSERT INTO quizzes (id, title, description, teacher_id) VALUES ($1, $2, $3, $4)
	          RETURNING id, title, description, teacher_id, created_at`
	var q models.Quiz
	err := r.db.Pool.QueryRow(ctx, query, id, title, description, teacherID).Scan(
		&q.ID, &q.Title, &q.Description, &q.TeacherID, &q.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// GetQuizByID returns a quiz by ID. Returns pgx.ErrNoRows if not found.
func (r *QuizRepository) GetQuizByID(ctx context.Context, id string) (*models.Quiz, error) {
	query := `SELECT id, title, description, teacher_id, created_at FROM quizzes WHERE id = $1`
	var q models.Quiz
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&q.ID, &q.Title, &q.Description, &q.TeacherID, &q.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// GetQuizWithQuestions returns a quiz and all its questions (including correct_answer).
// For teacher use.
func (r *QuizRepository) GetQuizWithQuestions(ctx context.Context, quizID string) (*models.Quiz, error) {
	q, err := r.GetQuizByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	questions, err := r.GetQuestionsByQuizID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	q.Questions = questions
	return q, nil
}

// GetQuizWithQuestionsForStudent returns quiz with questions but WITHOUT correct_answer.
func (r *QuizRepository) GetQuizWithQuestionsForStudent(ctx context.Context, quizID string) (*models.QuizView, error) {
	q, err := r.GetQuizByID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	questions, err := r.GetQuestionsByQuizID(ctx, quizID)
	if err != nil {
		return nil, err
	}
	view := &models.QuizView{
		ID:          q.ID,
		Title:       q.Title,
		Description: q.Description,
		Questions:    make([]models.QuestionView, 0, len(questions)),
	}
	for _, qu := range questions {
		view.Questions = append(view.Questions, models.QuestionView{
			ID:           qu.ID,
			QuestionText: qu.QuestionText,
			OptionA:      qu.OptionA,
			OptionB:      qu.OptionB,
			OptionC:      qu.OptionC,
			OptionD:      qu.OptionD,
		})
	}
	return view, nil
}

// GetQuestionsByQuizID returns all questions for a quiz (including correct_answer).
func (r *QuizRepository) GetQuestionsByQuizID(ctx context.Context, quizID string) ([]models.Question, error) {
	query := `SELECT id, quiz_id, question_text, option_a, option_b, option_c, option_d, correct_answer, created_at
	          FROM questions WHERE quiz_id = $1 ORDER BY created_at`
	rows, err := r.db.Pool.Query(ctx, query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.ID, &q.QuizID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.CorrectAnswer, &q.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}

// ListQuizzes returns all quizzes (id, title, description, teacher_id, created_at). No questions.
func (r *QuizRepository) ListQuizzes(ctx context.Context) ([]models.Quiz, error) {
	query := `SELECT id, title, description, teacher_id, created_at FROM quizzes ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Quiz
	for rows.Next() {
		var q models.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.TeacherID, &q.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}

// ListQuizzesByTeacherID returns quizzes created by the given teacher.
func (r *QuizRepository) ListQuizzesByTeacherID(ctx context.Context, teacherID string) ([]models.Quiz, error) {
	query := `SELECT id, title, description, teacher_id, created_at FROM quizzes WHERE teacher_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Pool.Query(ctx, query, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Quiz
	for rows.Next() {
		var q models.Quiz
		if err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.TeacherID, &q.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, q)
	}
	return list, rows.Err()
}

// AddQuestion inserts a question into a quiz.
func (r *QuizRepository) AddQuestion(ctx context.Context, quizID, questionText, optionA, optionB, optionC, optionD, correctAnswer string) (*models.Question, error) {
	id := uuid.New().String()
	query := `INSERT INTO questions (id, quiz_id, question_text, option_a, option_b, option_c, option_d, correct_answer)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id, quiz_id, question_text, option_a, option_b, option_c, option_d, correct_answer, created_at`
	var q models.Question
	err := r.db.Pool.QueryRow(ctx, query, id, quizID, questionText, optionA, optionB, optionC, optionD, correctAnswer).Scan(
		&q.ID, &q.QuizID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.CorrectAnswer, &q.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// UpdateQuestion updates a question by ID.
func (r *QuizRepository) UpdateQuestion(ctx context.Context, questionID, questionText, optionA, optionB, optionC, optionD, correctAnswer string) error {
	query := `UPDATE questions SET question_text = $2, option_a = $3, option_b = $4, option_c = $5, option_d = $6, correct_answer = $7 WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, query, questionID, questionText, optionA, optionB, optionC, optionD, correctAnswer)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// DeleteQuestion deletes a question by ID.
func (r *QuizRepository) DeleteQuestion(ctx context.Context, questionID string) error {
	query := `DELETE FROM questions WHERE id = $1`
	result, err := r.db.Pool.Exec(ctx, query, questionID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// GetQuestionByID returns a single question. For ownership checks.
func (r *QuizRepository) GetQuestionByID(ctx context.Context, questionID string) (*models.Question, error) {
	query := `SELECT id, quiz_id, question_text, option_a, option_b, option_c, option_d, correct_answer, created_at FROM questions WHERE id = $1`
	var q models.Question
	err := r.db.Pool.QueryRow(ctx, query, questionID).Scan(
		&q.ID, &q.QuizID, &q.QuestionText, &q.OptionA, &q.OptionB, &q.OptionC, &q.OptionD, &q.CorrectAnswer, &q.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
