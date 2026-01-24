package router

import (
	"github.com/erp-cosmetics/wms-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the HTTP router
func SetupRouter(
	warehouseHandler *handler.WarehouseHandler,
	stockHandler *handler.StockHandler,
	lotHandler *handler.LotHandler,
	grnHandler *handler.GRNHandler,
	issueHandler *handler.GoodsIssueHandler,
	reservationHandler *handler.ReservationHandler,
	adjustmentHandler *handler.AdjustmentHandler,
	healthHandler *handler.HealthHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health endpoints
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)
	r.GET("/live", healthHandler.Live)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Warehouse endpoints
		warehouses := v1.Group("/warehouses")
		{
			warehouses.GET("", warehouseHandler.ListWarehouses)
			warehouses.GET("/:id", warehouseHandler.GetWarehouse)
			warehouses.GET("/:id/zones", warehouseHandler.GetZones)
		}

		// Zone endpoints
		zones := v1.Group("/zones")
		{
			zones.GET("/:id/locations", warehouseHandler.GetLocations)
		}

		// Stock endpoints
		stock := v1.Group("/stock")
		{
			stock.GET("", stockHandler.ListStock)
			stock.GET("/by-material/:id", stockHandler.GetStockByMaterial)
			stock.GET("/expiring", stockHandler.GetExpiringStock)
			stock.GET("/low-stock", stockHandler.GetLowStock)
			stock.GET("/availability/:material_id", reservationHandler.CheckAvailability)
		}

		// Lot endpoints
		lots := v1.Group("/lots")
		{
			lots.GET("", lotHandler.ListLots)
			lots.GET("/:id", lotHandler.GetLot)
			lots.GET("/:id/movements", lotHandler.GetLotMovements)
		}

		// GRN endpoints
		grn := v1.Group("/grn")
		{
			grn.POST("", grnHandler.CreateGRN)
			grn.GET("", grnHandler.ListGRNs)
			grn.GET("/:id", grnHandler.GetGRN)
			grn.PATCH("/:id/complete", grnHandler.CompleteGRN)
		}

		// Goods Issue endpoints
		goodsIssue := v1.Group("/goods-issue")
		{
			goodsIssue.POST("", issueHandler.CreateGoodsIssue)
			goodsIssue.GET("", issueHandler.ListGoodsIssues)
			goodsIssue.GET("/:id", issueHandler.GetGoodsIssue)
		}

		// Reservation endpoints
		reservations := v1.Group("/reservations")
		{
			reservations.POST("", reservationHandler.CreateReservation)
			reservations.DELETE("/:id", reservationHandler.ReleaseReservation)
		}

		// Adjustment endpoints
		adjustments := v1.Group("/adjustments")
		{
			adjustments.POST("", adjustmentHandler.CreateAdjustment)
		}

		// Transfer endpoints
		transfers := v1.Group("/transfers")
		{
			transfers.POST("", adjustmentHandler.TransferStock)
		}
	}

	return r
}
