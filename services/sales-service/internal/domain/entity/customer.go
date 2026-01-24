package entity

import (
	"time"

	"github.com/google/uuid"
)

// CustomerType represents customer classification
type CustomerType string

const (
	CustomerTypeRetail      CustomerType = "RETAIL"
	CustomerTypeWholesale   CustomerType = "WHOLESALE"
	CustomerTypeDistributor CustomerType = "DISTRIBUTOR"
)

// CustomerStatus represents customer status
type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "ACTIVE"
	CustomerStatusInactive CustomerStatus = "INACTIVE"
	CustomerStatusBlocked  CustomerStatus = "BLOCKED"
)

// CustomerGroup represents pricing tier for customers
type CustomerGroup struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code            string    `json:"code" gorm:"type:varchar(20);unique;not null"`
	Name            string    `json:"name" gorm:"type:varchar(100);not null"`
	Description     string    `json:"description" gorm:"type:text"`
	DiscountPercent float64   `json:"discount_percent" gorm:"type:decimal(5,2);default:0"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (CustomerGroup) TableName() string {
	return "customer_groups"
}

// Customer represents a customer entity
type Customer struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerCode    string          `json:"customer_code" gorm:"type:varchar(20);unique;not null"`
	Name            string          `json:"name" gorm:"type:varchar(200);not null"`
	TaxCode         string          `json:"tax_code" gorm:"type:varchar(50)"`
	CustomerType    CustomerType    `json:"customer_type" gorm:"type:varchar(20);default:'RETAIL'"`
	CustomerGroupID *uuid.UUID      `json:"customer_group_id" gorm:"type:uuid"`
	CustomerGroup   *CustomerGroup  `json:"customer_group,omitempty" gorm:"foreignKey:CustomerGroupID"`
	Email           string          `json:"email" gorm:"type:varchar(100)"`
	Phone           string          `json:"phone" gorm:"type:varchar(20)"`
	Website         string          `json:"website" gorm:"type:varchar(200)"`
	PaymentTerms    string          `json:"payment_terms" gorm:"type:varchar(50);default:'Net 30'"`
	CreditLimit     float64         `json:"credit_limit" gorm:"type:decimal(18,2);default:0"`
	CurrentBalance  float64         `json:"current_balance" gorm:"type:decimal(18,2);default:0"`
	Currency        string          `json:"currency" gorm:"type:varchar(3);default:'VND'"`
	Status          CustomerStatus  `json:"status" gorm:"type:varchar(20);default:'ACTIVE'"`
	Notes           string          `json:"notes" gorm:"type:text"`
	CreatedBy       *uuid.UUID      `json:"created_by" gorm:"type:uuid"`
	UpdatedBy       *uuid.UUID      `json:"updated_by" gorm:"type:uuid"`
	CreatedAt       time.Time       `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Addresses []CustomerAddress `json:"addresses,omitempty" gorm:"foreignKey:CustomerID"`
	Contacts  []CustomerContact `json:"contacts,omitempty" gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string {
	return "customers"
}

// GetAvailableCredit returns the available credit for the customer
func (c *Customer) GetAvailableCredit() float64 {
	return c.CreditLimit - c.CurrentBalance
}

// CanPlaceOrder checks if customer can place an order of given amount
func (c *Customer) CanPlaceOrder(orderAmount float64) bool {
	if c.Status != CustomerStatusActive {
		return false
	}
	return c.GetAvailableCredit() >= orderAmount
}

// Block blocks the customer
func (c *Customer) Block() {
	c.Status = CustomerStatusBlocked
	c.UpdatedAt = time.Now()
}

// Activate activates the customer
func (c *Customer) Activate() {
	c.Status = CustomerStatusActive
	c.UpdatedAt = time.Now()
}

// AddressType represents the type of address
type AddressType string

const (
	AddressTypeBilling  AddressType = "BILLING"
	AddressTypeShipping AddressType = "SHIPPING"
	AddressTypeBoth     AddressType = "BOTH"
)

// CustomerAddress represents a customer address
type CustomerAddress struct {
	ID           uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID   uuid.UUID   `json:"customer_id" gorm:"type:uuid;not null"`
	AddressType  AddressType `json:"address_type" gorm:"type:varchar(20);not null"`
	AddressLine1 string      `json:"address_line1" gorm:"type:varchar(200);not null"`
	AddressLine2 string      `json:"address_line2" gorm:"type:varchar(200)"`
	City         string      `json:"city" gorm:"type:varchar(100)"`
	State        string      `json:"state" gorm:"type:varchar(100)"`
	PostalCode   string      `json:"postal_code" gorm:"type:varchar(20)"`
	Country      string      `json:"country" gorm:"type:varchar(100);default:'Vietnam'"`
	IsDefault    bool        `json:"is_default" gorm:"default:false"`
	CreatedAt    time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (CustomerAddress) TableName() string {
	return "customer_addresses"
}

// GetFullAddress returns the full formatted address
func (a *CustomerAddress) GetFullAddress() string {
	addr := a.AddressLine1
	if a.AddressLine2 != "" {
		addr += ", " + a.AddressLine2
	}
	if a.City != "" {
		addr += ", " + a.City
	}
	if a.Country != "" {
		addr += ", " + a.Country
	}
	return addr
}

// CustomerContact represents a customer contact person
type CustomerContact struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID  uuid.UUID `json:"customer_id" gorm:"type:uuid;not null"`
	ContactName string    `json:"contact_name" gorm:"type:varchar(100);not null"`
	Position    string    `json:"position" gorm:"type:varchar(100)"`
	Email       string    `json:"email" gorm:"type:varchar(100)"`
	Phone       string    `json:"phone" gorm:"type:varchar(20)"`
	Mobile      string    `json:"mobile" gorm:"type:varchar(20)"`
	IsPrimary   bool      `json:"is_primary" gorm:"default:false"`
	Notes       string    `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (CustomerContact) TableName() string {
	return "customer_contacts"
}
