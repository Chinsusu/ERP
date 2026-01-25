package handler

import (
	"net/http"

	"github.com/erp-cosmetics/reporting-service/internal/usecase/stats"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
)

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
		response.Error(c, http.StatusInternalServerError, "Failed to get inventory stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Inventory stats retrieved", stats)
}

// GetSales handles GET /api/v1/stats/sales
func (h *StatsHandler) GetSales(c *gin.Context) {
	stats, err := h.statsUC.GetSalesStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get sales stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Sales stats retrieved", stats)
}

// GetProduction handles GET /api/v1/stats/production
func (h *StatsHandler) GetProduction(c *gin.Context) {
	stats, err := h.statsUC.GetProductionStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get production stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Production stats retrieved", stats)
}

// GetProcurement handles GET /api/v1/stats/procurement
func (h *StatsHandler) GetProcurement(c *gin.Context) {
	stats, err := h.statsUC.GetProcurementStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get procurement stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Procurement stats retrieved", stats)
}

// GetDashboard handles GET /api/v1/stats/dashboard
func (h *StatsHandler) GetDashboard(c *gin.Context) {
	stats, err := h.statsUC.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get dashboard stats", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Dashboard stats retrieved", stats)
}
