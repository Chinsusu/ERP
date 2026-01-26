package handler

import (
	"net/http"
	"strconv"

	"github.com/erp-cosmetics/notification-service/internal/usecase/notification"
	"github.com/erp-cosmetics/notification-service/internal/usecase/template"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/erp-cosmetics/shared/pkg/errors")

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	notificationUC notification.UseCase
	templateUC     template.UseCase
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationUC notification.UseCase, templateUC template.UseCase) *NotificationHandler {
	return &NotificationHandler{
		notificationUC: notificationUC,
		templateUC:     templateUC,
	}
}

// SendNotification handles POST /api/v1/notifications/send
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var input notification.SendNotificationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}

	output, err := h.notificationUC.Send(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to send notification", http.StatusInternalServerError))
		return
	}

	response.Success(c, output)
}

// ListTemplates handles GET /api/v1/notifications/templates
func (h *NotificationHandler) ListTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	output, err := h.templateUC.List(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to list templates", http.StatusInternalServerError))
		return
	}

	response.SuccessWithMeta(c, output.Templates, response.NewMeta(page, pageSize, output.Total))
}

// GetTemplate handles GET /api/v1/notifications/templates/:id
func (h *NotificationHandler) GetTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid template ID", http.StatusBadRequest))
		return
	}

	tmpl, err := h.templateUC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Template not found", http.StatusNotFound))
		return
	}

	response.Success(c, tmpl)
}

// CreateTemplate handles POST /api/v1/notifications/templates
func (h *NotificationHandler) CreateTemplate(c *gin.Context) {
	var input template.CreateTemplateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}

	// Get user ID from context if available
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			input.CreatedBy = &uid
		}
	}

	tmpl, err := h.templateUC.Create(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to create template", http.StatusInternalServerError))
		return
	}

	response.Success(c, tmpl)
}

// UpdateTemplate handles PUT /api/v1/notifications/templates/:id
func (h *NotificationHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid template ID", http.StatusBadRequest))
		return
	}

	var input template.UpdateTemplateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, errors.New("ERROR", "Invalid request body", http.StatusBadRequest))
		return
	}
	input.ID = id

	tmpl, err := h.templateUC.Update(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to update template", http.StatusInternalServerError))
		return
	}

	response.Success(c, tmpl)
}

// DeleteTemplate handles DELETE /api/v1/notifications/templates/:id
func (h *NotificationHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid template ID", http.StatusBadRequest))
		return
	}

	if err := h.templateUC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, errors.New("ERROR", "Failed to delete template", http.StatusInternalServerError))
		return
	}

	response.Success(c, nil)
}
