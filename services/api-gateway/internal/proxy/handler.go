package proxy

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler handles reverse proxy requests
type Handler struct {
	registry *ServiceRegistry
	client   *http.Client
	logger   *zap.Logger
}

// NewHandler creates a new proxy handler
func NewHandler(registry *ServiceRegistry, logger *zap.Logger, timeout time.Duration) *Handler {
	return &Handler{
		registry: registry,
		client: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		logger: logger,
	}
}

// Proxy forwards request to target service
func (h *Handler) Proxy(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get service URL
		serviceURL, err := h.registry.GetServiceURL(serviceName)
		if err != nil {
			h.logger.Error("Service not found", zap.String("service", serviceName), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"error":   "Service not available",
				"code":    "SERVICE_NOT_FOUND",
				"service": serviceName,
			})
			return
		}

		// Build target URL
		targetURL, err := url.Parse(serviceURL)
		if err != nil {
			h.logger.Error("Invalid service URL", zap.String("url", serviceURL), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal configuration error",
			})
			return
		}

		// Set target path
		targetURL.Path = c.Request.URL.Path
		targetURL.RawQuery = c.Request.URL.RawQuery

		// Create proxy request
		proxyReq, err := http.NewRequestWithContext(
			c.Request.Context(),
			c.Request.Method,
			targetURL.String(),
			c.Request.Body,
		)
		if err != nil {
			h.logger.Error("Failed to create proxy request", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create request",
			})
			return
		}

		// Copy headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		// Add gateway headers
		proxyReq.Header.Set("X-Forwarded-For", c.ClientIP())
		proxyReq.Header.Set("X-Real-IP", c.ClientIP())
		proxyReq.Header.Set("X-Forwarded-Host", c.Request.Host)
		proxyReq.Header.Set("X-Forwarded-Proto", "http")

		// Forward request ID
		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
			proxyReq.Header.Set("X-Request-ID", requestID)
		}

		// Forward user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			proxyReq.Header.Set("X-User-ID", userID.(string))
		}

		start := time.Now()

		// Execute request
		resp, err := h.client.Do(proxyReq)
		if err != nil {
			h.logger.Error("Proxy request failed",
				zap.String("service", serviceName),
				zap.String("url", targetURL.String()),
				zap.Error(err),
			)
			
			// Mark service as unhealthy
			h.registry.MarkUnhealthy(serviceName)

			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"error":   "Backend service error",
				"code":    "BACKEND_ERROR",
				"service": serviceName,
			})
			return
		}
		defer resp.Body.Close()

		// Mark service as healthy on success
		if resp.StatusCode < 500 {
			h.registry.MarkHealthy(serviceName)
		}

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// Add gateway metadata headers
		c.Header("X-Response-Time", fmt.Sprintf("%dms", time.Since(start).Milliseconds()))
		c.Header("X-Service", serviceName)
		c.Header("X-Gateway-Version", "1.0.0")

		// Copy response body
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
