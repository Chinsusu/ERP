package http

import (
	"github.com/erp-cosmetics/sales-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the HTTP router
func NewRouter(
	customerHandler *handler.CustomerHandler,
	quotationHandler *handler.QuotationHandler,
	salesOrderHandler *handler.SalesOrderHandler,
	shipmentHandler *handler.ShipmentHandler,
) *gin.Engine {
	router := gin.New()

	// Global middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.CORS("*"))
	router.Use(middleware.RequestID())
	router.Use(gin.Logger())

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "sales-service"})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})
	router.GET("/live", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "live"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Customer Groups
		v1.GET("/customer-groups", customerHandler.ListGroups)
		v1.GET("/customer-groups/:id", customerHandler.GetGroup)

		// Customers
		customers := v1.Group("/customers")
		{
			customers.GET("", customerHandler.ListCustomers)
			customers.POST("", customerHandler.CreateCustomer)
			customers.GET("/:id", customerHandler.GetCustomer)
			customers.PUT("/:id", customerHandler.UpdateCustomer)
			customers.DELETE("/:id", customerHandler.DeleteCustomer)
			customers.GET("/:id/addresses", customerHandler.GetAddresses)
			customers.POST("/:id/addresses", customerHandler.CreateAddress)
			customers.GET("/:id/contacts", customerHandler.GetContacts)
			customers.POST("/:id/contacts", customerHandler.CreateContact)
			customers.GET("/:id/credit-check", customerHandler.CheckCredit)
		}

		// Quotations
		quotations := v1.Group("/quotations")
		{
			quotations.GET("", quotationHandler.ListQuotations)
			quotations.POST("", quotationHandler.CreateQuotation)
			quotations.GET("/:id", quotationHandler.GetQuotation)
			quotations.PUT("/:id", quotationHandler.UpdateQuotation)
			quotations.PATCH("/:id/send", quotationHandler.SendQuotation)
			quotations.POST("/:id/convert-to-order", quotationHandler.ConvertToOrder)
		}

		// Sales Orders
		orders := v1.Group("/sales-orders")
		{
			orders.GET("", salesOrderHandler.ListOrders)
			orders.POST("", salesOrderHandler.CreateOrder)
			orders.GET("/:id", salesOrderHandler.GetOrder)
			orders.PUT("/:id", salesOrderHandler.UpdateOrder)
			orders.PATCH("/:id/confirm", salesOrderHandler.ConfirmOrder)
			orders.PATCH("/:id/ship", salesOrderHandler.ShipOrder)
			orders.PATCH("/:id/deliver", salesOrderHandler.DeliverOrder)
			orders.PATCH("/:id/cancel", salesOrderHandler.CancelOrder)
		}

		// Shipments
		shipments := v1.Group("/shipments")
		{
			shipments.GET("", shipmentHandler.ListShipments)
			shipments.POST("", shipmentHandler.CreateShipment)
			shipments.GET("/:id", shipmentHandler.GetShipment)
			shipments.PATCH("/:id/ship", shipmentHandler.ShipShipment)
			shipments.PATCH("/:id/deliver", shipmentHandler.DeliverShipment)
		}
	}

	return router
}
