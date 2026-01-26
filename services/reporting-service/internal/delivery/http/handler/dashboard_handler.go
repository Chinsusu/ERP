package handler

import (
	"net/http"
	"strconv"

	"github.com/erp-cosmetics/reporting-service/internal/usecase/dashboard"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/erp-cosmetics/shared/pkg/errors")

// DashboardHandler handles dashboard requests
type DashboardHandler struct {
	dashboardUC dashboard.UseCase
}

// NewDashboardHandler creates new handler
func NewDashboardHandler(uc dashboard.UseCase) *DashboardHandler {
	return &DashboardHandler{dashboardUC: uc}
}

// List handles GET /api/v1/dashboards
func (h *DashboardHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	output, err := h.dashboardUC.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to list dashboards", http.StatusInternalServerError))
		return
	}

	response.SuccessWithMeta(c, output.Dashboards, response.NewMeta(page, pageSize, output.Total))
}

// Get handles GET /api/v1/dashboards/:id
func (h *DashboardHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid dashboard ID", http.StatusBadRequest))
		return
	}

	dashboard, err := h.dashboardUC.GetWithWidgets(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Dashboard not found", http.StatusNotFound))
		return
	}

	response.Success(c, dashboard)
}

// GetDefault handles GET /api/v1/dashboards/default
func (h *DashboardHandler) GetDefault(c *gin.Context) {
	dashboard, err := h.dashboardUC.GetDefault(c.Request.Context())
	if err != nil {
		response.Error(c, errors.New("ERROR", "Default dashboard not found", http.StatusNotFound))
		return
	}

	dashboardWithWidgets, _ := h.dashboardUC.GetWithWidgets(c.Request.Context(), dashboard.ID)
	response.Success(c, dashboardWithWidgets)
}

// Create handles POST /api/v1/dashboards
func (h *DashboardHandler) Create(c *gin.Context) {
	var input dashboard.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}

	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			input.CreatedBy = &uid
		}
	}

	result, err := h.dashboardUC.Create(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to create dashboard", http.StatusInternalServerError))
		return
	}

	response.Success(c, result)
}

// Update handles PUT /api/v1/dashboards/:id
func (h *DashboardHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid dashboard ID", http.StatusBadRequest))
		return
	}

	var input dashboard.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}
	input.ID = id

	result, err := h.dashboardUC.Update(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to update dashboard", http.StatusInternalServerError))
		return
	}

	response.Success(c, result)
}

// Delete handles DELETE /api/v1/dashboards/:id
func (h *DashboardHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid dashboard ID", http.StatusBadRequest))
		return
	}

	if err := h.dashboardUC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, errors.New("ERROR", "Failed to delete dashboard", http.StatusInternalServerError))
		return
	}

	response.Success(c, nil)
}

// AddWidget handles POST /api/v1/dashboards/:id/widgets
func (h *DashboardHandler) AddWidget(c *gin.Context) {
	dashboardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid dashboard ID", http.StatusBadRequest))
		return
	}

	var input dashboard.AddWidgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}

	widget, err := h.dashboardUC.AddWidget(c.Request.Context(), dashboardID, &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to add widget", http.StatusInternalServerError))
		return
	}

	response.Success(c, widget)
}

// UpdateWidget handles PUT /api/v1/widgets/:id
func (h *DashboardHandler) UpdateWidget(c *gin.Context) {
	widgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid widget ID", http.StatusBadRequest))
		return
	}

	var input dashboard.UpdateWidgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}

	widget, err := h.dashboardUC.UpdateWidget(c.Request.Context(), widgetID, &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to update widget", http.StatusInternalServerError))
		return
	}

	response.Success(c, widget)
}

// DeleteWidget handles DELETE /api/v1/widgets/:id
func (h *DashboardHandler) DeleteWidget(c *gin.Context) {
	widgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid widget ID", http.StatusBadRequest))
		return
	}

	if err := h.dashboardUC.RemoveWidget(c.Request.Context(), widgetID); err != nil {
		response.Error(c, errors.New("ERROR", "Failed to delete widget", http.StatusInternalServerError))
		return
	}

	response.Success(c, nil)
}
