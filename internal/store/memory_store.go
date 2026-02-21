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
