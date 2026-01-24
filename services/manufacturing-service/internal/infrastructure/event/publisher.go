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
	SubjectBOMCreated     = "manufacturing.bom.created"
	SubjectBOMApproved    = "manufacturing.bom.approved"
	SubjectWOCreated      = "manufacturing.wo.created"
	SubjectWOReleased     = "manufacturing.wo.released"
	SubjectWOStarted      = "manufacturing.wo.started"
	SubjectWOCompleted    = "manufacturing.wo.completed"
	SubjectQCPassed       = "manufacturing.qc.passed"
	SubjectQCFailed       = "manufacturing.qc.failed"
	SubjectNCRCreated     = "manufacturing.ncr.created"
)

// BOMEvent represents a BOM event payload
type BOMEvent struct {
	BOMID     string `json:"bom_id"`
	BOMNumber string `json:"bom_number"`
	ProductID string `json:"product_id"`
	Version   int    `json:"version"`
	Status    string `json:"status"`
}

// WOEvent represents a work order event payload
type WOEvent struct {
	WOID            string  `json:"wo_id"`
	WONumber        string  `json:"wo_number"`
	ProductID       string  `json:"product_id"`
	BOMID           string  `json:"bom_id"`
	BatchNumber     string  `json:"batch_number"`
	PlannedQuantity float64 `json:"planned_quantity"`
	Status          string  `json:"status"`
}

// WOCompletedEvent represents a completed work order event
type WOCompletedEvent struct {
	WOID         string  `json:"wo_id"`
	WONumber     string  `json:"wo_number"`
	ProductID    string  `json:"product_id"`
	BatchNumber  string  `json:"batch_number"`
	GoodQuantity float64 `json:"good_quantity"`
	OutputLotID  string  `json:"output_lot_id,omitempty"`
}

// QCEvent represents a QC event payload
type QCEvent struct {
	InspectionID     string `json:"inspection_id"`
	InspectionNumber string `json:"inspection_number"`
	InspectionType   string `json:"inspection_type"`
	ReferenceType    string `json:"reference_type"`
	ReferenceID      string `json:"reference_id"`
	Result           string `json:"result"`
	LotID            string `json:"lot_id,omitempty"`
}

// NCREvent represents an NCR event payload
type NCREvent struct {
	NCRID       string `json:"ncr_id"`
	NCRNumber   string `json:"ncr_number"`
	NCType      string `json:"nc_type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	LotID       string `json:"lot_id,omitempty"`
}

// Publish publishes an event
func (p *Publisher) Publish(subject string, payload interface{}) error {
	if p.client == nil {
		p.logger.Warn("NATS client not available, skipping event publish",
			zap.String("subject", subject))
		return nil
	}

	data, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("Failed to marshal event payload",
			zap.String("subject", subject),
			zap.Error(err))
		return err
	}

	if err := p.client.Publish(subject, data); err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("subject", subject),
			zap.Error(err))
		return err
	}

	p.logger.Info("Event published",
		zap.String("subject", subject))
	return nil
}

// PublishBOMCreated publishes BOM created event
func (p *Publisher) PublishBOMCreated(event BOMEvent) error {
	return p.Publish(SubjectBOMCreated, event)
}

// PublishBOMApproved publishes BOM approved event
func (p *Publisher) PublishBOMApproved(event BOMEvent) error {
	return p.Publish(SubjectBOMApproved, event)
}

// PublishWOCreated publishes WO created event
func (p *Publisher) PublishWOCreated(event WOEvent) error {
	return p.Publish(SubjectWOCreated, event)
}

// PublishWOReleased publishes WO released event
func (p *Publisher) PublishWOReleased(event WOEvent) error {
	return p.Publish(SubjectWOReleased, event)
}

// PublishWOStarted publishes WO started event - triggers WMS material reservation
func (p *Publisher) PublishWOStarted(event WOEvent) error {
	return p.Publish(SubjectWOStarted, event)
}

// PublishWOCompleted publishes WO completed event - triggers WMS finished goods receipt
func (p *Publisher) PublishWOCompleted(event WOCompletedEvent) error {
	return p.Publish(SubjectWOCompleted, event)
}

// PublishQCPassed publishes QC passed event
func (p *Publisher) PublishQCPassed(event QCEvent) error {
	return p.Publish(SubjectQCPassed, event)
}

// PublishQCFailed publishes QC failed event - triggers notification
func (p *Publisher) PublishQCFailed(event QCEvent) error {
	return p.Publish(SubjectQCFailed, event)
}

// PublishNCRCreated publishes NCR created event - triggers notification
func (p *Publisher) PublishNCRCreated(event NCREvent) error {
	return p.Publish(SubjectNCRCreated, event)
}
