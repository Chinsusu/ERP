package http

import (
	"github.com/erp-cosmetics/notification-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router handles HTTP routing
type Router struct {
	notificationHandler     *handler.NotificationHandler
	userNotificationHandler *handler.UserNotificationHandler
	alertRuleHandler        *handler.AlertRuleHandler
	healthHandler           *handler.HealthHandler
	logger                  *zap.Logger
}

// NewRouter creates a new router
func NewRouter(
	notificationHandler *handler.NotificationHandler,
	userNotificationHandler *handler.UserNotificationHandler,
	alertRuleHandler *handler.AlertRuleHandler,
	healthHandler *handler.HealthHandler,
	logger *zap.Logger,
) *Router {
	return &Router{
		notificationHandler:     notificationHandler,
		userNotificationHandler: userNotificationHandler,
		alertRuleHandler:        alertRuleHandler,
		healthHandler:           healthHandler,
		logger:                  logger,
	}
}

// Setup sets up the HTTP router
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS("*"))
	router.Use(middleware.Logger(r.logger))

	// Health check endpoints
	router.GET("/health", r.healthHandler.HealthCheck)
	router.GET("/ready", r.healthHandler.ReadyCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Notification routes
		notifications := v1.Group("/notifications")
		{
			// Send notification
			notifications.POST("/send", r.notificationHandler.SendNotification)

			// Template management
			templates := notifications.Group("/templates")
			{
				templates.GET("", r.notificationHandler.ListTemplates)
				templates.POST("", r.notificationHandler.CreateTemplate)
				templates.GET("/:id", r.notificationHandler.GetTemplate)
				templates.PUT("/:id", r.notificationHandler.UpdateTemplate)
				templates.DELETE("/:id", r.notificationHandler.DeleteTemplate)
			}

			// In-app notifications
			inApp := notifications.Group("/in-app")
			{
				inApp.GET("", r.userNotificationHandler.ListNotifications)
				inApp.POST("", r.userNotificationHandler.CreateNotification)
				inApp.GET("/unread-count", r.userNotificationHandler.GetUnreadCount)
				inApp.PATCH("/read-all", r.userNotificationHandler.MarkAllAsRead)
				inApp.PATCH("/:id/read", r.userNotificationHandler.MarkAsRead)
				inApp.DELETE("/:id", r.userNotificationHandler.DeleteNotification)
			}
		}

		// Alert rules routes
		alertRules := v1.Group("/alert-rules")
		{
			alertRules.GET("", r.alertRuleHandler.ListAlertRules)
			alertRules.POST("", r.alertRuleHandler.CreateAlertRule)
			alertRules.GET("/:id", r.alertRuleHandler.GetAlertRule)
			alertRules.PUT("/:id", r.alertRuleHandler.UpdateAlertRule)
			alertRules.DELETE("/:id", r.alertRuleHandler.DeleteAlertRule)
			alertRules.POST("/:id/activate", r.alertRuleHandler.ActivateAlertRule)
			alertRules.POST("/:id/deactivate", r.alertRuleHandler.DeactivateAlertRule)
		}
	}

	return router
}
