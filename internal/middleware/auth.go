package middleware

import (
	"strings"

	"quiz-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret is set at startup so middleware can validate tokens without importing config everywhere.
var JWTSecret []byte

// AuthMiddleware validates the Bearer JWT and injects user_id and role into context.
// Responds 401 if missing or invalid token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.JSONUnauthorized(c, "missing authorization header")
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JSONUnauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}
		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			utils.JSONUnauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.JSONUnauthorized(c, "invalid token claims")
			c.Abort()
			return
		}
		userID, _ := claims["user_id"].(string)
		role, _ := claims["role"].(string)
		if userID == "" || role == "" {
			utils.JSONUnauthorized(c, "invalid token claims")
			c.Abort()
			return
		}
		c.Set(string(utils.ContextKeyUserID), userID)
		c.Set(string(utils.ContextKeyRole), role)
		c.Next()
	}
}

// RequireRole returns a handler that ensures the authenticated user has one of the allowed roles.
// Use after AuthMiddleware.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	set := make(map[string]bool)
	for _, r := range allowedRoles {
		set[r] = true
	}
	return func(c *gin.Context) {
		roleVal, exists := c.Get(string(utils.ContextKeyRole))
		if !exists {
			utils.JSONUnauthorized(c, "unauthorized")
			c.Abort()
			return
		}
		role, _ := roleVal.(string)
		if !set[role] {
			utils.JSONForbidden(c, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
