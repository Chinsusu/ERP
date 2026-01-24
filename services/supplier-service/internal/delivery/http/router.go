package http

import (
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/erp-cosmetics/supplier-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router holds all HTTP handlers
type Router struct {
	supplierHandler     *handler.SupplierHandler
	certificationHandler *handler.CertificationHandler
	evaluationHandler   *handler.EvaluationHandler
	healthHandler       *handler.HealthHandler
	logger              *zap.Logger
}

// NewRouter creates a new Router
func NewRouter(
	supplierHandler *handler.SupplierHandler,
	certificationHandler *handler.CertificationHandler,
	evaluationHandler *handler.EvaluationHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		supplierHandler:     supplierHandler,
		certificationHandler: certificationHandler,
		evaluationHandler:   evaluationHandler,
		healthHandler:       healthHandler,
		logger:              logger,
	}
}

// Setup configures and returns the Gin engine
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	engine := gin.New()

	// Global middleware
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger(r.logger))
	engine.Use(middleware.Recovery(r.logger))
	engine.Use(middleware.CORS("*"))

	// Health endpoints
	engine.GET("/health", r.healthHandler.Health)
	engine.GET("/ready", r.healthHandler.Ready)

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Suppliers
		suppliers := v1.Group("/suppliers")
		{
			suppliers.GET("", r.supplierHandler.List)
			suppliers.POST("", r.supplierHandler.Create)
			suppliers.GET("/:id", r.supplierHandler.Get)
			suppliers.PUT("/:id", r.supplierHandler.Update)
			suppliers.PATCH("/:id/approve", r.supplierHandler.Approve)
			suppliers.PATCH("/:id/block", r.supplierHandler.Block)

			// Addresses
			suppliers.GET("/:id/addresses", r.supplierHandler.ListAddresses)
			suppliers.POST("/:id/addresses", r.supplierHandler.CreateAddress)

			// Contacts
			suppliers.GET("/:id/contacts", r.supplierHandler.ListContacts)
			suppliers.POST("/:id/contacts", r.supplierHandler.CreateContact)

			// Certifications
			suppliers.GET("/:id/certifications", r.certificationHandler.ListCertifications)
			suppliers.POST("/:id/certifications", r.certificationHandler.AddCertification)

			// Evaluations
			suppliers.GET("/:id/evaluations", r.evaluationHandler.ListEvaluations)
			suppliers.POST("/:id/evaluations", r.evaluationHandler.CreateEvaluation)
		}

		// Certifications (global)
		certifications := v1.Group("/certifications")
		{
			certifications.GET("/expiring", r.certificationHandler.GetExpiringCertifications)
		}
	}

	return engine
}
