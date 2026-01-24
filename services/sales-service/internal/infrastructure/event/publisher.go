package event

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Publisher handles NATS event publishing
type Publisher struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger *zap.Logger
}

// NewPublisher creates a new event publisher
func NewPublisher(nc *nats.Conn, logger *zap.Logger) (*Publisher, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Ensure stream exists
	_, err = js.StreamInfo("SALES")
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "SALES",
			Subjects: []string{"sales.>"},
			Storage:  nats.FileStorage,
			MaxAge:   30 * 24 * time.Hour, // 30 days retention
		})
		if err != nil {
			return nil, err
		}
	}

	return &Publisher{
		nc:     nc,
		js:     js,
		logger: logger,
	}, nil
}

// publish publishes an event to NATS
func (p *Publisher) publish(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		p.logger.Error("Failed to marshal event", zap.String("subject", subject), zap.Error(err))
		return err
	}

	_, err = p.js.Publish(subject, payload)
	if err != nil {
		p.logger.Error("Failed to publish event", zap.String("subject", subject), zap.Error(err))
		return err
	}

	p.logger.Info("Event published", zap.String("subject", subject))
	return nil
}

// CustomerCreatedEvent represents customer created event
type CustomerCreatedEvent struct {
	CustomerID   string `json:"customer_id"`
	CustomerCode string `json:"customer_code"`
	Name         string `json:"name"`
	CustomerType string `json:"customer_type"`
	Email        string `json:"email"`
	Timestamp    string `json:"timestamp"`
}

// PublishCustomerCreated publishes customer created event
func (p *Publisher) PublishCustomerCreated(event *CustomerCreatedEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.customer.created", event)
}

// QuotationSentEvent represents quotation sent event
type QuotationSentEvent struct {
	QuotationID     string  `json:"quotation_id"`
	QuotationNumber string  `json:"quotation_number"`
	CustomerID      string  `json:"customer_id"`
	TotalAmount     float64 `json:"total_amount"`
	ValidUntil      string  `json:"valid_until"`
	Timestamp       string  `json:"timestamp"`
}

// PublishQuotationSent publishes quotation sent event
func (p *Publisher) PublishQuotationSent(event *QuotationSentEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.quotation.sent", event)
}

// OrderLineItem represents a line item in order event
type OrderLineItem struct {
	ProductID   string  `json:"product_id"`
	ProductCode string  `json:"product_code"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

// OrderCreatedEvent represents order created event
type OrderCreatedEvent struct {
	SOID        string          `json:"so_id"`
	SONumber    string          `json:"so_number"`
	CustomerID  string          `json:"customer_id"`
	TotalAmount float64         `json:"total_amount"`
	Items       []OrderLineItem `json:"items"`
	Timestamp   string          `json:"timestamp"`
}

// PublishOrderCreated publishes order created event
func (p *Publisher) PublishOrderCreated(event *OrderCreatedEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.order.created", event)
}

// OrderConfirmedEvent represents order confirmed event - WMS subscribes to this
type OrderConfirmedEvent struct {
	SOID           string          `json:"so_id"`
	SONumber       string          `json:"so_number"`
	CustomerID     string          `json:"customer_id"`
	DeliveryDate   string          `json:"delivery_date"`
	DeliveryAddress string         `json:"delivery_address"`
	TotalAmount    float64         `json:"total_amount"`
	Items          []OrderLineItem `json:"items"`
	Timestamp      string          `json:"timestamp"`
}

// PublishOrderConfirmed publishes order confirmed event
func (p *Publisher) PublishOrderConfirmed(event *OrderConfirmedEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.order.confirmed", event)
}

// OrderShippedEvent represents order shipped event
type OrderShippedEvent struct {
	SOID           string `json:"so_id"`
	SONumber       string `json:"so_number"`
	CustomerID     string `json:"customer_id"`
	ShipmentID     string `json:"shipment_id"`
	ShipmentNumber string `json:"shipment_number"`
	Carrier        string `json:"carrier"`
	TrackingNumber string `json:"tracking_number"`
	Timestamp      string `json:"timestamp"`
}

// PublishOrderShipped publishes order shipped event
func (p *Publisher) PublishOrderShipped(event *OrderShippedEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.order.shipped", event)
}

// OrderDeliveredEvent represents order delivered event
type OrderDeliveredEvent struct {
	SOID        string `json:"so_id"`
	SONumber    string `json:"so_number"`
	CustomerID  string `json:"customer_id"`
	DeliveredAt string `json:"delivered_at"`
	Timestamp   string `json:"timestamp"`
}

// PublishOrderDelivered publishes order delivered event
func (p *Publisher) PublishOrderDelivered(event *OrderDeliveredEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.order.delivered", event)
}

// OrderCancelledEvent represents order cancelled event - WMS subscribes to release reservations
type OrderCancelledEvent struct {
	SOID        string `json:"so_id"`
	SONumber    string `json:"so_number"`
	CustomerID  string `json:"customer_id"`
	Reason      string `json:"reason"`
	CancelledBy string `json:"cancelled_by"`
	Timestamp   string `json:"timestamp"`
}

// PublishOrderCancelled publishes order cancelled event
func (p *Publisher) PublishOrderCancelled(event *OrderCancelledEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.order.cancelled", event)
}

// ReturnCreatedEvent represents return created event
type ReturnCreatedEvent struct {
	ReturnID     string  `json:"return_id"`
	ReturnNumber string  `json:"return_number"`
	SOID         string  `json:"so_id"`
	CustomerID   string  `json:"customer_id"`
	ReturnType   string  `json:"return_type"`
	RefundAmount float64 `json:"refund_amount"`
	Timestamp    string  `json:"timestamp"`
}

// PublishReturnCreated publishes return created event
func (p *Publisher) PublishReturnCreated(event *ReturnCreatedEvent) {
	event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	p.publish("sales.return.created", event)
}
