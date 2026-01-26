package handler

import (
	"net/http"

	"github.com/erp-cosmetics/reporting-service/internal/usecase/stats"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"

	"github.com/erp-cosmetics/shared/pkg/errors")

// StatsHandler handles stats requests
type StatsHandler struct {
	statsUC stats.UseCase
}

// NewStatsHandler creates new handler
func NewStatsHandler(uc stats.UseCase) *StatsHandler {
	return &StatsHandler{statsUC: uc}
}

// GetInventory handles GET /api/v1/stats/inventory
func (h *StatsHandler) GetInventory(c *gin.Context) {
	stats, err := h.statsUC.GetInventoryStats(c.Request.Context())
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to get inventory stats", http.StatusInternalServerError))
		return
	}

	response.Success(c, stats)
}

// GetSales handles GET /api/v1/stats/sales
func (h *StatsHandler) GetSales(c *gin.Context) {
	stats, err := h.statsUC.GetSalesStats(c.Request.Context())
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to get sales stats", http.StatusInternalServerError))
		return
	}

	response.Success(c, stats)
}

// GetProduction handles GET /api/v1/stats/production
func (h *StatsHandler) GetProduction(c *gin.Context) {
	stats, err := h.statsUC.GetProductionStats(c.Request.Context())
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to get production stats", http.StatusInternalServerError))
		return
	}

	response.Success(c, stats)
}

// GetProcurement handles GET /api/v1/stats/procurement
func (h *StatsHandler) GetProcurement(c *gin.Context) {
	stats, err := h.statsUC.GetProcurementStats(c.Request.Context())
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to get procurement stats", http.StatusInternalServerError))
		return
	}

	response.Success(c, stats)
}

// GetDashboard handles GET /api/v1/stats/dashboard
func (h *StatsHandler) GetDashboard(c *gin.Context) {
	stats, err := h.statsUC.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to get dashboard stats", http.StatusInternalServerError))
		return
	}

	response.Success(c, stats)
}
