package http

import (
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router handles HTTP routing
type Router struct {
	categoryHandler *handler.CategoryHandler
	unitHandler     *handler.UnitHandler
	materialHandler *handler.MaterialHandler
	productHandler  *handler.ProductHandler
	healthHandler   *handler.HealthHandler
	logger          *zap.Logger
}

// NewRouter creates a new router
func NewRouter(
	categoryHandler *handler.CategoryHandler,
	unitHandler *handler.UnitHandler,
	materialHandler *handler.MaterialHandler,
	productHandler *handler.ProductHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		categoryHandler: categoryHandler,
		unitHandler:     unitHandler,
		materialHandler: materialHandler,
		productHandler:  productHandler,
		healthHandler:   healthHandler,
		logger:          logger,
	}
}

// Setup configures the Gin router
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	engine := gin.New()
	
	// Global middlewares
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger(r.logger))
	engine.Use(middleware.Recovery(r.logger))
	engine.Use(middleware.CORS("*"))

	// Health check endpoints
	engine.GET("/health", r.healthHandler.Health)
	engine.GET("/ready", r.healthHandler.Ready)

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Categories
		categories := v1.Group("/categories")
		{
			categories.GET("", r.categoryHandler.List)
			categories.GET("/tree", r.categoryHandler.GetTree)
			categories.POST("", r.categoryHandler.Create)
			categories.GET("/:id", r.categoryHandler.GetByID)
			categories.PUT("/:id", r.categoryHandler.Update)
			categories.DELETE("/:id", r.categoryHandler.Delete)
		}

		// Units of Measure
		units := v1.Group("/units")
		{
			units.GET("", r.unitHandler.List)
			units.POST("", r.unitHandler.Create)
			units.GET("/:id", r.unitHandler.GetByID)
			units.POST("/convert", r.unitHandler.Convert)
			units.GET("/:id/conversions", r.unitHandler.GetConversions)
		}

		// Materials
		materials := v1.Group("/materials")
		{
			materials.GET("", r.materialHandler.List)
			materials.GET("/search", r.materialHandler.Search)
			materials.POST("", r.materialHandler.Create)
			materials.GET("/:id", r.materialHandler.GetByID)
			materials.PUT("/:id", r.materialHandler.Update)
			materials.DELETE("/:id", r.materialHandler.Delete)
			materials.POST("/:id/specifications", r.materialHandler.AddSpecification)
			materials.GET("/:id/specifications", r.materialHandler.GetSpecifications)
		}

		// Products
		products := v1.Group("/products")
		{
			products.GET("", r.productHandler.List)
			products.GET("/search", r.productHandler.Search)
			products.POST("", r.productHandler.Create)
			products.GET("/:id", r.productHandler.GetByID)
			products.PUT("/:id", r.productHandler.Update)
			products.DELETE("/:id", r.productHandler.Delete)
			products.POST("/:id/images", r.productHandler.AddImage)
			products.GET("/by-category/:category_id", r.productHandler.GetByCategory)
		}
	}

	return engine
}
