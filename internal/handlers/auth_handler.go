package handlers

import (
	"net/http"

	"quiz-backend/internal/models"
	"quiz-backend/internal/services"
	"quiz-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles registration and login; delegates to AuthService.
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler returns a new AuthHandler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, err.Error())
		return
	}
	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Role)
	if err != nil {
		switch err {
		case services.ErrEmailExists:
			utils.JSONError(c, http.StatusConflict, "email already registered")
			return
		case services.ErrInvalidRole:
			utils.JSONBadRequest(c, "role must be teacher or student")
			return
		default:
			utils.JSONInternal(c, "registration failed")
			return
		}
	}
	utils.JSONCreated(c, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONBadRequest(c, err.Error())
		return
	}
	user, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCreds {
			utils.JSONUnauthorized(c, "invalid email or password")
			return
		}
		utils.JSONInternal(c, "login failed")
		return
	}
	utils.JSONSuccess(c, models.LoginResponse{
		Token: token,
		User: struct {
			ID    string `json:"id"`
			Email string `json:"email"`
			Role  string `json:"role"`
		}{ID: user.ID, Email: user.Email, Role: user.Role},
	})
}
