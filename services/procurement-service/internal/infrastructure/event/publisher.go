package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/erp-cosmetics/shared/pkg/nats"
	"go.uber.org/zap"
)

// Publisher handles publishing events to NATS
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

// Event types
const (
	EventPRCreated    = "procurement.pr.created"
	EventPRSubmitted  = "procurement.pr.submitted"
	EventPRApproved   = "procurement.pr.approved"
	EventPRRejected   = "procurement.pr.rejected"
	EventPOCreated    = "procurement.po.created"
	EventPOConfirmed  = "procurement.po.confirmed"
	EventPOReceived   = "procurement.po.received"
	EventPOClosed     = "procurement.po.closed"
	EventPOCancelled  = "procurement.po.cancelled"
)

// PREvent represents a PR-related event
type PREvent struct {
	PRID        string    `json:"pr_id"`
	PRNumber    string    `json:"pr_number"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
	RequesterID string    `json:"requester_id,omitempty"`
	ApprovedBy  string    `json:"approved_by,omitempty"`
	RejectedBy  string    `json:"rejected_by,omitempty"`
	Reason      string    `json:"reason,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// POEvent represents a PO-related event
type POEvent struct {
	POID                 string    `json:"po_id"`
	PONumber             string    `json:"po_number"`
	Status               string    `json:"status"`
	SupplierID           string    `json:"supplier_id"`
	SupplierName         string    `json:"supplier_name"`
	TotalAmount          float64   `json:"total_amount"`
	ExpectedDeliveryDate string    `json:"expected_delivery_date,omitempty"`
	Timestamp            time.Time `json:"timestamp"`
}

// PublishPRCreated publishes a PR created event
func (p *Publisher) PublishPRCreated(ctx context.Context, event *PREvent) error {
	return p.publish(ctx, EventPRCreated, event)
}

// PublishPRSubmitted publishes a PR submitted event
func (p *Publisher) PublishPRSubmitted(ctx context.Context, event *PREvent) error {
	return p.publish(ctx, EventPRSubmitted, event)
}

// PublishPRApproved publishes a PR approved event
func (p *Publisher) PublishPRApproved(ctx context.Context, event *PREvent) error {
	return p.publish(ctx, EventPRApproved, event)
}

// PublishPRRejected publishes a PR rejected event
func (p *Publisher) PublishPRRejected(ctx context.Context, event *PREvent) error {
	return p.publish(ctx, EventPRRejected, event)
}

// PublishPOCreated publishes a PO created event
func (p *Publisher) PublishPOCreated(ctx context.Context, event *POEvent) error {
	return p.publish(ctx, EventPOCreated, event)
}

// PublishPOConfirmed publishes a PO confirmed event
func (p *Publisher) PublishPOConfirmed(ctx context.Context, event *POEvent) error {
	return p.publish(ctx, EventPOConfirmed, event)
}

// PublishPOReceived publishes a PO received event
func (p *Publisher) PublishPOReceived(ctx context.Context, event *POEvent) error {
	return p.publish(ctx, EventPOReceived, event)
}

// PublishPOClosed publishes a PO closed event
func (p *Publisher) PublishPOClosed(ctx context.Context, event *POEvent) error {
	return p.publish(ctx, EventPOClosed, event)
}

func (p *Publisher) publish(ctx context.Context, subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.String("subject", subject), zap.Error(err))
		return err
	}

	if err := p.client.Publish(subject, payload); err != nil {
		p.logger.Error("Failed to publish event", zap.String("subject", subject), zap.Error(err))
		return err
	}

	p.logger.Info("Event published", zap.String("subject", subject))
	return nil
}
