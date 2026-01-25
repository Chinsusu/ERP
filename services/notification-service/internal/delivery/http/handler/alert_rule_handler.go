package handler

import (
	"net/http"
	"strconv"

	alert_rule "github.com/erp-cosmetics/notification-service/internal/usecase/alert_rule"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AlertRuleHandler handles alert rule requests
type AlertRuleHandler struct {
	alertRuleUC alert_rule.UseCase
}

// NewAlertRuleHandler creates a new alert rule handler
func NewAlertRuleHandler(alertRuleUC alert_rule.UseCase) *AlertRuleHandler {
	return &AlertRuleHandler{
		alertRuleUC: alertRuleUC,
	}
}

// ListAlertRules handles GET /api/v1/alert-rules
func (h *AlertRuleHandler) ListAlertRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	output, err := h.alertRuleUC.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list alert rules", err.Error())
		return
	}

	response.SuccessWithPagination(c, http.StatusOK, "Alert rules retrieved", output.Rules, output.Total, page, pageSize)
}

// GetAlertRule handles GET /api/v1/alert-rules/:id
func (h *AlertRuleHandler) GetAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid alert rule ID", err.Error())
		return
	}

	rule, err := h.alertRuleUC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Alert rule not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Alert rule retrieved", rule)
}

// CreateAlertRule handles POST /api/v1/alert-rules
func (h *AlertRuleHandler) CreateAlertRule(c *gin.Context) {
	var input alert_rule.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Get user ID from context if available
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			input.CreatedBy = &uid
		}
	}

	rule, err := h.alertRuleUC.Create(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create alert rule", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Alert rule created", rule)
}

// UpdateAlertRule handles PUT /api/v1/alert-rules/:id
func (h *AlertRuleHandler) UpdateAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid alert rule ID", err.Error())
		return
	}

	var input alert_rule.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}
	input.ID = id

	rule, err := h.alertRuleUC.Update(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update alert rule", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Alert rule updated", rule)
}

// DeleteAlertRule handles DELETE /api/v1/alert-rules/:id
func (h *AlertRuleHandler) DeleteAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid alert rule ID", err.Error())
		return
	}

	if err := h.alertRuleUC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete alert rule", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Alert rule deleted", nil)
}

// ActivateAlertRule handles POST /api/v1/alert-rules/:id/activate
func (h *AlertRuleHandler) ActivateAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid alert rule ID", err.Error())
		return
	}

	if err := h.alertRuleUC.Activate(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to activate alert rule", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Alert rule activated", nil)
}

// DeactivateAlertRule handles POST /api/v1/alert-rules/:id/deactivate
func (h *AlertRuleHandler) DeactivateAlertRule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid alert rule ID", err.Error())
		return
	}

	if err := h.alertRuleUC.Deactivate(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to deactivate alert rule", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Alert rule deactivated", nil)
}
