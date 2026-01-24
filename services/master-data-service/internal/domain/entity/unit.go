package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UoMType represents the type of unit of measure
type UoMType string

const (
	UoMTypeWeight   UoMType = "WEIGHT"
	UoMTypeVolume   UoMType = "VOLUME"
	UoMTypeQuantity UoMType = "QUANTITY"
	UoMTypeLength   UoMType = "LENGTH"
	UoMTypeArea     UoMType = "AREA"
)

// UnitOfMeasure represents a unit of measurement
type UnitOfMeasure struct {
	ID               uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code             string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"code"`
	Name             string         `gorm:"type:varchar(100);not null" json:"name"`
	NameEN           string         `gorm:"type:varchar(100);column:name_en" json:"name_en,omitempty"`
	Symbol           string         `gorm:"type:varchar(10);not null" json:"symbol"`
	UoMType          UoMType        `gorm:"type:varchar(50);not null;column:uom_type" json:"uom_type"`
	IsBaseUnit       bool           `gorm:"not null;default:false" json:"is_base_unit"`
	BaseUnitID       *uuid.UUID     `gorm:"type:uuid" json:"base_unit_id,omitempty"`
	ConversionFactor float64        `gorm:"type:decimal(18,8);not null;default:1" json:"conversion_factor"`
	Description      string         `gorm:"type:text" json:"description,omitempty"`
	Status           string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	BaseUnit *UnitOfMeasure `gorm:"foreignKey:BaseUnitID" json:"base_unit,omitempty"`
}

// TableName specifies the table name
func (UnitOfMeasure) TableName() string {
	return "units_of_measure"
}

// Validate validates unit data
func (u *UnitOfMeasure) Validate() error {
	if u.Code == "" {
		return fmt.Errorf("code is required")
	}
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if u.UoMType == "" {
		return fmt.Errorf("uom_type is required")
	}
	return nil
}

// IsActive checks if unit is active
func (u *UnitOfMeasure) IsActive() bool {
	return u.Status == "active"
}

// UnitConversion represents a conversion factor between two units
type UnitConversion struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FromUnitID       uuid.UUID `gorm:"type:uuid;not null" json:"from_unit_id"`
	ToUnitID         uuid.UUID `gorm:"type:uuid;not null" json:"to_unit_id"`
	ConversionFactor float64   `gorm:"type:decimal(18,8);not null" json:"conversion_factor"`
	IsActive         bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	FromUnit *UnitOfMeasure `gorm:"foreignKey:FromUnitID" json:"from_unit,omitempty"`
	ToUnit   *UnitOfMeasure `gorm:"foreignKey:ToUnitID" json:"to_unit,omitempty"`
}

// TableName specifies the table name
func (UnitConversion) TableName() string {
	return "unit_conversions"
}

// Convert converts a value from source unit to target unit
func (uc *UnitConversion) Convert(value float64) float64 {
	return value * uc.ConversionFactor
}
