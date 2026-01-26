package handler

import (
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
)


type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "healthy",
		"service": "user-service",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "ready",
		"service": "user-service",
	})
}
