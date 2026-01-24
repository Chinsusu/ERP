package entity

import (
	"time"

	"github.com/google/uuid"
)

// QCStatus represents QC status
type QCStatus string

const (
	QCStatusPending    QCStatus = "PENDING"
	QCStatusPassed     QCStatus = "PASSED"
	QCStatusFailed     QCStatus = "FAILED"
	QCStatusQuarantine QCStatus = "QUARANTINE"
)

// LotStatus represents lot status
type LotStatus string

const (
	LotStatusAvailable LotStatus = "AVAILABLE"
	LotStatusReserved  LotStatus = "RESERVED"
	LotStatusBlocked   LotStatus = "BLOCKED"
	LotStatusExpired   LotStatus = "EXPIRED"
)

// Lot represents a lot/batch of material - CRITICAL for FEFO
type Lot struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LotNumber         string     `json:"lot_number" gorm:"type:varchar(30);unique;not null"` // LOT-YYYYMM-XXXX
	MaterialID        uuid.UUID  `json:"material_id" gorm:"type:uuid;not null"`
	SupplierID        *uuid.UUID `json:"supplier_id" gorm:"type:uuid"`
	SupplierLotNumber string     `json:"supplier_lot_number" gorm:"type:varchar(50)"`
	ManufacturedDate  *time.Time `json:"manufactured_date" gorm:"type:date"`
	ExpiryDate        time.Time  `json:"expiry_date" gorm:"type:date;not null"` // Critical for FEFO
	ReceivedDate      time.Time  `json:"received_date" gorm:"type:date;not null"`
	GRNID             *uuid.UUID `json:"grn_id" gorm:"type:uuid"`
	QCStatus          QCStatus   `json:"qc_status" gorm:"type:varchar(20);default:'PENDING'"`
	Status            LotStatus  `json:"status" gorm:"type:varchar(20);default:'AVAILABLE'"`
	Notes             string     `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (Lot) TableName() string {
	return "lots"
}

// DaysUntilExpiry returns the number of days until expiry
func (l *Lot) DaysUntilExpiry() int {
	duration := time.Until(l.ExpiryDate)
	return int(duration.Hours() / 24)
}

// IsExpired returns true if the lot is expired
func (l *Lot) IsExpired() bool {
	return time.Now().After(l.ExpiryDate)
}

// IsExpiringSoon returns true if the lot expires within the given days
func (l *Lot) IsExpiringSoon(days int) bool {
	threshold := time.Now().AddDate(0, 0, days)
	return l.ExpiryDate.Before(threshold) && !l.IsExpired()
}

// MarkExpired marks the lot as expired
func (l *Lot) MarkExpired() {
	l.Status = LotStatusExpired
	l.UpdatedAt = time.Now()
}

// PassQC marks the lot as QC passed
func (l *Lot) PassQC() {
	l.QCStatus = QCStatusPassed
	l.Status = LotStatusAvailable
	l.UpdatedAt = time.Now()
}

// FailQC marks the lot as QC failed
func (l *Lot) FailQC() {
	l.QCStatus = QCStatusFailed
	l.Status = LotStatusBlocked
	l.UpdatedAt = time.Now()
}

// CanBeIssued returns true if the lot can be issued
func (l *Lot) CanBeIssued() bool {
	return l.Status == LotStatusAvailable && 
		l.QCStatus == QCStatusPassed && 
		!l.IsExpired()
}

// Block blocks the lot
func (l *Lot) Block() {
	l.Status = LotStatusBlocked
	l.UpdatedAt = time.Now()
}

// Unblock unblocks the lot
func (l *Lot) Unblock() {
	if l.QCStatus == QCStatusPassed && !l.IsExpired() {
		l.Status = LotStatusAvailable
	}
	l.UpdatedAt = time.Now()
}
