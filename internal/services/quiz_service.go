package services

import (
	"context"
	"errors"
	"strings"

	"quiz-backend/internal/models"
	"quiz-backend/internal/repository"

	"github.com/jackc/pgx/v5"
)

var (
	ErrQuizNotFound      = errors.New("quiz not found")
	ErrQuestionNotFound  = errors.New("question not found")
	ErrForbidden        = errors.New("forbidden")
)

// QuizService handles quiz and submission business logic; no DB in handlers.
type QuizService struct {
	quizRepo       *repository.QuizRepository
	submissionRepo *repository.SubmissionRepository
}

// NewQuizService returns a new QuizService.
func NewQuizService(quizRepo *repository.QuizRepository, submissionRepo *repository.SubmissionRepository) *QuizService {
	return &QuizService{
		quizRepo:       quizRepo,
		submissionRepo: submissionRepo,
	}
}

// --- Teacher ---

// CreateQuiz creates a quiz for the given teacher.
func (s *QuizService) CreateQuiz(ctx context.Context, teacherID, title, description string) (*models.Quiz, error) {
	return s.quizRepo.CreateQuiz(ctx, title, description, teacherID)
}

// AddQuestion adds a question to a quiz. Returns ErrQuizNotFound or ErrForbidden if quiz not owned by teacher.
func (s *QuizService) AddQuestion(ctx context.Context, teacherID, quizID string, req *models.AddQuestionRequest) (*models.Question, error) {
	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err == pgx.ErrNoRows {
		return nil, ErrQuizNotFound
	}
	if err != nil {
		return nil, err
	}
	if quiz.TeacherID != teacherID {
		return nil, ErrForbidden
	}
	return s.quizRepo.AddQuestion(ctx, quizID, req.QuestionText, req.OptionA, req.OptionB, req.OptionC, req.OptionD, req.CorrectAnswer)
}

// UpdateQuestion updates a question. Returns ErrQuestionNotFound or ErrForbidden if not owned.
func (s *QuizService) UpdateQuestion(ctx context.Context, teacherID, questionID string, req *models.UpdateQuestionRequest) error {
	q, err := s.quizRepo.GetQuestionByID(ctx, questionID)
	if err == pgx.ErrNoRows {
		return ErrQuestionNotFound
	}
	if err != nil {
		return err
	}
	quiz, err := s.quizRepo.GetQuizByID(ctx, q.QuizID)
	if err != nil {
		return err
	}
	if quiz.TeacherID != teacherID {
		return ErrForbidden
	}
	return s.quizRepo.UpdateQuestion(ctx, questionID, req.QuestionText, req.OptionA, req.OptionB, req.OptionC, req.OptionD, req.CorrectAnswer)
}

// DeleteQuestion deletes a question. Returns ErrQuestionNotFound or ErrForbidden if not owned.
func (s *QuizService) DeleteQuestion(ctx context.Context, teacherID, questionID string) error {
	q, err := s.quizRepo.GetQuestionByID(ctx, questionID)
	if err == pgx.ErrNoRows {
		return ErrQuestionNotFound
	}
	if err != nil {
		return err
	}
	quiz, err := s.quizRepo.GetQuizByID(ctx, q.QuizID)
	if err != nil {
		return err
	}
	if quiz.TeacherID != teacherID {
		return ErrForbidden
	}
	return s.quizRepo.DeleteQuestion(ctx, questionID)
}

// ListQuizzesForTeacher returns all quizzes for the given teacher.
func (s *QuizService) ListQuizzesForTeacher(ctx context.Context, teacherID string) ([]models.Quiz, error) {
	return s.quizRepo.ListQuizzesByTeacherID(ctx, teacherID)
}

// GetQuizSubmissions returns submissions for a quiz. Teacher must own the quiz.
func (s *QuizService) GetQuizSubmissions(ctx context.Context, teacherID, quizID string) ([]models.Submission, error) {
	quiz, err := s.quizRepo.GetQuizByID(ctx, quizID)
	if err == pgx.ErrNoRows {
		return nil, ErrQuizNotFound
	}
	if err != nil {
		return nil, err
	}
	if quiz.TeacherID != teacherID {
		return nil, ErrForbidden
	}
	return s.submissionRepo.GetSubmissionsByQuizID(ctx, quizID)
}

// --- Student ---

// ListQuizzes returns all quizzes (for student browse).
func (s *QuizService) ListQuizzes(ctx context.Context) ([]models.Quiz, error) {
	return s.quizRepo.ListQuizzes(ctx)
}

// GetQuizForStudent returns a quiz with questions but without correct answers.
func (s *QuizService) GetQuizForStudent(ctx context.Context, quizID string) (*models.QuizView, error) {
	view, err := s.quizRepo.GetQuizWithQuestionsForStudent(ctx, quizID)
	if err == pgx.ErrNoRows {
		return nil, ErrQuizNotFound
	}
	if err != nil {
		return nil, err
	}
	return view, nil
}

// SubmitQuizResult holds the response for a quiz submission.
type SubmitQuizResult struct {
	QuizID   string                      `json:"quiz_id"`
	Score    int                         `json:"score"`
	Total    int                         `json:"total"`
	Details  []models.SubmissionDetailItem `json:"details"`
}

// SubmitQuiz scores the answers, persists submission and submission_answers, returns result. Does not expose correct answers.
func (s *QuizService) SubmitQuiz(ctx context.Context, studentID, quizID string, answers map[string]string) (*SubmitQuizResult, error) {
	quiz, err := s.quizRepo.GetQuizWithQuestions(ctx, quizID)
	if err == pgx.ErrNoRows {
		return nil, ErrQuizNotFound
	}
	if err != nil {
		return nil, err
	}

	// Score: compare each question's correct_answer to submitted answer (normalize case).
	total := len(quiz.Questions)
	score := 0
	details := make([]models.SubmissionDetailItem, 0, total)
	submissionAnswers := make([]models.SubmissionAnswer, 0, total)

	for _, q := range quiz.Questions {
		selected := strings.TrimSpace(strings.ToUpper(answers[q.ID]))
		correct := strings.TrimSpace(strings.ToUpper(q.CorrectAnswer))
		isCorrect := selected == correct
		if isCorrect {
			score++
		}
		details = append(details, models.SubmissionDetailItem{QuestionID: q.ID, Correct: isCorrect})
		submissionAnswers = append(submissionAnswers, models.SubmissionAnswer{
			QuestionID:     q.ID,
			SelectedAnswer: answers[q.ID],
			IsCorrect:      isCorrect,
		})
	}

	sub, err := s.submissionRepo.CreateSubmission(ctx, quizID, studentID, score, total)
	if err != nil {
		return nil, err
	}
	if err := s.submissionRepo.CreateSubmissionAnswers(ctx, sub.ID, submissionAnswers); err != nil {
		return nil, err
	}

	return &SubmitQuizResult{
		QuizID:  quizID,
		Score:   score,
		Total:   total,
		Details: details,
	}, nil
}

// GetStudentSubmissions returns all submissions for the given student.
func (s *QuizService) GetStudentSubmissions(ctx context.Context, studentID string) ([]models.Submission, error) {
	return s.submissionRepo.GetSubmissionsByStudentID(ctx, studentID)
}
