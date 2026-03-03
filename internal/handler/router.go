package handler

import (
	"net/http"

	"github.com/alikurb12/auth_service_jwt_golang/internal/domain"
	"github.com/alikurb12/auth_service_jwt_golang/internal/infrastructure/jwt"
	"github.com/alikurb12/auth_service_jwt_golang/internal/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	authUC usecase.AuthUseCase,
	jwtService *jwt.JWTService,
) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	authHandler := NewAuthHandler(authUC)

	auth := AuthMiddleware(jwtService)
	adminOnly := RequireRole(domain.RoleAdmin)

	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/refresh", authHandler.Refresh)
			authRoutes.GET("/verify-email", authHandler.VerifyEmail)
			authRoutes.POST("/resend-verification", authHandler.ResendVerification)

			private := authRoutes.Group("", auth)
			{
				private.GET("/me", authHandler.GetProfile)
				private.POST("/logout", authHandler.Logout)
				private.POST("/logout-all", authHandler.LogoutAll)
			}
		}

		adminRoutes := v1.Group("/admin", auth, adminOnly)
		{
			adminRoutes.GET("/users", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin area"})
			})
		}
	}

	return router
}
