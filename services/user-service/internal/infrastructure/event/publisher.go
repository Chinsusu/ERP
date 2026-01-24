package event

import (
	"github.com/erp-cosmetics/shared/pkg/nats"
	"go.uber.org/zap"
)

// Publisher wraps NATS for event publishing
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

// Publish publishes an event to NATS
func (p *Publisher) Publish(subject string, data interface{}) error {
	if err := p.client.Publish(subject, data); err != nil {
		p.logger.Error("Failed to publish event",
			zap.String("subject", subject),
			zap.Error(err),
		)
		return err
	}
	
	p.logger.Debug("Event published",
		zap.String("subject", subject),
	)
	
	return nil
}
