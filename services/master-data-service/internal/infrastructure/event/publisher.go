package event

import (
	"github.com/erp-cosmetics/shared/pkg/nats"
	"go.uber.org/zap"
)

// Publisher handles event publishing to NATS
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
	EventMaterialCreated     = "master_data.material.created"
	EventMaterialUpdated     = "master_data.material.updated"
	EventMaterialDeactivated = "master_data.material.deactivated"
	EventProductCreated      = "master_data.product.created"
	EventProductUpdated      = "master_data.product.updated"
	EventProductDeactivated  = "master_data.product.deactivated"
	EventCategoryCreated     = "master_data.category.created"
	EventCategoryUpdated     = "master_data.category.updated"
	EventLicenseExpiring     = "master_data.product.license_expiring"
)

// Publish publishes an event to NATS
func (p *Publisher) Publish(subject string, data interface{}) error {
	if p.client == nil {
		p.logger.Warn("NATS client not available, skipping event publish", zap.String("subject", subject))
		return nil
	}
	
	if err := p.client.Publish(subject, data); err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return err
	}
	
	p.logger.Info("Event published",
		zap.String("subject", subject),
	)
	return nil
}

// MaterialCreatedEvent represents material creation event data
type MaterialCreatedEvent struct {
	MaterialID   string `json:"material_id"`
	MaterialCode string `json:"material_code"`
	Name         string `json:"name"`
	MaterialType string `json:"material_type"`
	CreatedBy    string `json:"created_by,omitempty"`
}

// ProductCreatedEvent represents product creation event data
type ProductCreatedEvent struct {
	ProductID   string `json:"product_id"`
	ProductCode string `json:"product_code"`
	SKU         string `json:"sku"`
	Name        string `json:"name"`
	CreatedBy   string `json:"created_by,omitempty"`
}

// CategoryCreatedEvent represents category creation event data
type CategoryCreatedEvent struct {
	CategoryID   string `json:"category_id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	CategoryType string `json:"category_type"`
}
