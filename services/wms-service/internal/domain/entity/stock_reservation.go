package entity

import (
	"time"

	"github.com/google/uuid"
)

// ReservationType represents the type of reservation
type ReservationType string

const (
	ReservationTypeSalesOrder ReservationType = "SALES_ORDER"
	ReservationTypeWorkOrder  ReservationType = "WORK_ORDER"
	ReservationTypeTransfer   ReservationType = "TRANSFER"
)

// ReservationStatus represents the status of reservation
type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "ACTIVE"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

// StockReservation represents a stock reservation
type StockReservation struct {
	ID              uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MaterialID      uuid.UUID         `json:"material_id" gorm:"type:uuid;not null"`
	LotID           *uuid.UUID        `json:"lot_id" gorm:"type:uuid"`
	LocationID      *uuid.UUID        `json:"location_id" gorm:"type:uuid"`
	Quantity        float64           `json:"quantity" gorm:"type:decimal(15,4);not null"`
	UnitID          uuid.UUID         `json:"unit_id" gorm:"type:uuid;not null"`
	ReservationType ReservationType   `json:"reservation_type" gorm:"type:varchar(30);not null"`
	ReferenceID     uuid.UUID         `json:"reference_id" gorm:"type:uuid;not null"`
	ReferenceNumber string            `json:"reference_number" gorm:"type:varchar(30)"`
	Status          ReservationStatus `json:"status" gorm:"type:varchar(20);default:'ACTIVE'"`
	ExpiresAt       *time.Time        `json:"expires_at"`
	CreatedBy       uuid.UUID         `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt       time.Time         `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	ReleasedAt      *time.Time        `json:"released_at"`

	// Relations
	Lot      *Lot      `json:"lot,omitempty" gorm:"foreignKey:LotID"`
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
}

// TableName returns the table name
func (StockReservation) TableName() string {
	return "stock_reservations"
}

// IsActive returns true if reservation is active
func (r *StockReservation) IsActive() bool {
	return r.Status == ReservationStatusActive
}

// IsExpired returns true if reservation is expired
func (r *StockReservation) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// Release releases the reservation
func (r *StockReservation) Release() {
	now := time.Now()
	r.Status = ReservationStatusReleased
	r.ReleasedAt = &now
}

// Fulfill marks reservation as fulfilled
func (r *StockReservation) Fulfill() {
	now := time.Now()
	r.Status = ReservationStatusFulfilled
	r.ReleasedAt = &now
}

// MarkExpired marks reservation as expired
func (r *StockReservation) MarkExpired() {
	r.Status = ReservationStatusExpired
}
