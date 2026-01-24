package dto

import (
	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/google/uuid"
)

// ======== Category DTOs ========

type CreateCategoryRequest struct {
	Code         string `json:"code" validate:"required"`
	Name         string `json:"name" validate:"required"`
	NameEN       string `json:"name_en,omitempty"`
	Description  string `json:"description,omitempty"`
	CategoryType string `json:"category_type" validate:"required,oneof=MATERIAL PRODUCT"`
	ParentID     string `json:"parent_id,omitempty"`
	SortOrder    int    `json:"sort_order,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	NameEN      string `json:"name_en,omitempty"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CategoryResponse struct {
	ID           string              `json:"id"`
	Code         string              `json:"code"`
	Name         string              `json:"name"`
	NameEN       string              `json:"name_en,omitempty"`
	Description  string              `json:"description,omitempty"`
	CategoryType string              `json:"category_type"`
	ParentID     *string             `json:"parent_id,omitempty"`
	Path         string              `json:"path"`
	Level        int                 `json:"level"`
	Status       string              `json:"status"`
	Children     []CategoryResponse  `json:"children,omitempty"`
}

func ToCategoryResponse(c *entity.Category) CategoryResponse {
	resp := CategoryResponse{
		ID:           c.ID.String(),
		Code:         c.Code,
		Name:         c.Name,
		NameEN:       c.NameEN,
		Description:  c.Description,
		CategoryType: string(c.CategoryType),
		Path:         c.Path,
		Level:        c.Level,
		Status:       c.Status,
	}
	if c.ParentID != nil {
		s := c.ParentID.String()
		resp.ParentID = &s
	}
	for _, child := range c.Children {
		resp.Children = append(resp.Children, ToCategoryResponse(&child))
	}
	return resp
}

// ======== Unit DTOs ========

type CreateUnitRequest struct {
	Code             string  `json:"code" validate:"required"`
	Name             string  `json:"name" validate:"required"`
	NameEN           string  `json:"name_en,omitempty"`
	Symbol           string  `json:"symbol" validate:"required"`
	UoMType          string  `json:"uom_type" validate:"required,oneof=WEIGHT VOLUME QUANTITY LENGTH AREA"`
	IsBaseUnit       bool    `json:"is_base_unit,omitempty"`
	BaseUnitID       string  `json:"base_unit_id,omitempty"`
	ConversionFactor float64 `json:"conversion_factor,omitempty"`
	Description      string  `json:"description,omitempty"`
}

type ConvertUnitRequest struct {
	Value      float64 `json:"value" validate:"required"`
	FromUnitID string  `json:"from_unit_id" validate:"required,uuid"`
	ToUnitID   string  `json:"to_unit_id" validate:"required,uuid"`
}

type ConvertUnitResponse struct {
	OriginalValue  float64 `json:"original_value"`
	OriginalUoM    string  `json:"original_uom"`
	ConvertedValue float64 `json:"converted_value"`
	ConvertedUoM   string  `json:"converted_uom"`
}

type UnitResponse struct {
	ID               string  `json:"id"`
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	NameEN           string  `json:"name_en,omitempty"`
	Symbol           string  `json:"symbol"`
	UoMType          string  `json:"uom_type"`
	IsBaseUnit       bool    `json:"is_base_unit"`
	ConversionFactor float64 `json:"conversion_factor"`
	Status           string  `json:"status"`
}

func ToUnitResponse(u *entity.UnitOfMeasure) UnitResponse {
	return UnitResponse{
		ID:               u.ID.String(),
		Code:             u.Code,
		Name:             u.Name,
		NameEN:           u.NameEN,
		Symbol:           u.Symbol,
		UoMType:          string(u.UoMType),
		IsBaseUnit:       u.IsBaseUnit,
		ConversionFactor: u.ConversionFactor,
		Status:           u.Status,
	}
}

// ======== Material DTOs ========

type CreateMaterialRequest struct {
	Code             string   `json:"code,omitempty"`
	Name             string   `json:"name" validate:"required"`
	NameEN           string   `json:"name_en,omitempty"`
	Description      string   `json:"description,omitempty"`
	MaterialType     string   `json:"material_type" validate:"required,oneof=RAW_MATERIAL PACKAGING CONSUMABLE SEMI_FINISHED"`
	CategoryID       string   `json:"category_id,omitempty"`
	BaseUnitID       string   `json:"base_unit_id" validate:"required,uuid"`
	INCIName         string   `json:"inci_name,omitempty"`
	CASNumber        string   `json:"cas_number,omitempty"`
	IsAllergen       bool     `json:"is_allergen,omitempty"`
	AllergenInfo     string   `json:"allergen_info,omitempty"`
	IsOrganic        bool     `json:"is_organic,omitempty"`
	IsNatural        bool     `json:"is_natural,omitempty"`
	IsVegan          bool     `json:"is_vegan,omitempty"`
	OriginCountry    string   `json:"origin_country,omitempty"`
	StorageCondition string   `json:"storage_condition,omitempty"`
	MinTemp          *float64 `json:"min_temp,omitempty"`
	MaxTemp          *float64 `json:"max_temp,omitempty"`
	ShelfLifeDays    int      `json:"shelf_life_days,omitempty"`
	IsHazardous      bool     `json:"is_hazardous,omitempty"`
	LeadTimeDays     int      `json:"lead_time_days,omitempty"`
	MinOrderQty      float64  `json:"min_order_qty,omitempty"`
	ReorderPoint     float64  `json:"reorder_point,omitempty"`
	SafetyStock      float64  `json:"safety_stock,omitempty"`
	StandardCost     float64  `json:"standard_cost,omitempty"`
	Currency         string   `json:"currency,omitempty"`
}

type MaterialResponse struct {
	ID               string           `json:"id"`
	Code             string           `json:"code"`
	Name             string           `json:"name"`
	NameEN           string           `json:"name_en,omitempty"`
	MaterialType     string           `json:"material_type"`
	Category         *CategoryResponse `json:"category,omitempty"`
	BaseUnit         *UnitResponse    `json:"base_unit,omitempty"`
	INCIName         string           `json:"inci_name,omitempty"`
	CASNumber        string           `json:"cas_number,omitempty"`
	IsAllergen       bool             `json:"is_allergen"`
	IsOrganic        bool             `json:"is_organic"`
	IsNatural        bool             `json:"is_natural"`
	StorageCondition string           `json:"storage_condition"`
	ShelfLifeDays    int              `json:"shelf_life_days"`
	StandardCost     float64          `json:"standard_cost"`
	Currency         string           `json:"currency"`
	Status           string           `json:"status"`
}

func ToMaterialResponse(m *entity.Material) MaterialResponse {
	resp := MaterialResponse{
		ID:               m.ID.String(),
		Code:             m.Code,
		Name:             m.Name,
		NameEN:           m.NameEN,
		MaterialType:     string(m.MaterialType),
		INCIName:         m.INCIName,
		CASNumber:        m.CASNumber,
		IsAllergen:       m.IsAllergen,
		IsOrganic:        m.IsOrganic,
		IsNatural:        m.IsNatural,
		StorageCondition: string(m.StorageCondition),
		ShelfLifeDays:    m.ShelfLifeDays,
		StandardCost:     m.StandardCost,
		Currency:         m.Currency,
		Status:           m.Status,
	}
	if m.Category != nil {
		cat := ToCategoryResponse(m.Category)
		resp.Category = &cat
	}
	if m.BaseUnit != nil {
		unit := ToUnitResponse(m.BaseUnit)
		resp.BaseUnit = &unit
	}
	return resp
}

// ======== Product DTOs ========

type CreateProductRequest struct {
	Code                  string   `json:"code,omitempty"`
	SKU                   string   `json:"sku" validate:"required"`
	Barcode               string   `json:"barcode,omitempty"`
	Name                  string   `json:"name" validate:"required"`
	NameEN                string   `json:"name_en,omitempty"`
	Description           string   `json:"description,omitempty"`
	CategoryID            string   `json:"category_id,omitempty"`
	ProductLine           string   `json:"product_line,omitempty"`
	Brand                 string   `json:"brand,omitempty"`
	Volume                *float64 `json:"volume,omitempty"`
	VolumeUnit            string   `json:"volume_unit,omitempty"`
	CosmeticLicenseNumber string   `json:"cosmetic_license_number,omitempty"`
	IngredientsSummary    string   `json:"ingredients_summary,omitempty"`
	TargetSkinType        string   `json:"target_skin_type,omitempty"`
	PackagingType         string   `json:"packaging_type,omitempty"`
	StandardCost          float64  `json:"standard_cost,omitempty"`
	StandardPrice         float64  `json:"standard_price,omitempty"`
	Currency              string   `json:"currency,omitempty"`
	BaseUnitID            string   `json:"base_unit_id" validate:"required,uuid"`
	ShelfLifeMonths       int      `json:"shelf_life_months,omitempty"`
}

type ProductResponse struct {
	ID                    string            `json:"id"`
	Code                  string            `json:"code"`
	SKU                   string            `json:"sku"`
	Barcode               string            `json:"barcode,omitempty"`
	Name                  string            `json:"name"`
	NameEN                string            `json:"name_en,omitempty"`
	Category              *CategoryResponse `json:"category,omitempty"`
	ProductLine           string            `json:"product_line,omitempty"`
	Brand                 string            `json:"brand,omitempty"`
	Volume                *float64          `json:"volume,omitempty"`
	VolumeUnit            string            `json:"volume_unit,omitempty"`
	CosmeticLicenseNumber string            `json:"cosmetic_license_number,omitempty"`
	IngredientsSummary    string            `json:"ingredients_summary,omitempty"`
	StandardPrice         float64           `json:"standard_price"`
	Currency              string            `json:"currency"`
	Status                string            `json:"status"`
	Images                []ImageResponse   `json:"images,omitempty"`
}

type ImageResponse struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}

func ToProductResponse(p *entity.Product) ProductResponse {
	resp := ProductResponse{
		ID:                    p.ID.String(),
		Code:                  p.Code,
		SKU:                   p.SKU,
		Barcode:               p.Barcode,
		Name:                  p.Name,
		NameEN:                p.NameEN,
		ProductLine:           p.ProductLine,
		Brand:                 p.Brand,
		Volume:                p.Volume,
		VolumeUnit:            p.VolumeUnit,
		CosmeticLicenseNumber: p.CosmeticLicenseNumber,
		IngredientsSummary:    p.IngredientsSummary,
		StandardPrice:         p.StandardPrice,
		Currency:              p.Currency,
		Status:                p.Status,
	}
	if p.Category != nil {
		cat := ToCategoryResponse(p.Category)
		resp.Category = &cat
	}
	for _, img := range p.Images {
		resp.Images = append(resp.Images, ImageResponse{
			ID:        img.ID.String(),
			ImageURL:  img.ImageURL,
			IsPrimary: img.IsPrimary,
			SortOrder: img.SortOrder,
		})
	}
	return resp
}

// ======== Common DTOs ========

type ListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
	Search   string `form:"search"`
	Status   string `form:"status"`
}

type MaterialListQuery struct {
	ListQuery
	MaterialType     string `form:"material_type"`
	CategoryID       string `form:"category_id"`
	StorageCondition string `form:"storage_condition"`
	IsOrganic        *bool  `form:"is_organic"`
	IsNatural        *bool  `form:"is_natural"`
}

type ProductListQuery struct {
	ListQuery
	CategoryID  string `form:"category_id"`
	ProductLine string `form:"product_line"`
	Brand       string `form:"brand"`
}

// Helper to parse UUID
func ParseUUID(s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
