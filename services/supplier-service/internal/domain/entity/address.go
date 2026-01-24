package entity

import (
	"time"

	"github.com/google/uuid"
)

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeBilling   AddressType = "BILLING"
	AddressTypeShipping  AddressType = "SHIPPING"
	AddressTypeFactory   AddressType = "FACTORY"
	AddressTypeOffice    AddressType = "OFFICE"
	AddressTypeWarehouse AddressType = "WAREHOUSE"
)

// Address represents a supplier address
type Address struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SupplierID   uuid.UUID   `json:"supplier_id" gorm:"type:uuid;not null"`
	AddressType  AddressType `json:"address_type" gorm:"type:varchar(20);not null"`
	AddressLine1 string      `json:"address_line1" gorm:"type:varchar(255);not null"`
	AddressLine2 string      `json:"address_line2" gorm:"type:varchar(255)"`
	Ward         string      `json:"ward" gorm:"type:varchar(100)"`
	District     string      `json:"district" gorm:"type:varchar(100)"`
	City         string      `json:"city" gorm:"type:varchar(100);not null"`
	Province     string      `json:"province" gorm:"type:varchar(100)"`
	Country      string      `json:"country" gorm:"type:varchar(100);not null;default:'Vietnam'"`
	PostalCode   string      `json:"postal_code" gorm:"type:varchar(20)"`
	IsPrimary    bool        `json:"is_primary" gorm:"default:false"`
	Notes        string      `json:"notes" gorm:"type:text"`
	CreatedAt    time.Time   `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for GORM
func (Address) TableName() string {
	return "supplier_addresses"
}

// FullAddress returns formatted full address
func (a *Address) FullAddress() string {
	addr := a.AddressLine1
	if a.AddressLine2 != "" {
		addr += ", " + a.AddressLine2
	}
	if a.Ward != "" {
		addr += ", " + a.Ward
	}
	if a.District != "" {
		addr += ", " + a.District
	}
	addr += ", " + a.City
	if a.Province != "" && a.Province != a.City {
		addr += ", " + a.Province
	}
	addr += ", " + a.Country
	if a.PostalCode != "" {
		addr += " " + a.PostalCode
	}
	return addr
}
