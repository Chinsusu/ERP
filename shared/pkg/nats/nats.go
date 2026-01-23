package nats

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Client wraps a NATS connection
type Client struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	logger *zap.Logger
}

// Config holds NATS configuration
type Config struct {
	URL            string
	MaxReconnects  int
	ReconnectWait  time.Duration
	Logger         *zap.Logger
}

// NewClient creates a new NATS client
func NewClient(cfg *Config) (*Client, error) {
	opts := []nats.Option{
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				cfg.Logger.Error("NATS disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			cfg.Logger.Info("NATS reconnected", zap.String("url", nc.ConnectedUrl()))
		}),
	}

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	return &Client{
		conn:   conn,
		js:     js,
		logger: cfg.Logger,
	}, nil
}

// Publish publishes a message to a subject
func (c *Client) Publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if _, err := c.js.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	c.logger.Debug("Message published",
		zap.String("subject", subject),
		zap.Int("size", len(payload)),
	)

	return nil
}

// Subscribe subscribes to a subject with a handler
func (c *Client) Subscribe(subject, queue string, handler func(data []byte) error) (*nats.Subscription, error) {
	sub, err := c.js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		c.logger.Debug("Message received",
			zap.String("subject", msg.Subject),
			zap.Int("size", len(msg.Data)),
		)

		if err := handler(msg.Data); err != nil {
			c.logger.Error("Handler error",
				zap.String("subject", msg.Subject),
				zap.Error(err),
			)
			// Negative acknowledgment - message will be redelivered
			msg.Nak()
			return
		}

		// Acknowledge successful processing
		msg.Ack()
	}, nats.Durable(queue), nats.ManualAck())

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	c.logger.Info("Subscribed to subject",
		zap.String("subject", subject),
		zap.String("queue", queue),
	)

	return sub, nil
}

// CreateStream creates a JetStream stream
func (c *Client) CreateStream(name string, subjects []string) error {
	streamConfig := &nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
		Storage:  nats.FileStorage,
		Retention: nats.WorkQueuePolicy,
	}

	_, err := c.js.AddStream(streamConfig)
	if err != nil {
		return fmt.Errorf("failed to create stream: %w", err)
	}

	c.logger.Info("Stream created",
		zap.String("name", name),
		zap.Strings("subjects", subjects),
	)

	return nil
}

// Close closes the NATS connection
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.logger.Info("NATS connection closed")
	}
}

// Unmarshal unmarshals message data into a struct
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
