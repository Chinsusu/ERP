package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/erp-cosmetics/notification-service/internal/domain/entity"
	"github.com/erp-cosmetics/notification-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/nats"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Subscribed event subjects from other services
const (
	// Auth events
	SubjectAuthPasswordChanged = "auth.user.password_changed"
	SubjectAuthAccountLocked   = "auth.user.account_locked"

	// Procurement events
	SubjectPRSubmitted = "procurement.pr.submitted"
	SubjectPOCreated   = "procurement.po.created"

	// WMS events
	SubjectStockLowAlert   = "wms.stock.low_stock_alert"
	SubjectLotExpiringSoon = "wms.lot.expiring_soon"

	// Supplier events
	SubjectCertExpiring = "supplier.certification.expiring"

	// Manufacturing events
	SubjectQCFailed = "manufacturing.qc.failed"

	// Sales events
	SubjectOrderConfirmed = "sales.order.confirmed"
)

// Subscriber handles event subscriptions
type Subscriber struct {
	client               *nats.Client
	logger               *zap.Logger
	templateRepo         repository.TemplateRepository
	userNotificationRepo repository.UserNotificationRepository
	notificationRepo     repository.NotificationRepository
}

// NewSubscriber creates a new event subscriber
func NewSubscriber(
	client *nats.Client,
	logger *zap.Logger,
	templateRepo repository.TemplateRepository,
	userNotificationRepo repository.UserNotificationRepository,
	notificationRepo repository.NotificationRepository,
) *Subscriber {
	return &Subscriber{
		client:               client,
		logger:               logger,
		templateRepo:         templateRepo,
		userNotificationRepo: userNotificationRepo,
		notificationRepo:     notificationRepo,
	}
}

// SubscribeAll subscribes to all relevant events
func (s *Subscriber) SubscribeAll() error {
	subscriptions := []struct {
		subject string
		handler func([]byte) error
	}{
		{SubjectStockLowAlert, s.handleStockLowAlert},
		{SubjectLotExpiringSoon, s.handleLotExpiringSoon},
		{SubjectCertExpiring, s.handleCertExpiring},
		{SubjectPRSubmitted, s.handlePRSubmitted},
		{SubjectPOCreated, s.handlePOCreated},
		{SubjectQCFailed, s.handleQCFailed},
		{SubjectOrderConfirmed, s.handleOrderConfirmed},
	}

	for _, sub := range subscriptions {
		subject := sub.subject
		handler := sub.handler

		_, err := s.client.Subscribe(subject, func(msg []byte) {
			s.logger.Info("Received event", zap.String("subject", subject))
			if err := handler(msg); err != nil {
				s.logger.Error("Failed to handle event",
					zap.String("subject", subject),
					zap.Error(err),
				)
			}
		})
		if err != nil {
			s.logger.Error("Failed to subscribe to event",
				zap.String("subject", subject),
				zap.Error(err),
			)
			return err
		}
		s.logger.Info("Subscribed to event", zap.String("subject", subject))
	}

	return nil
}

// Event data structures
type StockLowAlertData struct {
	MaterialID      string  `json:"material_id"`
	MaterialCode    string  `json:"material_code"`
	MaterialName    string  `json:"material_name"`
	CurrentQuantity float64 `json:"current_quantity"`
	ReorderPoint    float64 `json:"reorder_point"`
	UoM             string  `json:"uom"`
	WarehouseID     string  `json:"warehouse_id"`
	WarehouseName   string  `json:"warehouse_name"`
}

type LotExpiringData struct {
	LotID         string  `json:"lot_id"`
	LotNumber     string  `json:"lot_number"`
	MaterialID    string  `json:"material_id"`
	MaterialCode  string  `json:"material_code"`
	MaterialName  string  `json:"material_name"`
	ExpiryDate    string  `json:"expiry_date"`
	DaysRemaining int     `json:"days_remaining"`
	Quantity      float64 `json:"quantity"`
	UoM           string  `json:"uom"`
	Location      string  `json:"location"`
}

type CertExpiringData struct {
	SupplierID        string `json:"supplier_id"`
	SupplierName      string `json:"supplier_name"`
	CertificateType   string `json:"certificate_type"`
	CertificateNumber string `json:"certificate_number"`
	ExpiryDate        string `json:"expiry_date"`
	DaysRemaining     int    `json:"days_remaining"`
}

type PRSubmittedData struct {
	PRID          string   `json:"pr_id"`
	PRNumber      string   `json:"pr_number"`
	RequesterID   string   `json:"requester_id"`
	RequesterName string   `json:"requester_name"`
	Department    string   `json:"department"`
	TotalAmount   float64  `json:"total_amount"`
	Currency      string   `json:"currency"`
	ApproverIDs   []string `json:"approver_ids"`
}

type QCFailedData struct {
	WorkOrderID     string `json:"work_order_id"`
	WorkOrderNumber string `json:"work_order_number"`
	ProductName     string `json:"product_name"`
	BatchNumber     string `json:"batch_number"`
	QCType          string `json:"qc_type"`
	FailedItems     string `json:"failed_items"`
	InspectorName   string `json:"inspector_name"`
	ManagerID       string `json:"manager_id"`
}

func (s *Subscriber) handleStockLowAlert(msg []byte) error {
	var data StockLowAlertData
	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	shortage := data.ReorderPoint - data.CurrentQuantity

	// Create in-app notification
	notification := &entity.UserNotification{
		Title:            "Low Stock Alert",
		Message:          formatLowStockMessage(data),
		NotificationType: entity.UserNotifTypeWarning,
		Category:         entity.CategoryAlert,
		LinkURL:          "/warehouse/stock?material=" + data.MaterialID,
		EntityType:       "MATERIAL",
	}

	if materialUUID, err := uuid.Parse(data.MaterialID); err == nil {
		notification.EntityID = &materialUUID
	}

	s.logger.Info("Stock low alert processed",
		zap.String("material", data.MaterialCode),
		zap.Float64("shortage", shortage),
	)

	return nil
}

func (s *Subscriber) handleLotExpiringSoon(msg []byte) error {
	var data LotExpiringData
	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	s.logger.Info("Lot expiring soon alert processed",
		zap.String("lot", data.LotNumber),
		zap.Int("days_remaining", data.DaysRemaining),
	)

	return nil
}

func (s *Subscriber) handleCertExpiring(msg []byte) error {
	var data CertExpiringData
	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	s.logger.Info("Certificate expiring alert processed",
		zap.String("supplier", data.SupplierName),
		zap.String("cert_type", data.CertificateType),
		zap.Int("days_remaining", data.DaysRemaining),
	)

	return nil
}

func (s *Subscriber) handlePRSubmitted(msg []byte) error {
	var data PRSubmittedData
	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	// Create notifications for each approver
	for _, approverID := range data.ApproverIDs {
		userID, err := uuid.Parse(approverID)
		if err != nil {
			continue
		}

		notification := &entity.UserNotification{
			UserID:           userID,
			Title:            "Purchase Requisition Pending Approval",
			Message:          formatPRMessage(data),
			NotificationType: entity.UserNotifTypeWarning,
			Category:         entity.CategoryApproval,
			LinkURL:          "/procurement/requisitions/" + data.PRID,
			EntityType:       "PR",
		}

		if prUUID, err := uuid.Parse(data.PRID); err == nil {
			notification.EntityID = &prUUID
		}

		if err := s.userNotificationRepo.Create(context.Background(), notification); err != nil {
			s.logger.Error("Failed to create PR approval notification",
				zap.String("approver_id", approverID),
				zap.Error(err),
			)
		}
	}

	s.logger.Info("PR submitted notification processed",
		zap.String("pr_number", data.PRNumber),
		zap.Int("approvers", len(data.ApproverIDs)),
	)

	return nil
}

func (s *Subscriber) handlePOCreated(msg []byte) error {
	var eventData map[string]interface{}
	if err := json.Unmarshal(msg, &eventData); err != nil {
		return err
	}

	s.logger.Info("PO created notification processed", zap.Any("data", eventData))
	return nil
}

func (s *Subscriber) handleQCFailed(msg []byte) error {
	var data QCFailedData
	if err := json.Unmarshal(msg, &data); err != nil {
		return err
	}

	// Notify production manager
	if data.ManagerID != "" {
		userID, err := uuid.Parse(data.ManagerID)
		if err == nil {
			notification := &entity.UserNotification{
				UserID:           userID,
				Title:            "Quality Check Failed",
				Message:          formatQCFailedMessage(data),
				NotificationType: entity.UserNotifTypeError,
				Category:         entity.CategoryAlert,
				LinkURL:          "/manufacturing/work-orders/" + data.WorkOrderID,
				EntityType:       "WO",
			}

			if woUUID, err := uuid.Parse(data.WorkOrderID); err == nil {
				notification.EntityID = &woUUID
			}

			if err := s.userNotificationRepo.Create(context.Background(), notification); err != nil {
				s.logger.Error("Failed to create QC failed notification", zap.Error(err))
			}
		}
	}

	s.logger.Info("QC failed notification processed",
		zap.String("work_order", data.WorkOrderNumber),
		zap.String("qc_type", data.QCType),
	)

	return nil
}

func (s *Subscriber) handleOrderConfirmed(msg []byte) error {
	var eventData map[string]interface{}
	if err := json.Unmarshal(msg, &eventData); err != nil {
		return err
	}

	s.logger.Info("Order confirmed notification processed", zap.Any("data", eventData))
	return nil
}

// Helper functions
func formatLowStockMessage(data StockLowAlertData) string {
	shortage := data.ReorderPoint - data.CurrentQuantity
	return fmt.Sprintf(
		"%s (%s) is below reorder point. Current: %.2f %s, Reorder Point: %.2f %s, Shortage: %.2f %s",
		data.MaterialName, data.MaterialCode,
		data.CurrentQuantity, data.UoM,
		data.ReorderPoint, data.UoM,
		shortage, data.UoM,
	)
}

func formatPRMessage(data PRSubmittedData) string {
	return fmt.Sprintf(
		"PR %s from %s (%s) requires your approval. Total: %.0f %s",
		data.PRNumber, data.RequesterName, data.Department,
		data.TotalAmount, data.Currency,
	)
}

func formatQCFailedMessage(data QCFailedData) string {
	return fmt.Sprintf(
		"QC failed for Work Order %s (%s), Batch: %s. Failed items: %s. Inspector: %s",
		data.WorkOrderNumber, data.ProductName, data.BatchNumber,
		data.FailedItems, data.InspectorName,
	)
}
