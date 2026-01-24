package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Publisher handles NATS event publishing
type Publisher struct {
	js     nats.JetStreamContext
	logger *zap.Logger
}

// NewPublisher creates a new event publisher
func NewPublisher(nc *nats.Conn, logger *zap.Logger) (*Publisher, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Ensure MARKETING stream exists
	_, err = js.StreamInfo("MARKETING")
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "MARKETING",
			Subjects: []string{"marketing.>"},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return nil, err
		}
	}

	return &Publisher{js: js, logger: logger}, nil
}

// Event represents a generic event
type Event struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

func (p *Publisher) publish(ctx context.Context, subject string, data interface{}) error {
	event := Event{
		Type:      subject,
		Timestamp: time.Now(),
		Data:      data,
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(subject, bytes)
	if err != nil {
		p.logger.Error("Failed to publish event", zap.String("subject", subject), zap.Error(err))
		return err
	}

	p.logger.Info("Event published", zap.String("subject", subject))
	return nil
}

// PublishCampaignCreated publishes campaign.created event
func (p *Publisher) PublishCampaignCreated(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.campaign.created", data)
}

// PublishCampaignLaunched publishes campaign.launched event
func (p *Publisher) PublishCampaignLaunched(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.campaign.launched", data)
}

// PublishSampleRequestCreated publishes sample_request.created event
func (p *Publisher) PublishSampleRequestCreated(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.sample_request.created", data)
}

// PublishSampleRequestApproved publishes sample_request.approved event
func (p *Publisher) PublishSampleRequestApproved(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.sample_request.approved", data)
}

// PublishSampleShipped publishes sample.shipped event
func (p *Publisher) PublishSampleShipped(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.sample.shipped", data)
}

// PublishKOLPostRecorded publishes kol_post.recorded event
func (p *Publisher) PublishKOLPostRecorded(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.kol_post.recorded", data)
}

// PublishKOLCreated publishes kol.created event
func (p *Publisher) PublishKOLCreated(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.kol.created", data)
}

// PublishCollaborationCreated publishes collaboration.created event
func (p *Publisher) PublishCollaborationCreated(ctx context.Context, data interface{}) error {
	return p.publish(ctx, "marketing.collaboration.created", data)
}
