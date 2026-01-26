package http

import (
	"github.com/erp-cosmetics/file-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router handles HTTP routing
type Router struct {
	fileHandler   *handler.FileHandler
	healthHandler *handler.HealthHandler
	logger        *zap.Logger
}

// NewRouter creates new router
func NewRouter(
	fileHandler *handler.FileHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		fileHandler:   fileHandler,
		healthHandler: healthHandler,
		logger:        logger,
	}
}

// Setup sets up the HTTP router
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS("*"))
	router.Use(middleware.Logger(r.logger))

	// Increase max multipart memory
	router.MaxMultipartMemory = 32 << 20 // 32 MB

	// Health check
	router.GET("/health", r.healthHandler.HealthCheck)
	router.GET("/ready", r.healthHandler.ReadyCheck)

	// API v1
	v1 := router.Group("/api/v1")
	{
		files := v1.Group("/files")
		{
			files.POST("/upload", r.fileHandler.Upload)
			files.POST("/upload/multiple", r.fileHandler.UploadMultiple)
			files.GET("/categories", r.fileHandler.ListCategories)
			files.GET("/entity/:type/:id", r.fileHandler.GetByEntity)
			files.GET("/:id", r.fileHandler.Get)
			files.GET("/:id/download", r.fileHandler.Download)
			files.GET("/:id/url", r.fileHandler.GetDownloadURL)
			files.DELETE("/:id", r.fileHandler.Delete)
		}
	}

	return router
}
