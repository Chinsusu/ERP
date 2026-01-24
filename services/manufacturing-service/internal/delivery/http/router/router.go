package router

import (
	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the HTTP router
func SetupRouter(
	bomHandler *handler.BOMHandler,
	woHandler *handler.WOHandler,
	qcHandler *handler.QCHandler,
	ncrHandler *handler.NCRHandler,
	traceHandler *handler.TraceHandler,
	healthHandler *handler.HealthHandler,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Health endpoints
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	r.GET("/live", healthHandler.Live)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// BOM routes
		boms := v1.Group("/boms")
		{
			boms.POST("", bomHandler.CreateBOM)
			boms.GET("", bomHandler.ListBOMs)
			boms.GET("/:id", bomHandler.GetBOM)
			boms.POST("/:id/approve", bomHandler.ApproveBOM)
		}

		// Work Order routes
		workOrders := v1.Group("/work-orders")
		{
			workOrders.POST("", woHandler.CreateWO)
			workOrders.GET("", woHandler.ListWOs)
			workOrders.GET("/:id", woHandler.GetWO)
			workOrders.PATCH("/:id/release", woHandler.ReleaseWO)
			workOrders.PATCH("/:id/start", woHandler.StartWO)
			workOrders.PATCH("/:id/complete", woHandler.CompleteWO)
		}

		// QC routes
		v1.GET("/qc-checkpoints", qcHandler.GetCheckpoints)
		
		qcInspections := v1.Group("/qc-inspections")
		{
			qcInspections.POST("", qcHandler.CreateInspection)
			qcInspections.GET("", qcHandler.ListInspections)
			qcInspections.GET("/:id", qcHandler.GetInspection)
			qcInspections.PATCH("/:id/approve", qcHandler.ApproveInspection)
		}

		// NCR routes
		ncrs := v1.Group("/ncrs")
		{
			ncrs.POST("", ncrHandler.CreateNCR)
			ncrs.GET("", ncrHandler.ListNCRs)
			ncrs.GET("/:id", ncrHandler.GetNCR)
			ncrs.PATCH("/:id/close", ncrHandler.CloseNCR)
		}

		// Traceability routes
		trace := v1.Group("/traceability")
		{
			trace.GET("/backward/:lot_id", traceHandler.TraceBackward)
			trace.GET("/forward/:lot_id", traceHandler.TraceForward)
		}
	}

	return r
}
