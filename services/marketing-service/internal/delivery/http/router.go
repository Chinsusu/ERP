package http

import (
	"github.com/erp-cosmetics/marketing-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// NewRouter creates and configures the HTTP router
func NewRouter(
	kolHandler *handler.KOLHandler,
	campaignHandler *handler.CampaignHandler,
	sampleHandler *handler.SampleHandler,
) *gin.Engine {
	router := gin.New()

	// Global middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.CORS("*"))
	router.Use(middleware.RequestID())
	router.Use(gin.Logger())

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"service": "marketing-service", "status": "healthy"})
	})
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"ready": true})
	})
	router.GET("/live", func(c *gin.Context) {
		c.JSON(200, gin.H{"live": true})
	})

	// API v1
	v1 := router.Group("/api/v1/marketing")
	{
		// KOL Tiers
		v1.GET("/kol-tiers", kolHandler.ListTiers)

		// KOLs
		kols := v1.Group("/kols")
		{
			kols.GET("", kolHandler.ListKOLs)
			kols.POST("", kolHandler.CreateKOL)
			kols.GET("/:id", kolHandler.GetKOL)
			kols.PUT("/:id", kolHandler.UpdateKOL)
			kols.DELETE("/:id", kolHandler.DeleteKOL)
			kols.GET("/:id/posts", kolHandler.GetKOLPosts)
		}

		// Campaigns
		campaigns := v1.Group("/campaigns")
		{
			campaigns.GET("", campaignHandler.ListCampaigns)
			campaigns.POST("", campaignHandler.CreateCampaign)
			campaigns.GET("/:id", campaignHandler.GetCampaign)
			campaigns.PUT("/:id", campaignHandler.UpdateCampaign)
			campaigns.PATCH("/:id/launch", campaignHandler.LaunchCampaign)
			campaigns.GET("/:id/collaborations", campaignHandler.GetCampaignCollaborations)
			campaigns.GET("/:id/performance", campaignHandler.GetCampaignPerformance)
		}

		// Sample Requests
		samples := v1.Group("/samples")
		{
			samples.GET("/requests", sampleHandler.ListSampleRequests)
			samples.POST("/requests", sampleHandler.CreateSampleRequest)
			samples.GET("/requests/:id", sampleHandler.GetSampleRequest)
			samples.PATCH("/requests/:id/approve", sampleHandler.ApproveSampleRequest)
			samples.PATCH("/requests/:id/ship", sampleHandler.ShipSampleRequest)
			samples.GET("/shipments", sampleHandler.ListShipments)
		}
	}

	return router
}
