package entity

import (
	"time"

	"github.com/google/uuid"
)

// WarehouseType represents the type of warehouse
type WarehouseType string

const (
	WarehouseTypeMain          WarehouseType = "MAIN"
	WarehouseTypeColdStorage   WarehouseType = "COLD_STORAGE"
	WarehouseTypeFinishedGoods WarehouseType = "FINISHED_GOODS"
	WarehouseTypeQuarantine    WarehouseType = "QUARANTINE"
)

// Warehouse represents a warehouse entity
type Warehouse struct {
	ID            uuid.UUID     `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code          string        `json:"code" gorm:"type:varchar(20);unique;not null"`
	Name          string        `json:"name" gorm:"type:varchar(100);not null"`
	WarehouseType WarehouseType `json:"warehouse_type" gorm:"type:varchar(20);not null"`
	Address       string        `json:"address" gorm:"type:text"`
	IsActive      bool          `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time     `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Zones []Zone `json:"zones,omitempty" gorm:"foreignKey:WarehouseID"`
}

// TableName returns the table name
func (Warehouse) TableName() string {
	return "warehouses"
}

// ZoneType represents the type of zone
type ZoneType string

const (
	ZoneTypeReceiving  ZoneType = "RECEIVING"
	ZoneTypeQuarantine ZoneType = "QUARANTINE"
	ZoneTypeStorage    ZoneType = "STORAGE"
	ZoneTypeCold       ZoneType = "COLD"
	ZoneTypePicking    ZoneType = "PICKING"
	ZoneTypeShipping   ZoneType = "SHIPPING"
)

// Zone represents a zone within a warehouse
type Zone struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	WarehouseID    uuid.UUID  `json:"warehouse_id" gorm:"type:uuid;not null"`
	Code           string     `json:"code" gorm:"type:varchar(20);not null"`
	Name           string     `json:"name" gorm:"type:varchar(100);not null"`
	ZoneType       ZoneType   `json:"zone_type" gorm:"type:varchar(30);not null"`
	TemperatureMin *float64   `json:"temperature_min" gorm:"type:decimal(5,2)"`
	TemperatureMax *float64   `json:"temperature_max" gorm:"type:decimal(5,2)"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Warehouse *Warehouse `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Locations []Location `json:"locations,omitempty" gorm:"foreignKey:ZoneID"`
}

// TableName returns the table name
func (Zone) TableName() string {
	return "zones"
}

// IsColdZone returns true if zone requires temperature control
func (z *Zone) IsColdZone() bool {
	return z.ZoneType == ZoneTypeCold
}

// IsQuarantineZone returns true if zone is quarantine
func (z *Zone) IsQuarantineZone() bool {
	return z.ZoneType == ZoneTypeQuarantine
}
