package http

import (
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/erp-cosmetics/user-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	userHandler       *handler.UserHandler
	departmentHandler *handler.DepartmentHandler
	healthHandler     *handler.HealthHandler
	logger            *zap.Logger
}

func NewRouter(
	userHandler *handler.UserHandler,
	departmentHandler *handler.DepartmentHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		userHandler:       userHandler,
		departmentHandler: departmentHandler,
		healthHandler:     healthHandler,
		logger:            logger,
	}
}

func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Global middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(r.logger))
	router.Use(middleware.Recovery(r.logger))
	router.Use(middleware.CORS("*"))

	// Health endpoints
	router.GET("/health", r.healthHandler.Health)
	router.GET("/ready", r.healthHandler.Ready)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Users
		users := v1.Group("/users")
		{
			users.GET("/me", r.userHandler.GetMe)
			users.GET("", r.userHandler.ListUsers)
			users.POST("", r.userHandler.CreateUser)
			users.GET("/:id", r.userHandler.GetUser)
			// users.PUT("/:id", r.userHandler.UpdateUser)
			// users.DELETE("/:id", r.userHandler.DeleteUser)
		}

		// Departments
		departments := v1.Group("/departments")
		{
			departments.GET("", r.departmentHandler.GetDepartmentTree)
			departments.POST("", r.departmentHandler.CreateDepartment)
			// departments.GET("/:id", r.departmentHandler.GetDepartment)
			// departments.PUT("/:id", r.departmentHandler.UpdateDepartment)
			// departments.DELETE("/:id", r.departmentHandler.DeleteDepartment)
		}
	}

	return router
}
