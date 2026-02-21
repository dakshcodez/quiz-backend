package utils

// Context keys for values set by middleware. Use type to avoid collisions.
type contextKey string

const (
	ContextKeyUserID contextKey = "user_id"
	ContextKeyRole   contextKey = "role"
)
