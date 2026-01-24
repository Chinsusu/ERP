package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health handles GET /health
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "wms-service",
	})
}

// Ready handles GET /ready
func (h *HealthHandler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Live handles GET /live
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}
