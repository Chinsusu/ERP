package http

import (
	"github.com/erp-cosmetics/auth-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	authHandler   *handler.AuthHandler
	healthHandler *handler.HealthHandler
	logger        *zap.Logger
}

func NewRouter(
	authHandler *handler.AuthHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		authHandler:   authHandler,
		healthHandler: healthHandler,
		logger:        logger,
	}
}

func (r *Router) Setup() *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Global middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(r.logger))
	router.Use(middleware.Recovery(r.logger))
	router.Use(middleware.CORS("*"))

	// Health endpoints (no auth required)
	router.GET("/health", r.healthHandler.Health)
	router.GET("/ready", r.healthHandler.Ready)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Auth endpoints (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/logout", r.authHandler.Logout)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			
			// Protected endpoints
			// auth.GET("/me", middleware.AuthMiddleware(jwtManager), r.authHandler.GetMe)
			// auth.GET("/permissions", middleware.AuthMiddleware(jwtManager), r.authHandler.GetPermissions)
		}
	}

	return router
}
