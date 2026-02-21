package store

import (
	"sync"

	"quiz-backend/internal/models"
)

// MemoryStore holds all quizzes in memory. Safe for concurrent use.
// Data is lost when the server restarts.
//
// We use sync.RWMutex so that many goroutines can read (ViewQuiz) at the same time,
// while any write would lock exclusively. For Phase 1 we only read; the mutex
// keeps the map safe if you add writes later (e.g. create quiz).
type MemoryStore struct {
	mu      sync.RWMutex
	quizzes map[string]models.Quiz
}

// NewMemoryStore creates a new in-memory store and seeds it with sample data.
func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		quizzes: make(map[string]models.Quiz),
	}
	s.seed()
	return s
}

// seed adds pre-defined quiz data for teaching/demo. Called once at startup.
func (s *MemoryStore) seed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.quizzes["quiz1"] = models.Quiz{
		ID:          "quiz1",
		Title:       "Sample Quiz",
		Description: "A short quiz for the workshop",
		Questions: []models.Question{
			{
				ID:            "question1",
				QuestionText:  "What is 2 + 2?",
				OptionA:       "3",
				OptionB:       "4",
				OptionC:       "5",
				OptionD:       "6",
				CorrectAnswer: "B",
			},
			{
				ID:            "question2",
				QuestionText:  "What is the capital of France?",
				OptionA:       "London",
				OptionB:       "Berlin",
				OptionC:       "Paris",
				OptionD:       "Madrid",
				CorrectAnswer: "C",
			},
		},
	}
}

// GetQuiz returns the quiz by ID. Second return is false if not found.
// We hold RLock for the whole function so the map isn't modified while we read.
func (s *MemoryStore) GetQuiz(id string) (models.Quiz, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	quiz, ok := s.quizzes[id]
	if !ok {
		return models.Quiz{}, false
	}
	return quiz, true
}

// CreateQuiz adds a new quiz. Returns false if a quiz with the same ID already exists.
func (s *MemoryStore) CreateQuiz(quiz models.Quiz) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.quizzes[quiz.ID]; exists {
		return false
	}
	if quiz.Questions == nil {
		quiz.Questions = []models.Question{}
	}
	s.quizzes[quiz.ID] = quiz
	return true
}

// AddQuestion appends a question to a quiz. Returns false if the quiz is not found.
func (s *MemoryStore) AddQuestion(quizID string, q models.Question) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	quiz, ok := s.quizzes[quizID]
	if !ok {
		return false
	}
	newQuestions := make([]models.Question, len(quiz.Questions), len(quiz.Questions)+1)
	copy(newQuestions, quiz.Questions)
	quiz.Questions = append(newQuestions, q)
	s.quizzes[quizID] = quiz
	return true
}

// UpdateQuestion finds a question by ID across all quizzes and updates it. Returns false if not found.
func (s *MemoryStore) UpdateQuestion(questionID string, update models.Question) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for quizID, quiz := range s.quizzes {
		for i := range quiz.Questions {
			if quiz.Questions[i].ID == questionID {
				newQuestions := make([]models.Question, len(quiz.Questions))
				copy(newQuestions, quiz.Questions)
				update.ID = questionID
				newQuestions[i] = update
				quiz.Questions = newQuestions
				s.quizzes[quizID] = quiz
				return true
			}
		}
	}
	return false
}

// DeleteQuestion removes a question by ID from its quiz. Returns false if not found.
func (s *MemoryStore) DeleteQuestion(questionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for quizID, quiz := range s.quizzes {
		for i, q := range quiz.Questions {
			if q.ID == questionID {
				newQuestions := make([]models.Question, 0, len(quiz.Questions)-1)
				newQuestions = append(newQuestions, quiz.Questions[:i]...)
				newQuestions = append(newQuestions, quiz.Questions[i+1:]...)
				quiz.Questions = newQuestions
				s.quizzes[quizID] = quiz
				return true
			}
		}
	}
	return false
}

// GetAllQuizzes returns a copy of all quizzes (for teacher list view).
func (s *MemoryStore) GetAllQuizzes() []models.Quiz {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.Quiz, 0, len(s.quizzes))
	for _, quiz := range s.quizzes {
		out = append(out, quiz)
	}
	return out
}
