package event

import (
	"encoding/json"

	"github.com/erp-cosmetics/shared/pkg/nats"
	"go.uber.org/zap"
)

// Publisher handles event publishing
type Publisher struct {
	client *nats.Client
	logger *zap.Logger
}

// NewPublisher creates a new event publisher
func NewPublisher(client *nats.Client, logger *zap.Logger) *Publisher {
	return &Publisher{
		client: client,
		logger: logger,
	}
}

// Event subjects
const (
	SubjectGRNCreated      = "wms.grn.created"
	SubjectGRNCompleted    = "wms.grn.completed"
	SubjectStockReceived   = "wms.stock.received"
	SubjectStockIssued     = "wms.stock.issued"
	SubjectStockReserved   = "wms.stock.reserved"
	SubjectLowStockAlert   = "wms.stock.low_stock_alert"
	SubjectLotExpiringSoon = "wms.lot.expiring_soon"
	SubjectLotExpired      = "wms.lot.expired"
)

// GRNCreatedEvent represents GRN created event
type GRNCreatedEvent struct {
	GRNID     string `json:"grn_id"`
	GRNNumber string `json:"grn_number"`
	POID      string `json:"po_id,omitempty"`
}

// GRNCompletedEvent represents GRN completed event
type GRNCompletedEvent struct {
	GRNID       string                  `json:"grn_id"`
	GRNNumber   string                  `json:"grn_number"`
	POID        string                  `json:"po_id,omitempty"`
	WarehouseID string                  `json:"warehouse_id"`
	Items       []GRNCompletedEventItem `json:"items"`
}

// GRNCompletedEventItem represents an item in GRN completed event
type GRNCompletedEventItem struct {
	MaterialID   string  `json:"material_id"`
	LotID        string  `json:"lot_id"`
	LotNumber    string  `json:"lot_number"`
	ReceivedQty  float64 `json:"received_qty"`
	AcceptedQty  float64 `json:"accepted_qty"`
}

// StockReceivedEvent represents stock received event
type StockReceivedEvent struct {
	MaterialID  string  `json:"material_id"`
	LotID       string  `json:"lot_id"`
	Quantity    float64 `json:"quantity"`
	LocationID  string  `json:"location_id"`
	WarehouseID string  `json:"warehouse_id"`
}

// StockIssuedEvent represents stock issued event
type StockIssuedEvent struct {
	MaterialID    string           `json:"material_id"`
	Quantity      float64          `json:"quantity"`
	LotsUsed      []LotUsedInIssue `json:"lots_used"`
	ReferenceType string           `json:"reference_type"`
	ReferenceID   string           `json:"reference_id"`
}

// LotUsedInIssue represents a lot used in an issue
type LotUsedInIssue struct {
	LotID      string  `json:"lot_id"`
	LotNumber  string  `json:"lot_number"`
	Quantity   float64 `json:"quantity"`
	ExpiryDate string  `json:"expiry_date"`
}

// StockReservedEvent represents stock reserved event
type StockReservedEvent struct {
	ReservationID   string  `json:"reservation_id"`
	MaterialID      string  `json:"material_id"`
	Quantity        float64 `json:"quantity"`
	ReservationType string  `json:"reservation_type"`
	ReferenceID     string  `json:"reference_id"`
}

// LowStockAlertEvent represents low stock alert event
type LowStockAlertEvent struct {
	MaterialID      string  `json:"material_id"`
	MaterialCode    string  `json:"material_code"`
	CurrentQuantity float64 `json:"current_quantity"`
	ReorderPoint    float64 `json:"reorder_point"`
}

// LotExpiringEvent represents lot expiring event
type LotExpiringEvent struct {
	LotID           string  `json:"lot_id"`
	LotNumber       string  `json:"lot_number"`
	MaterialID      string  `json:"material_id"`
	ExpiryDate      string  `json:"expiry_date"`
	DaysUntilExpiry int     `json:"days_until_expiry"`
	Quantity        float64 `json:"quantity"`
}

// PublishGRNCreated publishes GRN created event
func (p *Publisher) PublishGRNCreated(event *GRNCreatedEvent) error {
	return p.publish(SubjectGRNCreated, event)
}

// PublishGRNCompleted publishes GRN completed event
func (p *Publisher) PublishGRNCompleted(event *GRNCompletedEvent) error {
	return p.publish(SubjectGRNCompleted, event)
}

// PublishStockReceived publishes stock received event
func (p *Publisher) PublishStockReceived(event *StockReceivedEvent) error {
	return p.publish(SubjectStockReceived, event)
}

// PublishStockIssued publishes stock issued event
func (p *Publisher) PublishStockIssued(event *StockIssuedEvent) error {
	return p.publish(SubjectStockIssued, event)
}

// PublishStockReserved publishes stock reserved event
func (p *Publisher) PublishStockReserved(event *StockReservedEvent) error {
	return p.publish(SubjectStockReserved, event)
}

// PublishLowStockAlert publishes low stock alert event
func (p *Publisher) PublishLowStockAlert(event *LowStockAlertEvent) error {
	return p.publish(SubjectLowStockAlert, event)
}

// PublishLotExpiringSoon publishes lot expiring soon event
func (p *Publisher) PublishLotExpiringSoon(event *LotExpiringEvent) error {
	return p.publish(SubjectLotExpiringSoon, event)
}

// PublishLotExpired publishes lot expired event
func (p *Publisher) PublishLotExpired(event *LotExpiringEvent) error {
	return p.publish(SubjectLotExpired, event)
}

func (p *Publisher) publish(subject string, data interface{}) error {
	if p.client == nil {
		p.logger.Warn("NATS client not available, skipping event publish",
			zap.String("subject", subject))
		return nil
	}

	payload, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("Failed to marshal event",
			zap.String("subject", subject),
			zap.Error(err))
		return err
	}

	if err := p.client.Publish(subject, payload); err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("subject", subject),
			zap.Error(err))
		return err
	}

	p.logger.Info("Event published",
		zap.String("subject", subject))
	return nil
}
