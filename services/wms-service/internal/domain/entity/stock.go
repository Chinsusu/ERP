package entity

import (
	"time"

	"github.com/google/uuid"
)

// Stock represents stock at a specific location with lot
type Stock struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WarehouseID  uuid.UUID  `json:"warehouse_id" gorm:"type:uuid;not null"`
	ZoneID       uuid.UUID  `json:"zone_id" gorm:"type:uuid;not null"`
	LocationID   uuid.UUID  `json:"location_id" gorm:"type:uuid;not null"`
	MaterialID   uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	LotID        *uuid.UUID `json:"lot_id" gorm:"type:uuid"`
	Quantity     float64    `json:"quantity" gorm:"type:decimal(15,4);not null;default:0"`
	ReservedQty  float64    `json:"reserved_qty" gorm:"type:decimal(15,4);not null;default:0"`
	AvailableQty float64    `json:"available_qty" gorm:"type:decimal(15,4)"` // Computed: Quantity - ReservedQty
	UnitID       uuid.UUID  `json:"unit_id" gorm:"type:uuid;not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Warehouse *Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Zone      *Zone      `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
	Location  *Location  `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	Lot       *Lot       `json:"lot,omitempty" gorm:"foreignKey:LotID"`
}

// TableName returns the table name
func (Stock) TableName() string {
	return "stock"
}

// GetAvailableQuantity returns available quantity (quantity - reserved)
func (s *Stock) GetAvailableQuantity() float64 {
	return s.Quantity - s.ReservedQty
}

// CanIssue returns true if quantity can be issued
func (s *Stock) CanIssue(qty float64) bool {
	return s.GetAvailableQuantity() >= qty
}

// Reserve reserves stock
func (s *Stock) Reserve(qty float64) error {
	if !s.CanIssue(qty) {
		return ErrInsufficientStock
	}
	s.ReservedQty += qty
	s.AvailableQty = s.GetAvailableQuantity()
	s.UpdatedAt = time.Now()
	return nil
}

// ReleaseReservation releases reserved stock
func (s *Stock) ReleaseReservation(qty float64) {
	s.ReservedQty -= qty
	if s.ReservedQty < 0 {
		s.ReservedQty = 0
	}
	s.AvailableQty = s.GetAvailableQuantity()
	s.UpdatedAt = time.Now()
}

// Issue issues stock (deduct quantity)
func (s *Stock) Issue(qty float64) error {
	if s.Quantity < qty {
		return ErrInsufficientStock
	}
	s.Quantity -= qty
	s.AvailableQty = s.GetAvailableQuantity()
	s.UpdatedAt = time.Now()
	return nil
}

// Receive adds stock
func (s *Stock) Receive(qty float64) {
	s.Quantity += qty
	s.AvailableQty = s.GetAvailableQuantity()
	s.UpdatedAt = time.Now()
}

// LotIssued represents a lot that was issued (for FEFO tracking)
type LotIssued struct {
	LotID      uuid.UUID `json:"lot_id"`
	LotNumber  string    `json:"lot_number"`
	Quantity   float64   `json:"quantity"`
	ExpiryDate time.Time `json:"expiry_date"`
	LocationID uuid.UUID `json:"location_id"`
}

// StockSummary represents aggregated stock summary
type StockSummary struct {
	MaterialID    uuid.UUID `json:"material_id"`
	MaterialCode  string    `json:"material_code"`
	MaterialName  string    `json:"material_name"`
	TotalQuantity float64   `json:"total_quantity"`
	TotalReserved float64   `json:"total_reserved"`
	TotalAvailable float64  `json:"total_available"`
	UnitCode      string    `json:"unit_code"`
}
