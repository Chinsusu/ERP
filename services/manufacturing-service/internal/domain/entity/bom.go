package entity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

// BOMStatus represents BOM status
type BOMStatus string

const (
	BOMStatusDraft           BOMStatus = "DRAFT"
	BOMStatusPendingApproval BOMStatus = "PENDING_APPROVAL"
	BOMStatusApproved        BOMStatus = "APPROVED"
	BOMStatusObsolete        BOMStatus = "OBSOLETE"
)

// ConfidentialityLevel represents the access level for BOM
type ConfidentialityLevel string

const (
	ConfidentialityPublic      ConfidentialityLevel = "PUBLIC"
	ConfidentialityInternal    ConfidentialityLevel = "INTERNAL"
	ConfidentialityConfidental ConfidentialityLevel = "CONFIDENTIAL"
	ConfidentialityRestricted  ConfidentialityLevel = "RESTRICTED"
)

// BOMItemType represents the type of BOM item
type BOMItemType string

const (
	BOMItemTypeMaterial   BOMItemType = "MATERIAL"
	BOMItemTypePackaging  BOMItemType = "PACKAGING"
	BOMItemTypeConsumable BOMItemType = "CONSUMABLE"
)

// BOM represents a Bill of Materials with encrypted formula
type BOM struct {
	ID                   uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BOMNumber            string               `json:"bom_number" gorm:"type:varchar(50);unique;not null"`
	ProductID            uuid.UUID            `json:"product_id" gorm:"type:uuid;not null"`
	Version              int                  `json:"version" gorm:"default:1"`
	Name                 string               `json:"name" gorm:"type:varchar(200);not null"`
	Description          string               `json:"description" gorm:"type:text"`
	Status               BOMStatus            `json:"status" gorm:"type:varchar(30);default:'DRAFT'"`
	BatchSize            float64              `json:"batch_size" gorm:"type:decimal(15,4);not null"`
	BatchUnitID          uuid.UUID            `json:"batch_unit_id" gorm:"type:uuid;not null"`
	FormulaDetails       []byte               `json:"-" gorm:"type:bytea"` // Encrypted - never expose directly
	ConfidentialityLevel ConfidentialityLevel `json:"confidentiality_level" gorm:"type:varchar(30);default:'RESTRICTED'"`
	MaterialCost         float64              `json:"material_cost" gorm:"type:decimal(18,2);default:0"`
	LaborCost            float64              `json:"labor_cost" gorm:"type:decimal(18,2);default:0"`
	OverheadCost         float64              `json:"overhead_cost" gorm:"type:decimal(18,2);default:0"`
	TotalCost            float64              `json:"total_cost" gorm:"type:decimal(18,2);default:0"`
	EffectiveFrom        *time.Time           `json:"effective_from" gorm:"type:date"`
	EffectiveTo          *time.Time           `json:"effective_to" gorm:"type:date"`
	ApprovedBy           *uuid.UUID           `json:"approved_by" gorm:"type:uuid"`
	ApprovedAt           *time.Time           `json:"approved_at"`
	CreatedBy            *uuid.UUID           `json:"created_by" gorm:"type:uuid"`
	UpdatedBy            *uuid.UUID           `json:"updated_by" gorm:"type:uuid"`
	CreatedAt            time.Time            `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time            `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Associations
	Items []BOMLineItem `json:"items,omitempty" gorm:"foreignKey:BOMID"`
}

// TableName returns the table name
func (BOM) TableName() string {
	return "boms"
}

// BOMLineItem represents a line item in a BOM
type BOMLineItem struct {
	ID              uuid.UUID   `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BOMID           uuid.UUID   `json:"bom_id" gorm:"type:uuid;not null"`
	LineNumber      int         `json:"line_number" gorm:"not null"`
	MaterialID      uuid.UUID   `json:"material_id" gorm:"type:uuid;not null"`
	ItemType        BOMItemType `json:"item_type" gorm:"type:varchar(30);default:'MATERIAL'"`
	Quantity        float64     `json:"quantity" gorm:"type:decimal(15,4);not null"`
	UOMID           uuid.UUID   `json:"uom_id" gorm:"type:uuid;not null"`
	QuantityMin     *float64    `json:"quantity_min" gorm:"type:decimal(15,4)"`
	QuantityMax     *float64    `json:"quantity_max" gorm:"type:decimal(15,4)"`
	IsCritical      bool        `json:"is_critical" gorm:"default:false"`
	ScrapPercentage float64     `json:"scrap_percentage" gorm:"type:decimal(5,2);default:0"`
	UnitCost        float64     `json:"unit_cost" gorm:"type:decimal(18,4);default:0"`
	TotalCost       float64     `json:"total_cost" gorm:"type:decimal(18,2);default:0"`
	Notes           string      `json:"notes" gorm:"type:text"`
	CreatedAt       time.Time   `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time   `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (BOMLineItem) TableName() string {
	return "bom_line_items"
}

// BOMVersion represents a version snapshot of a BOM
type BOMVersion struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	BOMID        uuid.UUID  `json:"bom_id" gorm:"type:uuid;not null"`
	Version      int        `json:"version" gorm:"not null"`
	ChangeReason string     `json:"change_reason" gorm:"type:text"`
	ChangedBy    *uuid.UUID `json:"changed_by" gorm:"type:uuid"`
	ChangedAt    time.Time  `json:"changed_at" gorm:"default:CURRENT_TIMESTAMP"`
	Snapshot     []byte     `json:"snapshot" gorm:"type:jsonb;not null"` // JSON snapshot of BOM
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name
func (BOMVersion) TableName() string {
	return "bom_versions"
}

// FormulaDetails represents the decrypted formula content
type FormulaDetails struct {
	ProcessingSteps    []string          `json:"processing_steps"`
	CriticalParameters map[string]string `json:"critical_parameters"`
	Notes              string            `json:"notes"`
}

// EncryptFormula encrypts formula details using AES-256-GCM
func EncryptFormula(formula *FormulaDetails, key []byte) ([]byte, error) {
	if formula == nil {
		return nil, nil
	}

	plaintext, err := json.Marshal(formula)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptFormula decrypts formula details using AES-256-GCM
func DecryptFormula(ciphertext []byte, key []byte) (*FormulaDetails, error) {
	if len(ciphertext) == 0 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	var formula FormulaDetails
	if err := json.Unmarshal(plaintext, &formula); err != nil {
		return nil, err
	}

	return &formula, nil
}

// BOM business methods

// IsDraft returns true if BOM is in draft status
func (b *BOM) IsDraft() bool {
	return b.Status == BOMStatusDraft
}

// CanBeApproved returns true if BOM can be approved
func (b *BOM) CanBeApproved() bool {
	return b.Status == BOMStatusPendingApproval
}

// CanBeModified returns true if BOM can be modified
func (b *BOM) CanBeModified() bool {
	return b.Status == BOMStatusDraft
}

// Submit submits BOM for approval
func (b *BOM) Submit() error {
	if b.Status != BOMStatusDraft {
		return errors.New("only draft BOMs can be submitted")
	}
	b.Status = BOMStatusPendingApproval
	b.UpdatedAt = time.Now()
	return nil
}

// Approve approves the BOM
func (b *BOM) Approve(approverID uuid.UUID) error {
	if b.Status != BOMStatusPendingApproval {
		return errors.New("only pending approval BOMs can be approved")
	}
	b.Status = BOMStatusApproved
	b.ApprovedBy = &approverID
	now := time.Now()
	b.ApprovedAt = &now
	b.EffectiveFrom = &now
	b.UpdatedAt = now
	return nil
}

// MarkObsolete marks the BOM as obsolete
func (b *BOM) MarkObsolete() {
	b.Status = BOMStatusObsolete
	now := time.Now()
	b.EffectiveTo = &now
	b.UpdatedAt = now
}

// CalculateTotalCost calculates the total cost of the BOM
func (b *BOM) CalculateTotalCost() {
	var materialCost float64
	for _, item := range b.Items {
		materialCost += item.TotalCost
	}
	b.MaterialCost = materialCost
	b.TotalCost = b.MaterialCost + b.LaborCost + b.OverheadCost
}

// IsActive returns true if BOM is active (approved and within effective dates)
func (b *BOM) IsActive() bool {
	if b.Status != BOMStatusApproved {
		return false
	}
	now := time.Now()
	if b.EffectiveFrom != nil && now.Before(*b.EffectiveFrom) {
		return false
	}
	if b.EffectiveTo != nil && now.After(*b.EffectiveTo) {
		return false
	}
	return true
}
