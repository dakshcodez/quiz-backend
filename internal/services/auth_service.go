package services

import (
	"context"
	"errors"
	"time"

	"quiz-backend/internal/models"
	"quiz-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists   = errors.New("email already registered")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrInvalidRole  = errors.New("role must be teacher or student")
)

// AuthService handles registration and login; no JWT logic in handlers.
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

// NewAuthService returns a new AuthService with dependencies injected.
func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: jwtExpiry,
	}
}

// Register creates a user. Returns ErrEmailExists if email is taken, ErrInvalidRole if role is invalid.
func (s *AuthService) Register(ctx context.Context, email, password, role string) (*models.User, error) {
	if !models.ValidRole(role) {
		return nil, ErrInvalidRole
	}
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u, err := s.userRepo.Create(ctx, email, string(hash), role)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Login validates credentials and returns user + JWT token, or ErrInvalidCreds.
func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCreds
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCreds
	}
	token, err := s.generateToken(u.ID, u.Role)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

// generateToken creates a JWT with user_id, role, and exp.
func (s *AuthService) generateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(s.jwtExpiry).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}
