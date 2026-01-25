package handler

import (
	"net/http"
	"strconv"

	"github.com/erp-cosmetics/reporting-service/internal/usecase/dashboard"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		response.Error(c, http.StatusInternalServerError, "Failed to list dashboards", err.Error())
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, "Dashboards retrieved", output.Dashboards, output.Total, page, pageSize)
}

// Get handles GET /api/v1/dashboards/:id
func (h *DashboardHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid dashboard ID", err.Error())
		return
	}

	dashboard, err := h.dashboardUC.GetWithWidgets(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Dashboard not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Dashboard retrieved", dashboard)
}

// GetDefault handles GET /api/v1/dashboards/default
func (h *DashboardHandler) GetDefault(c *gin.Context) {
	dashboard, err := h.dashboardUC.GetDefault(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusNotFound, "Default dashboard not found", err.Error())
		return
	}

	dashboardWithWidgets, _ := h.dashboardUC.GetWithWidgets(c.Request.Context(), dashboard.ID)
	response.Success(c, http.StatusOK, "Default dashboard retrieved", dashboardWithWidgets)
}

// Create handles POST /api/v1/dashboards
func (h *DashboardHandler) Create(c *gin.Context) {
	var input dashboard.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			input.CreatedBy = &uid
		}
	}

	result, err := h.dashboardUC.Create(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create dashboard", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Dashboard created", result)
}

// Update handles PUT /api/v1/dashboards/:id
func (h *DashboardHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid dashboard ID", err.Error())
		return
	}

	var input dashboard.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}
	input.ID = id

	result, err := h.dashboardUC.Update(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update dashboard", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Dashboard updated", result)
}

// Delete handles DELETE /api/v1/dashboards/:id
func (h *DashboardHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid dashboard ID", err.Error())
		return
	}

	if err := h.dashboardUC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete dashboard", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Dashboard deleted", nil)
}

// AddWidget handles POST /api/v1/dashboards/:id/widgets
func (h *DashboardHandler) AddWidget(c *gin.Context) {
	dashboardID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid dashboard ID", err.Error())
		return
	}

	var input dashboard.AddWidgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	widget, err := h.dashboardUC.AddWidget(c.Request.Context(), dashboardID, &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to add widget", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Widget added", widget)
}

// UpdateWidget handles PUT /api/v1/widgets/:id
func (h *DashboardHandler) UpdateWidget(c *gin.Context) {
	widgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid widget ID", err.Error())
		return
	}

	var input dashboard.UpdateWidgetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	widget, err := h.dashboardUC.UpdateWidget(c.Request.Context(), widgetID, &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update widget", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Widget updated", widget)
}

// DeleteWidget handles DELETE /api/v1/widgets/:id
func (h *DashboardHandler) DeleteWidget(c *gin.Context) {
	widgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid widget ID", err.Error())
		return
	}

	if err := h.dashboardUC.RemoveWidget(c.Request.Context(), widgetID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete widget", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Widget deleted", nil)
}
