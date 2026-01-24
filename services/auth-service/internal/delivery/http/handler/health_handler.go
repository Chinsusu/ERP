package handler

import (
	"net/http"

	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health returns service health status
// @Summary Health check
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "auth-service",
	})
}

// Ready returns service readiness status
// @Summary Readiness check
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ready",
		"service": "auth-service",
	})
}
