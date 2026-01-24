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
	EventSupplierCreated           = "supplier.created"
	EventSupplierApproved          = "supplier.approved"
	EventSupplierBlocked           = "supplier.blocked"
	EventSupplierUpdated           = "supplier.updated"
	EventCertificationAdded        = "supplier.certification.added"
	EventCertificationExpiring     = "supplier.certification.expiring"
	EventCertificationExpired      = "supplier.certification.expired"
	EventEvaluationCompleted       = "supplier.evaluation.completed"
)

// SupplierEvent represents a supplier-related event
type SupplierEvent struct {
	SupplierID   string    `json:"supplier_id"`
	SupplierCode string    `json:"supplier_code"`
	Name         string    `json:"name"`
	Status       string    `json:"status,omitempty"`
	Reason       string    `json:"reason,omitempty"`
	ActionBy     string    `json:"action_by,omitempty"`
	ActionAt     time.Time `json:"action_at"`
}

// CertificationEvent represents a certification-related event
type CertificationEvent struct {
	SupplierID        string    `json:"supplier_id"`
	SupplierName      string    `json:"supplier_name"`
	CertificationType string    `json:"certification_type"`
	CertNumber        string    `json:"certificate_number"`
	ExpiryDate        time.Time `json:"expiry_date"`
	DaysUntilExpiry   int       `json:"days_until_expiry,omitempty"`
}

// EvaluationEvent represents an evaluation-related event
type EvaluationEvent struct {
	SupplierID       string    `json:"supplier_id"`
	EvaluationID     string    `json:"evaluation_id"`
	EvaluationPeriod string    `json:"evaluation_period"`
	OverallScore     float64   `json:"overall_score"`
	EvaluatedBy      string    `json:"evaluated_by"`
	EvaluatedAt      time.Time `json:"evaluated_at"`
}

// PublishSupplierCreated publishes a supplier created event
func (p *Publisher) PublishSupplierCreated(ctx context.Context, event *SupplierEvent) error {
	return p.publish(ctx, EventSupplierCreated, event)
}

// PublishSupplierApproved publishes a supplier approved event
func (p *Publisher) PublishSupplierApproved(ctx context.Context, event *SupplierEvent) error {
	return p.publish(ctx, EventSupplierApproved, event)
}

// PublishSupplierBlocked publishes a supplier blocked event
func (p *Publisher) PublishSupplierBlocked(ctx context.Context, event *SupplierEvent) error {
	return p.publish(ctx, EventSupplierBlocked, event)
}

// PublishCertificationAdded publishes a certification added event
func (p *Publisher) PublishCertificationAdded(ctx context.Context, event *CertificationEvent) error {
	return p.publish(ctx, EventCertificationAdded, event)
}

// PublishCertificationExpiring publishes a certification expiring event
func (p *Publisher) PublishCertificationExpiring(ctx context.Context, event *CertificationEvent) error {
	return p.publish(ctx, EventCertificationExpiring, event)
}

// PublishCertificationExpired publishes a certification expired event
func (p *Publisher) PublishCertificationExpired(ctx context.Context, event *CertificationEvent) error {
	return p.publish(ctx, EventCertificationExpired, event)
}

// PublishEvaluationCompleted publishes an evaluation completed event
func (p *Publisher) PublishEvaluationCompleted(ctx context.Context, event *EvaluationEvent) error {
	return p.publish(ctx, EventEvaluationCompleted, event)
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
