package handler

import (
	"net/http"
	"strings"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/alikurb12/auth_service_jwt_golang/internal/infrastructure/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextUserID = "userID"
	ContextRole   = "userRole"
)

func AuthMiddleware(jwtService *jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			AbortWithError(c, http.StatusUnauthorized, "authorization header required")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			AbortWithError(c, http.StatusUnauthorized, "invalid authorization header format")
			return
		}

		claims, err := jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			if err == domain.ErrExpiredToken {
				AbortWithError(c, http.StatusUnauthorized, "token has expired")
			} else {
				AbortWithError(c, http.StatusUnauthorized, "invalid token")
			}
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			AbortWithError(c, http.StatusUnauthorized, "invalid token claims")
			return
		}

		c.Set(ContextUserID, userID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...domain.Role) gin.HandlerFunc {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get(ContextRole)
		if !exists {
			AbortWithError(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		userRole, ok := role.(domain.Role)
		if !ok {
			AbortWithError(c, http.StatusInternalServerError, "internal server error")
			return
		}
		if _, ok := allowed[userRole]; !ok {
			AbortWithError(c, http.StatusForbidden, "forbidden: insufficient permissions")
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextUserID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
