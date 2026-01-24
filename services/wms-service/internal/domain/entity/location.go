package entity

import (
	"time"

	"github.com/google/uuid"
)

// Location represents a storage location within a zone
type Location struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ZoneID    uuid.UUID `json:"zone_id" gorm:"type:uuid;not null"`
	Code      string    `json:"code" gorm:"type:varchar(30);not null"` // A01-R02-S03-B01
	Aisle     string    `json:"aisle" gorm:"type:varchar(10)"`
	Rack      string    `json:"rack" gorm:"type:varchar(10)"`
	Shelf     string    `json:"shelf" gorm:"type:varchar(10)"`
	Bin       string    `json:"bin" gorm:"type:varchar(10)"`
	Capacity  *float64  `json:"capacity" gorm:"type:decimal(10,2)"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Zone *Zone `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
}

// TableName returns the table name
func (Location) TableName() string {
	return "locations"
}

// GetFullPath returns the full location path including warehouse and zone
func (l *Location) GetFullPath() string {
	if l.Zone == nil || l.Zone.Warehouse == nil {
		return l.Code
	}
	return l.Zone.Warehouse.Code + "/" + l.Zone.Code + "/" + l.Code
}
