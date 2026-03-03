package handler

import (
	"errors"
	"net/http"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func successResponse(c *gin.Context, status int, data any) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
	})
}

func messageResponse(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
	})
}

func errorResponse(c *gin.Context, status int, err string) {
	c.JSON(status, Response{
		Success: false,
		Error:   err,
	})
}

func AbortWithError(c *gin.Context, status int, err string) {
	c.AbortWithStatusJSON(status, Response{
		Success: false,
		Error:   err,
	})
}

func domainErrorResponce(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		errorResponse(c, http.StatusNotFound, "user not found")
	case errors.Is(err, domain.ErrUserAlreadyExsists):
		errorResponse(c, http.StatusConflict, "email already registered")
	case errors.Is(err, domain.ErrInvalidCredentials):
		errorResponse(c, http.StatusUnauthorized, "invalid email or password")
	case errors.Is(err, domain.ErrInvalidToken):
		errorResponse(c, http.StatusUnauthorized, "invalid token")
	case errors.Is(err, domain.ErrExpiredToken):
		errorResponse(c, http.StatusUnauthorized, "token has expired")
	case errors.Is(err, domain.ErrNotVerified):
		errorResponse(c, http.StatusForbidden, "email not verified")
	case errors.Is(err, domain.ErrAlreadyVerified):
		errorResponse(c, http.StatusConflict, "email already verified")
	case errors.Is(err, domain.ErrForbidden):
		errorResponse(c, http.StatusForbidden, "forbidden")
	default:
		errorResponse(c, http.StatusInternalServerError, "internal server error")
	}
}
