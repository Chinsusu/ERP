package router

import (
	"github.com/erp-cosmetics/procurement-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router with all routes
func SetupRouter(
	prHandler *handler.PRHandler,
	poHandler *handler.POHandler,
	healthHandler *handler.HealthHandler,
) *gin.Engine {
	r := gin.Default()

	// Health endpoints
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Purchase Requisition routes
		prs := v1.Group("/purchase-requisitions")
		{
			prs.POST("", prHandler.Create)
			prs.GET("", prHandler.List)
			prs.GET("/:id", prHandler.Get)
			prs.POST("/:id/submit", prHandler.Submit)
			prs.POST("/:id/approve", prHandler.Approve)
			prs.POST("/:id/reject", prHandler.Reject)
			prs.POST("/:id/convert-to-po", poHandler.ConvertPRToPO)
		}

		// Purchase Order routes
		pos := v1.Group("/purchase-orders")
		{
			pos.GET("", poHandler.List)
			pos.GET("/:id", poHandler.Get)
			pos.POST("/:id/confirm", poHandler.Confirm)
			pos.POST("/:id/cancel", poHandler.Cancel)
			pos.POST("/:id/close", poHandler.Close)
			pos.GET("/:id/receipts", poHandler.GetReceipts)
		}
	}

	return r
}
