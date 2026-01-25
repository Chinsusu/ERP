package handler

import (
	"net/http"
	"strconv"

	user_notification "github.com/erp-cosmetics/notification-service/internal/usecase/user_notification"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserNotificationHandler handles user notification requests
type UserNotificationHandler struct {
	userNotificationUC user_notification.UseCase
}

// NewUserNotificationHandler creates a new user notification handler
func NewUserNotificationHandler(userNotificationUC user_notification.UseCase) *UserNotificationHandler {
	return &UserNotificationHandler{
		userNotificationUC: userNotificationUC,
	}
}

// ListNotifications handles GET /api/v1/notifications/in-app
func (h *UserNotificationHandler) ListNotifications(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	var userID uuid.UUID
	switch v := userIDVal.(type) {
	case uuid.UUID:
		userID = v
	case string:
		var err error
		userID, err = uuid.Parse(v)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}
	default:
		response.Error(c, http.StatusBadRequest, "Invalid user ID type", "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	output, err := h.userNotificationUC.GetByUserID(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list notifications", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Notifications retrieved",
		"data":         output.Notifications,
		"total":        output.Total,
		"unread_count": output.UnreadCount,
		"page":         output.Page,
		"page_size":    output.PageSize,
	})
}

// CreateNotification handles POST /api/v1/notifications/in-app
func (h *UserNotificationHandler) CreateNotification(c *gin.Context) {
	var input user_notification.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	notification, err := h.userNotificationUC.Create(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create notification", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "Notification created", notification)
}

// MarkAsRead handles PATCH /api/v1/notifications/in-app/:id/read
func (h *UserNotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID", err.Error())
		return
	}

	if err := h.userNotificationUC.MarkAsRead(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to mark as read", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Notification marked as read", nil)
}

// MarkAllAsRead handles PATCH /api/v1/notifications/in-app/read-all
func (h *UserNotificationHandler) MarkAllAsRead(c *gin.Context) {
	// Get user ID from context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	var userID uuid.UUID
	switch v := userIDVal.(type) {
	case uuid.UUID:
		userID = v
	case string:
		var err error
		userID, err = uuid.Parse(v)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}
	default:
		response.Error(c, http.StatusBadRequest, "Invalid user ID type", "")
		return
	}

	if err := h.userNotificationUC.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to mark all as read", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "All notifications marked as read", nil)
}

// DeleteNotification handles DELETE /api/v1/notifications/in-app/:id
func (h *UserNotificationHandler) DeleteNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID", err.Error())
		return
	}

	if err := h.userNotificationUC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete notification", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Notification deleted", nil)
}

// GetUnreadCount handles GET /api/v1/notifications/in-app/unread-count
func (h *UserNotificationHandler) GetUnreadCount(c *gin.Context) {
	// Get user ID from context
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated", "")
		return
	}

	var userID uuid.UUID
	switch v := userIDVal.(type) {
	case uuid.UUID:
		userID = v
	case string:
		var err error
		userID, err = uuid.Parse(v)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}
	default:
		response.Error(c, http.StatusBadRequest, "Invalid user ID type", "")
		return
	}

	count, err := h.userNotificationUC.CountUnread(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get unread count", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Unread count retrieved", gin.H{"unread_count": count})
}
