package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError is a structured error response.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSONSuccess sends a 200 JSON response.
func JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// JSONCreated sends a 201 JSON response.
func JSONCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// JSONError sends an error response with the given status code and message.
func JSONError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIError{Code: statusCode, Message: message})
}

// JSONBadRequest sends 400.
func JSONBadRequest(c *gin.Context, message string) {
	JSONError(c, http.StatusBadRequest, message)
}

// JSONUnauthorized sends 401.
func JSONUnauthorized(c *gin.Context, message string) {
	JSONError(c, http.StatusUnauthorized, message)
}

// JSONForbidden sends 403.
func JSONForbidden(c *gin.Context, message string) {
	JSONError(c, http.StatusForbidden, message)
}

// JSONNotFound sends 404.
func JSONNotFound(c *gin.Context, message string) {
	JSONError(c, http.StatusNotFound, message)
}

// JSONInternal sends 500.
func JSONInternal(c *gin.Context, message string) {
	JSONError(c, http.StatusInternalServerError, message)
}
