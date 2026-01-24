package entity

import (
	"time"

	"github.com/google/uuid"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeIn         MovementType = "IN"
	MovementTypeOut        MovementType = "OUT"
	MovementTypeTransfer   MovementType = "TRANSFER"
	MovementTypeAdjustment MovementType = "ADJUSTMENT"
)

// ReferenceType represents the reference type for movement
type ReferenceType string

const (
	ReferenceTypeGRN         ReferenceType = "GRN"
	ReferenceTypeGI          ReferenceType = "GI"
	ReferenceTypeWO          ReferenceType = "WO"
	ReferenceTypeTransfer    ReferenceType = "TRANSFER"
	ReferenceTypeAdjustment  ReferenceType = "ADJUSTMENT"
	ReferenceTypeReservation ReferenceType = "RESERVATION"
)

// StockMovement represents a stock movement transaction
type StockMovement struct {
	ID             uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	MovementNumber string        `json:"movement_number" gorm:"type:varchar(30);not null"`
	MovementType   MovementType  `json:"movement_type" gorm:"type:varchar(20);not null"`
	ReferenceType  ReferenceType `json:"reference_type" gorm:"type:varchar(30)"`
	ReferenceID    *uuid.UUID    `json:"reference_id" gorm:"type:uuid"`
	MaterialID     uuid.UUID     `json:"material_id" gorm:"type:uuid;not null"`
	LotID          *uuid.UUID    `json:"lot_id" gorm:"type:uuid"`
	FromLocationID *uuid.UUID    `json:"from_location_id" gorm:"type:uuid"`
	ToLocationID   *uuid.UUID    `json:"to_location_id" gorm:"type:uuid"`
	Quantity       float64       `json:"quantity" gorm:"type:decimal(15,4);not null"`
	UnitID         uuid.UUID     `json:"unit_id" gorm:"type:uuid;not null"`
	Notes          string        `json:"notes" gorm:"type:text"`
	CreatedBy      uuid.UUID     `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt      time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Lot          *Lot      `json:"lot,omitempty" gorm:"foreignKey:LotID"`
	FromLocation *Location `json:"from_location,omitempty" gorm:"foreignKey:FromLocationID"`
	ToLocation   *Location `json:"to_location,omitempty" gorm:"foreignKey:ToLocationID"`
}

// TableName returns the table name
func (StockMovement) TableName() string {
	return "stock_movements"
}

// NewStockMovementIn creates an IN movement
func NewStockMovementIn(materialID, lotID, locationID, unitID, createdBy uuid.UUID, qty float64, refType ReferenceType, refID *uuid.UUID, movementNumber string) *StockMovement {
	return &StockMovement{
		MovementNumber: movementNumber,
		MovementType:   MovementTypeIn,
		ReferenceType:  refType,
		ReferenceID:    refID,
		MaterialID:     materialID,
		LotID:          &lotID,
		ToLocationID:   &locationID,
		Quantity:       qty,
		UnitID:         unitID,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}
}

// NewStockMovementOut creates an OUT movement
func NewStockMovementOut(materialID uuid.UUID, lotID, locationID *uuid.UUID, unitID, createdBy uuid.UUID, qty float64, refType ReferenceType, refID *uuid.UUID, movementNumber string) *StockMovement {
	return &StockMovement{
		MovementNumber: movementNumber,
		MovementType:   MovementTypeOut,
		ReferenceType:  refType,
		ReferenceID:    refID,
		MaterialID:     materialID,
		LotID:          lotID,
		FromLocationID: locationID,
		Quantity:       qty,
		UnitID:         unitID,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}
}

// NewStockMovementTransfer creates a TRANSFER movement
func NewStockMovementTransfer(materialID uuid.UUID, lotID *uuid.UUID, fromLocationID, toLocationID, unitID, createdBy uuid.UUID, qty float64, movementNumber string) *StockMovement {
	return &StockMovement{
		MovementNumber: movementNumber,
		MovementType:   MovementTypeTransfer,
		ReferenceType:  ReferenceTypeTransfer,
		MaterialID:     materialID,
		LotID:          lotID,
		FromLocationID: &fromLocationID,
		ToLocationID:   &toLocationID,
		Quantity:       qty,
		UnitID:         unitID,
		CreatedBy:      createdBy,
		CreatedAt:      time.Now(),
	}
}
