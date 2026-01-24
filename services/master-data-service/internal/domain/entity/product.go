package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a finished cosmetic product
type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	SKU         string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"sku"`
	Barcode     string         `gorm:"type:varchar(50)" json:"barcode,omitempty"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	NameEN      string         `gorm:"type:varchar(255);column:name_en" json:"name_en,omitempty"`
	Description string         `gorm:"type:text" json:"description,omitempty"`

	// Classification
	CategoryID  *uuid.UUID `gorm:"type:uuid" json:"category_id,omitempty"`
	ProductLine string     `gorm:"type:varchar(100)" json:"product_line,omitempty"`
	Brand       string     `gorm:"type:varchar(100)" json:"brand,omitempty"`

	// Physical properties
	Volume     *float64 `gorm:"type:decimal(10,2)" json:"volume,omitempty"`
	VolumeUnit string   `gorm:"type:varchar(20)" json:"volume_unit,omitempty"`
	Weight     *float64 `gorm:"type:decimal(10,2)" json:"weight,omitempty"`
	WeightUnit string   `gorm:"type:varchar(20)" json:"weight_unit,omitempty"`

	// Cosmetics regulatory
	CosmeticLicenseNumber string     `gorm:"type:varchar(100)" json:"cosmetic_license_number,omitempty"`
	LicenseExpiryDate     *time.Time `gorm:"type:date" json:"license_expiry_date,omitempty"`
	RegistrationCountry   string     `gorm:"type:varchar(100)" json:"registration_country,omitempty"`

	// Product info
	IngredientsSummary string `gorm:"type:text" json:"ingredients_summary,omitempty"`
	TargetSkinType     string `gorm:"type:varchar(255)" json:"target_skin_type,omitempty"`
	UsageInstructions  string `gorm:"type:text" json:"usage_instructions,omitempty"`
	Warnings           string `gorm:"type:text" json:"warnings,omitempty"`

	// Packaging
	PackagingType     string `gorm:"type:varchar(100)" json:"packaging_type,omitempty"`
	PackagingMaterial string `gorm:"type:varchar(100)" json:"packaging_material,omitempty"`

	// Pricing
	StandardCost            float64  `gorm:"type:decimal(18,4);not null;default:0" json:"standard_cost"`
	StandardPrice           float64  `gorm:"type:decimal(18,4);not null;default:0" json:"standard_price"`
	RecommendedRetailPrice  *float64 `gorm:"type:decimal(18,4)" json:"recommended_retail_price,omitempty"`
	Currency                string   `gorm:"type:varchar(3);not null;default:'VND'" json:"currency"`

	// Units
	BaseUnitID  uuid.UUID  `gorm:"type:uuid;not null" json:"base_unit_id"`
	SalesUnitID *uuid.UUID `gorm:"type:uuid" json:"sales_unit_id,omitempty"`

	// Shelf life
	ShelfLifeMonths int `gorm:"not null;default:24" json:"shelf_life_months"`

	// Launch & Status
	LaunchDate      *time.Time     `gorm:"type:date" json:"launch_date,omitempty"`
	DiscontinueDate *time.Time     `gorm:"type:date" json:"discontinue_date,omitempty"`
	Status          string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// Audit
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Category  *Category       `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	BaseUnit  *UnitOfMeasure  `gorm:"foreignKey:BaseUnitID" json:"base_unit,omitempty"`
	SalesUnit *UnitOfMeasure  `gorm:"foreignKey:SalesUnitID" json:"sales_unit,omitempty"`
	Images    []ProductImage  `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

// TableName specifies the table name
func (Product) TableName() string {
	return "products"
}

// Validate validates product data
func (p *Product) Validate() error {
	if p.Code == "" {
		return fmt.Errorf("code is required")
	}
	if p.SKU == "" {
		return fmt.Errorf("sku is required")
	}
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.BaseUnitID == uuid.Nil {
		return fmt.Errorf("base_unit_id is required")
	}
	return nil
}

// IsActive checks if product is active
func (p *Product) IsActive() bool {
	return p.Status == "active"
}

// IsLicenseExpiring checks if license is expiring within days
func (p *Product) IsLicenseExpiring(days int) bool {
	if p.LicenseExpiryDate == nil {
		return false
	}
	return time.Until(*p.LicenseExpiryDate).Hours()/24 <= float64(days)
}

// GenerateCode generates product code based on category
func (p *Product) GenerateCode(categoryCode string, sequence int) {
	if p.Code == "" {
		if categoryCode == "" {
			categoryCode = "PROD"
		}
		p.Code = fmt.Sprintf("FG-%s-%04d", categoryCode, sequence)
	}
}

// ProductImage represents an image for a product
type ProductImage struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProductID    uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	ImageURL     string    `gorm:"type:varchar(500);not null" json:"image_url"`
	ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url,omitempty"`
	AltText      string    `gorm:"type:varchar(255)" json:"alt_text,omitempty"`
	SortOrder    int       `gorm:"not null;default:0" json:"sort_order"`
	IsPrimary    bool      `gorm:"not null;default:false" json:"is_primary"`
	ImageType    string    `gorm:"type:varchar(50)" json:"image_type,omitempty"` // FRONT, BACK, SIDE, DETAIL
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name
func (ProductImage) TableName() string {
	return "product_images"
}
