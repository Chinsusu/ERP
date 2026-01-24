package entity

import (
	"time"

	"github.com/google/uuid"
)

// ShipmentStatus represents shipment status
type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "PENDING"
	ShipmentStatusPicked    ShipmentStatus = "PICKED"
	ShipmentStatusPacked    ShipmentStatus = "PACKED"
	ShipmentStatusShipped   ShipmentStatus = "SHIPPED"
	ShipmentStatusInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
	ShipmentStatusReturned  ShipmentStatus = "RETURNED"
)

// Shipment represents a shipment
type Shipment struct {
	ID                    uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ShipmentNumber        string          `json:"shipment_number" gorm:"type:varchar(20);unique;not null"`
	SalesOrderID          uuid.UUID       `json:"sales_order_id" gorm:"type:uuid;not null"`
	SalesOrder            *SalesOrder     `json:"sales_order,omitempty" gorm:"foreignKey:SalesOrderID"`
	ShippedDate           *time.Time      `json:"shipped_date" gorm:"type:date"`
	EstimatedDeliveryDate *time.Time      `json:"estimated_delivery_date" gorm:"type:date"`
	ActualDeliveryDate    *time.Time      `json:"actual_delivery_date" gorm:"type:date"`
	Carrier               string          `json:"carrier" gorm:"type:varchar(100)"`
	TrackingNumber        string          `json:"tracking_number" gorm:"type:varchar(100)"`
	ShippingMethod        string          `json:"shipping_method" gorm:"type:varchar(50)"`
	ShippingCost          float64         `json:"shipping_cost" gorm:"type:decimal(18,2);default:0"`
	Status                ShipmentStatus  `json:"status" gorm:"type:varchar(20);default:'PENDING'"`
	RecipientName         string          `json:"recipient_name" gorm:"type:varchar(200)"`
	RecipientPhone        string          `json:"recipient_phone" gorm:"type:varchar(20)"`
	DeliveryAddress       string          `json:"delivery_address" gorm:"type:text"`
	Notes                 string          `json:"notes" gorm:"type:text"`
	CreatedBy             *uuid.UUID      `json:"created_by" gorm:"type:uuid"`
	UpdatedBy             *uuid.UUID      `json:"updated_by" gorm:"type:uuid"`
	CreatedAt             time.Time       `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time       `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (Shipment) TableName() string {
	return "shipments"
}

// MarkPicked marks shipment as picked
func (s *Shipment) MarkPicked() {
	s.Status = ShipmentStatusPicked
	s.UpdatedAt = time.Now()
}

// MarkPacked marks shipment as packed
func (s *Shipment) MarkPacked() {
	s.Status = ShipmentStatusPacked
	s.UpdatedAt = time.Now()
}

// Ship marks shipment as shipped
func (s *Shipment) Ship(carrier, trackingNumber string) {
	now := time.Now()
	s.Status = ShipmentStatusShipped
	s.Carrier = carrier
	s.TrackingNumber = trackingNumber
	s.ShippedDate = &now
	s.UpdatedAt = now
}

// MarkInTransit marks shipment as in transit
func (s *Shipment) MarkInTransit() {
	s.Status = ShipmentStatusInTransit
	s.UpdatedAt = time.Now()
}

// MarkDelivered marks shipment as delivered
func (s *Shipment) MarkDelivered() {
	now := time.Now()
	s.Status = ShipmentStatusDelivered
	s.ActualDeliveryDate = &now
	s.UpdatedAt = now
}

// MarkReturned marks shipment as returned
func (s *Shipment) MarkReturned() {
	s.Status = ShipmentStatusReturned
	s.UpdatedAt = time.Now()
}

// CanBeShipped checks if shipment can be shipped
func (s *Shipment) CanBeShipped() bool {
	return s.Status == ShipmentStatusPending || s.Status == ShipmentStatusPicked || s.Status == ShipmentStatusPacked
}

// CanBeDelivered checks if shipment can be marked as delivered
func (s *Shipment) CanBeDelivered() bool {
	return s.Status == ShipmentStatusShipped || s.Status == ShipmentStatusInTransit
}
