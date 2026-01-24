package entity

import (
	"time"

	"github.com/google/uuid"
)

// InventoryCountStatus represents the status of inventory count
type InventoryCountStatus string

const (
	InventoryCountStatusDraft      InventoryCountStatus = "DRAFT"
	InventoryCountStatusInProgress InventoryCountStatus = "IN_PROGRESS"
	InventoryCountStatusCompleted  InventoryCountStatus = "COMPLETED"
	InventoryCountStatusCancelled  InventoryCountStatus = "CANCELLED"
)

// InventoryCountType represents the type of count
type InventoryCountType string

const (
	InventoryCountTypeFull   InventoryCountType = "FULL"
	InventoryCountTypeCycle  InventoryCountType = "CYCLE"
	InventoryCountTypeSpot   InventoryCountType = "SPOT"
)

// InventoryCount represents a physical inventory count
type InventoryCount struct {
	ID            uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CountNumber   string               `json:"count_number" gorm:"type:varchar(30);uniqueIndex;not null"`
	CountDate     time.Time            `json:"count_date" gorm:"type:date;not null"`
	CountType     InventoryCountType   `json:"count_type" gorm:"type:varchar(20);not null"`
	WarehouseID   uuid.UUID            `json:"warehouse_id" gorm:"type:uuid;not null"`
	ZoneID        *uuid.UUID           `json:"zone_id" gorm:"type:uuid"`
	Status        InventoryCountStatus `json:"status" gorm:"type:varchar(20);default:'DRAFT'"`
	Notes         string               `json:"notes" gorm:"type:text"`
	StartedAt     *time.Time           `json:"started_at"`
	CompletedAt   *time.Time           `json:"completed_at"`
	CreatedBy     uuid.UUID            `json:"created_by" gorm:"type:uuid;not null"`
	ApprovedBy    *uuid.UUID           `json:"approved_by" gorm:"type:uuid"`
	CreatedAt     time.Time            `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time            `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Warehouse *Warehouse               `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Zone      *Zone                    `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
	LineItems []InventoryCountLineItem `json:"line_items,omitempty" gorm:"foreignKey:InventoryCountID"`
}

// TableName returns the table name
func (InventoryCount) TableName() string {
	return "inventory_counts"
}

// InventoryCountLineItem represents a line item in inventory count
type InventoryCountLineItem struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InventoryCountID uuid.UUID  `json:"inventory_count_id" gorm:"type:uuid;not null"`
	LocationID       uuid.UUID  `json:"location_id" gorm:"type:uuid;not null"`
	MaterialID       uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	LotID            *uuid.UUID `json:"lot_id" gorm:"type:uuid"`
	UnitID           uuid.UUID  `json:"unit_id" gorm:"type:uuid;not null"`
	SystemQty        float64    `json:"system_qty" gorm:"type:decimal(15,4);not null"`
	CountedQty       *float64   `json:"counted_qty" gorm:"type:decimal(15,4)"`
	Variance         float64    `json:"variance" gorm:"type:decimal(15,4);default:0"`
	VariancePercent  float64    `json:"variance_percent" gorm:"type:decimal(8,4);default:0"`
	IsCounted        bool       `json:"is_counted" gorm:"default:false"`
	CountedBy        *uuid.UUID `json:"counted_by" gorm:"type:uuid"`
	CountedAt        *time.Time `json:"counted_at"`
	Notes            string     `json:"notes" gorm:"type:text"`
	CreatedAt        time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Location *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`
	Lot      *Lot      `json:"lot,omitempty" gorm:"foreignKey:LotID"`
}

// TableName returns the table name
func (InventoryCountLineItem) TableName() string {
	return "inventory_count_lines"
}

// Start starts the inventory count
func (ic *InventoryCount) Start() {
	now := time.Now()
	ic.Status = InventoryCountStatusInProgress
	ic.StartedAt = &now
	ic.UpdatedAt = now
}

// Complete completes the inventory count
func (ic *InventoryCount) Complete(approvedBy uuid.UUID) {
	now := time.Now()
	ic.Status = InventoryCountStatusCompleted
	ic.CompletedAt = &now
	ic.ApprovedBy = &approvedBy
	ic.UpdatedAt = now
}

// Cancel cancels the inventory count
func (ic *InventoryCount) Cancel() {
	ic.Status = InventoryCountStatusCancelled
	ic.UpdatedAt = time.Now()
}

// CanStart checks if count can be started
func (ic *InventoryCount) CanStart() bool {
	return ic.Status == InventoryCountStatusDraft
}

// CanComplete checks if count can be completed
func (ic *InventoryCount) CanComplete() bool {
	return ic.Status == InventoryCountStatusInProgress
}

// RecordCount records the counted quantity
func (li *InventoryCountLineItem) RecordCount(countedQty float64, countedBy uuid.UUID) {
	now := time.Now()
	li.CountedQty = &countedQty
	li.Variance = countedQty - li.SystemQty
	if li.SystemQty > 0 {
		li.VariancePercent = (li.Variance / li.SystemQty) * 100
	}
	li.IsCounted = true
	li.CountedBy = &countedBy
	li.CountedAt = &now
}

// HasVariance checks if there's a variance
func (li *InventoryCountLineItem) HasVariance() bool {
	return li.Variance != 0
}
