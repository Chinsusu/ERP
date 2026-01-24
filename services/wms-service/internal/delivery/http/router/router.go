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

		// TODO: Add these endpoints in next iteration
		// Goods Issue endpoints
		// goodsIssue := v1.Group("/goods-issue")
		// Reserve endpoints
		// reserve := v1.Group("/reserve")
		// Adjustment endpoints
		// adjustment := v1.Group("/adjustment")
		// Transfer endpoints
		// transfer := v1.Group("/transfer")
	}

	return r
}
