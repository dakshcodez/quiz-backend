package models

import "time"

// User represents a row in the users table.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never expose in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

// ValidRole returns true if role is teacher or student.
func ValidRole(role string) bool {
	return role == "teacher" || role == "student"
}
