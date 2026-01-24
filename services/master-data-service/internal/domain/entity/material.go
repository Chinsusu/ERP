package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MaterialType represents the type of material
type MaterialType string

const (
	MaterialTypeRaw         MaterialType = "RAW_MATERIAL"
	MaterialTypePackaging   MaterialType = "PACKAGING"
	MaterialTypeConsumable  MaterialType = "CONSUMABLE"
	MaterialTypeSemiFinished MaterialType = "SEMI_FINISHED"
)

// StorageCondition represents storage requirements
type StorageCondition string

const (
	StorageAmbient StorageCondition = "AMBIENT"
	StorageCold    StorageCondition = "COLD"   // 2-8°C
	StorageFrozen  StorageCondition = "FROZEN" // <-18°C
)

// Material represents a raw material, packaging, or consumable
type Material struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	NameEN      string         `gorm:"type:varchar(255);column:name_en" json:"name_en,omitempty"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	
	// Classification
	MaterialType MaterialType   `gorm:"type:varchar(50);not null" json:"material_type"`
	CategoryID   *uuid.UUID     `gorm:"type:uuid" json:"category_id,omitempty"`
	
	// Units
	BaseUnitID     uuid.UUID  `gorm:"type:uuid;not null" json:"base_unit_id"`
	PurchaseUnitID *uuid.UUID `gorm:"type:uuid" json:"purchase_unit_id,omitempty"`
	StockUnitID    *uuid.UUID `gorm:"type:uuid" json:"stock_unit_id,omitempty"`

	// Cosmetics specific
	INCIName    string `gorm:"type:varchar(500);column:inci_name" json:"inci_name,omitempty"`
	CASNumber   string `gorm:"type:varchar(50);column:cas_number" json:"cas_number,omitempty"`
	IsAllergen  bool   `gorm:"not null;default:false" json:"is_allergen"`
	AllergenInfo string `gorm:"type:text" json:"allergen_info,omitempty"`
	IsOrganic   bool   `gorm:"not null;default:false" json:"is_organic"`
	IsNatural   bool   `gorm:"not null;default:false" json:"is_natural"`
	IsVegan     bool   `gorm:"not null;default:false" json:"is_vegan"`
	OriginCountry string `gorm:"type:varchar(100)" json:"origin_country,omitempty"`

	// Storage
	StorageCondition    StorageCondition `gorm:"type:varchar(50);not null;default:'AMBIENT'" json:"storage_condition"`
	MinTemp             *float64         `gorm:"type:decimal(5,2)" json:"min_temp,omitempty"`
	MaxTemp             *float64         `gorm:"type:decimal(5,2)" json:"max_temp,omitempty"`
	StorageInstructions string           `gorm:"type:text" json:"storage_instructions,omitempty"`
	ShelfLifeDays       int              `gorm:"not null;default:365" json:"shelf_life_days"`

	// Safety
	IsHazardous       bool   `gorm:"not null;default:false" json:"is_hazardous"`
	HazardClass       string `gorm:"type:varchar(50)" json:"hazard_class,omitempty"`
	SafetyDataSheetURL string `gorm:"type:varchar(500);column:safety_data_sheet_url" json:"safety_data_sheet_url,omitempty"`

	// Procurement
	DefaultSupplierID *uuid.UUID `gorm:"type:uuid" json:"default_supplier_id,omitempty"`
	LeadTimeDays      int        `gorm:"not null;default:14" json:"lead_time_days"`
	MinOrderQty       float64    `gorm:"type:decimal(18,4);not null;default:1" json:"min_order_qty"`
	ReorderPoint      float64    `gorm:"type:decimal(18,4);not null;default:0" json:"reorder_point"`
	SafetyStock       float64    `gorm:"type:decimal(18,4);not null;default:0" json:"safety_stock"`
	MaxStockQty       *float64   `gorm:"type:decimal(18,4)" json:"max_stock_qty,omitempty"`

	// Costing
	StandardCost     float64 `gorm:"type:decimal(18,4);not null;default:0" json:"standard_cost"`
	LastPurchaseCost *float64 `gorm:"type:decimal(18,4)" json:"last_purchase_cost,omitempty"`
	Currency         string  `gorm:"type:varchar(3);not null;default:'VND'" json:"currency"`

	// Status
	Status    string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Category     *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	BaseUnit     *UnitOfMeasure `gorm:"foreignKey:BaseUnitID" json:"base_unit,omitempty"`
	PurchaseUnit *UnitOfMeasure `gorm:"foreignKey:PurchaseUnitID" json:"purchase_unit,omitempty"`
	StockUnit    *UnitOfMeasure `gorm:"foreignKey:StockUnitID" json:"stock_unit,omitempty"`
	Specifications []MaterialSpecification `gorm:"foreignKey:MaterialID" json:"specifications,omitempty"`
}

// TableName specifies the table name
func (Material) TableName() string {
	return "materials"
}

// Validate validates material data
func (m *Material) Validate() error {
	if m.Code == "" {
		return fmt.Errorf("code is required")
	}
	if m.Name == "" {
		return fmt.Errorf("name is required")
	}
	if m.MaterialType == "" {
		return fmt.Errorf("material_type is required")
	}
	if m.BaseUnitID == uuid.Nil {
		return fmt.Errorf("base_unit_id is required")
	}
	return nil
}

// IsActive checks if material is active
func (m *Material) IsActive() bool {
	return m.Status == "active"
}

// RequiresColdStorage checks if material needs cold storage
func (m *Material) RequiresColdStorage() bool {
	return m.StorageCondition == StorageCold || m.StorageCondition == StorageFrozen
}

// GenerateCode generates material code based on type
func (m *Material) GenerateCode(sequence int) {
	if m.Code == "" {
		prefix := "MAT"
		switch m.MaterialType {
		case MaterialTypeRaw:
			prefix = "RM"
		case MaterialTypePackaging:
			prefix = "PKG"
		case MaterialTypeConsumable:
			prefix = "CON"
		case MaterialTypeSemiFinished:
			prefix = "SF"
		}
		m.Code = fmt.Sprintf("%s-%04d", prefix, sequence)
	}
}

// MaterialSpecification represents extended specifications for a material
type MaterialSpecification struct {
	ID               uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	MaterialID       uuid.UUID  `gorm:"type:uuid;not null" json:"material_id"`
	
	// Specifications
	SpecName  string   `gorm:"type:varchar(255);not null" json:"spec_name"`
	SpecValue string   `gorm:"type:text" json:"spec_value,omitempty"`
	SpecUnit  string   `gorm:"type:varchar(50)" json:"spec_unit,omitempty"`
	MinValue  *float64 `gorm:"type:decimal(18,6)" json:"min_value,omitempty"`
	MaxValue  *float64 `gorm:"type:decimal(18,6)" json:"max_value,omitempty"`
	
	// Certificates
	CertificateType       string     `gorm:"type:varchar(100)" json:"certificate_type,omitempty"`
	CertificateNumber     string     `gorm:"type:varchar(100)" json:"certificate_number,omitempty"`
	CertificateIssuer     string     `gorm:"type:varchar(255)" json:"certificate_issuer,omitempty"`
	CertificateExpiryDate *time.Time `gorm:"type:date" json:"certificate_expiry_date,omitempty"`
	CertificateFileURL    string     `gorm:"type:varchar(500)" json:"certificate_file_url,omitempty"`
	
	// Quality
	QualityGrade     string   `gorm:"type:varchar(50)" json:"quality_grade,omitempty"`
	PurityPercentage *float64 `gorm:"type:decimal(5,2)" json:"purity_percentage,omitempty"`
	TestMethod       string   `gorm:"type:varchar(255)" json:"test_method,omitempty"`
	
	// Custom attributes
	CustomAttributes map[string]interface{} `gorm:"type:jsonb" json:"custom_attributes,omitempty"`
	
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid" json:"updated_by,omitempty"`
}

// TableName specifies the table name
func (MaterialSpecification) TableName() string {
	return "material_specifications"
}
