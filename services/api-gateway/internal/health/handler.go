package health

import (
	"net/http"
	"time"

	"github.com/erp-cosmetics/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

// Handler handles health check requests
type Handler struct {
	registry  *proxy.ServiceRegistry
	startTime time.Time
	version   string
}

// NewHandler creates a new health handler
func NewHandler(registry *proxy.ServiceRegistry, version string) *Handler {
	return &Handler{
		registry:  registry,
		startTime: time.Now(),
		version:   version,
	}
}

// Health returns overall gateway health
func (h *Handler) Health(c *gin.Context) {
	uptime := int64(time.Since(h.startTime).Seconds())
	
	services := h.registry.GetAllStatuses()
	
	// Determine overall status
	status := "healthy"
	for _, svcStatus := range services {
		if svcStatus == "unhealthy" {
			status = "degraded"
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   status,
		"version":  h.version,
		"uptime":   uptime,
		"services": services,
	})
}

// Ready returns readiness status
func (h *Handler) Ready(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Live returns liveness status
func (h *Handler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
	})
}

// ServiceHealth returns health of a specific service
func (h *Handler) ServiceHealth(c *gin.Context) {
	serviceName := c.Param("service")
	
	status := h.registry.GetStatus(serviceName)
	
	if status == proxy.StatusUnknown {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Service not found",
			"service": serviceName,
		})
		return
	}

	httpStatus := http.StatusOK
	if status == proxy.StatusUnhealthy {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"service": serviceName,
		"status":  status.String(),
	})
}
