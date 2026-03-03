package handler

import (
	"net/http"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/alikurb12/auth_service_jwt_golang/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc usecase.AuthUseCase
}

func NewAuthHandler(uc usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input domain.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.uc.Register(c.Request.Context(), input)
	if err != nil {
		domainErrorResponse(c, err)
		return
	}

	successResponse(c, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input domain.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.uc.Login(c.Request.Context(), input)
	if err != nil {
		domainErrorResponse(c, err)
		return
	}

	successResponse(c, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input domain.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.uc.Refresh(c.Request.Context(), input.RefreshToken)
	if err != nil {
		domainErrorResponse(c, err)
		return
	}

	successResponse(c, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var input domain.RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.uc.Logout(c.Request.Context(), input.RefreshToken); err != nil {
		domainErrorResponse(c, err)
		return
	}

	messageResponse(c, http.StatusOK, "logged out successfully")
}

func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.uc.LogoutAll(c.Request.Context(), userID); err != nil {
		domainErrorResponse(c, err)
		return
	}

	messageResponse(c, http.StatusOK, "logged out from all devices")
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		errorResponse(c, http.StatusBadRequest, "token is required")
		return
	}

	if err := h.uc.VerifyEmail(c.Request.Context(), token); err != nil {
		domainErrorResponse(c, err)
		return
	}

	messageResponse(c, http.StatusOK, "email verified successfully")
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.uc.ResendVerifiction(c.Request.Context(), input.Email); err != nil {
		domainErrorResponse(c, err)
		return
	}

	messageResponse(c, http.StatusOK, "if this email is registered, a verification link has been sent")
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, ok := GetUserID(c)
	if !ok {
		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.uc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		domainErrorResponse(c, err)
		return
	}

	successResponse(c, http.StatusOK, user)
}