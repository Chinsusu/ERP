package http

import (
	"github.com/erp-cosmetics/reporting-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router handles HTTP routing
type Router struct {
	dashboardHandler *handler.DashboardHandler
	reportHandler    *handler.ReportHandler
	statsHandler     *handler.StatsHandler
	healthHandler    *handler.HealthHandler
	logger           *zap.Logger
}

// NewRouter creates new router
func NewRouter(
	dashboardHandler *handler.DashboardHandler,
	reportHandler *handler.ReportHandler,
	statsHandler *handler.StatsHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		dashboardHandler: dashboardHandler,
		reportHandler:    reportHandler,
		statsHandler:     statsHandler,
		healthHandler:    healthHandler,
		logger:           logger,
	}
}

// Setup sets up the HTTP router
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(r.logger))

	// Health check
	router.GET("/health", r.healthHandler.HealthCheck)
	router.GET("/ready", r.healthHandler.ReadyCheck)

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Dashboards
		dashboards := v1.Group("/dashboards")
		{
			dashboards.GET("", r.dashboardHandler.List)
			dashboards.POST("", r.dashboardHandler.Create)
			dashboards.GET("/default", r.dashboardHandler.GetDefault)
			dashboards.GET("/:id", r.dashboardHandler.Get)
			dashboards.PUT("/:id", r.dashboardHandler.Update)
			dashboards.DELETE("/:id", r.dashboardHandler.Delete)
			dashboards.POST("/:id/widgets", r.dashboardHandler.AddWidget)
		}

		// Widgets (for update/delete)
		widgets := v1.Group("/widgets")
		{
			widgets.PUT("/:id", r.dashboardHandler.UpdateWidget)
			widgets.DELETE("/:id", r.dashboardHandler.DeleteWidget)
		}

		// Reports
		reports := v1.Group("/reports")
		{
			reports.GET("", r.reportHandler.ListDefinitions)
			reports.GET("/:id", r.reportHandler.GetDefinition)
			reports.POST("/:id/execute", r.reportHandler.Execute)
			reports.GET("/:id/executions", r.reportHandler.ListExecutions)
			reports.GET("/executions/:id", r.reportHandler.GetExecution)
			reports.GET("/executions/:id/download", r.reportHandler.Download)
		}

		// Stats
		stats := v1.Group("/stats")
		{
			stats.GET("/dashboard", r.statsHandler.GetDashboard)
			stats.GET("/inventory", r.statsHandler.GetInventory)
			stats.GET("/sales", r.statsHandler.GetSales)
			stats.GET("/production", r.statsHandler.GetProduction)
			stats.GET("/procurement", r.statsHandler.GetProcurement)
		}
	}

	return router
}
