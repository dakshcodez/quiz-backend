package repository

import (
	"context"

	"quiz-backend/internal/database"
	"quiz-backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserRepository handles user table access.
type UserRepository struct {
	db *database.DB
}

// NewUserRepository returns a new UserRepository.
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a user. Returns the created user or error.
func (r *UserRepository) Create(ctx context.Context, email, passwordHash, role string) (*models.User, error) {
	id := uuid.New().String()
	query := `INSERT INTO users (id, email, password_hash, role) VALUES ($1, $2, $3, $4)
	          RETURNING id, email, password_hash, role, created_at`
	var u models.User
	err := r.db.Pool.QueryRow(ctx, query, id, email, passwordHash, role).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail returns a user by email. Returns pgx.ErrNoRows if not found.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1`
	var u models.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByID returns a user by ID. Returns pgx.ErrNoRows if not found.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1`
	var u models.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ExistsByEmail returns true if a user with the email exists.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1`
	var exists int
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&exists)
	if err == pgx.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
