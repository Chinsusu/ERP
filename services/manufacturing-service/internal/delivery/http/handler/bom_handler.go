package handler

import (
	"fmt"

	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/bom"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BOMHandler handles BOM-related requests
type BOMHandler struct {
	createBOMUC    *bom.CreateBOMUseCase
	getBOMUC       *bom.GetBOMUseCase
	listBOMsUC     *bom.ListBOMsUseCase
	approveBOMUC   *bom.ApproveBOMUseCase
	getActiveBOMUC *bom.GetActiveBOMUseCase
}

// NewBOMHandler creates a new BOMHandler
func NewBOMHandler(
	createBOMUC *bom.CreateBOMUseCase,
	getBOMUC *bom.GetBOMUseCase,
	listBOMsUC *bom.ListBOMsUseCase,
	approveBOMUC *bom.ApproveBOMUseCase,
	getActiveBOMUC *bom.GetActiveBOMUseCase,
) *BOMHandler {
	return &BOMHandler{
		createBOMUC:    createBOMUC,
		getBOMUC:       getBOMUC,
		listBOMsUC:     listBOMsUC,
		approveBOMUC:   approveBOMUC,
		getActiveBOMUC: getActiveBOMUC,
	}
}

// CreateBOM creates a new BOM
func (h *BOMHandler) CreateBOM(c *gin.Context) {
	var req dto.CreateBOMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	// Convert item type
	var items []bom.CreateBOMItemInput
	for i, item := range req.Items {
		itemType := entity.BOMItemTypeMaterial
		if item.ItemType != "" {
			itemType = entity.BOMItemType(item.ItemType)
		}
		items = append(items, bom.CreateBOMItemInput{
			LineNumber:      i + 1,
			MaterialID:      item.MaterialID,
			ItemType:        itemType,
			Quantity:        item.Quantity,
			UOMID:           item.UOMID,
			QuantityMin:     item.QuantityMin,
			QuantityMax:     item.QuantityMax,
			IsCritical:      item.IsCritical,
			ScrapPercentage: item.ScrapPercentage,
			UnitCost:        item.UnitCost,
			Notes:           item.Notes,
		})
	}

	confidentiality := entity.ConfidentialityRestricted
	if req.ConfidentialityLevel != "" {
		confidentiality = entity.ConfidentialityLevel(req.ConfidentialityLevel)
	}

	input := bom.CreateBOMInput{
		BOMNumber:            req.BOMNumber,
		ProductID:            req.ProductID,
		Version:              req.Version,
		Name:                 req.Name,
		Description:          req.Description,
		BatchSize:            req.BatchSize,
		BatchUnitID:          req.BatchUnitID,
		ConfidentialityLevel: confidentiality,
		LaborCost:            req.LaborCost,
		OverheadCost:         req.OverheadCost,
		FormulaDetails:       req.FormulaDetails,
		Items:                items,
		CreatedBy:            userID,
	}

	result, err := h.createBOMUC.Execute(c.Request.Context(), input)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	created(c, toBOMResponse(result, nil, false))
}

// GetBOM gets a BOM by ID
func (h *BOMHandler) GetBOM(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid BOM ID")
		return
	}

	// Check if user can view formula (would check permission via auth service)
	canViewFormula := true // TODO: check permission manufacturing:bom:formula_view

	result, err := h.getBOMUC.Execute(c.Request.Context(), id, canViewFormula)
	if err != nil {
		notFound(c, "BOM not found")
		return
	}

	success(c, toBOMResponse(result.BOM, result.FormulaDetails, result.CanViewFormula))
}

// ListBOMs lists BOMs
func (h *BOMHandler) ListBOMs(c *gin.Context) {
	filter := repository.BOMFilter{
		Page:     getPageFromQuery(c),
		PageSize: getPageSizeFromQuery(c),
		Search:   c.Query("search"),
	}

	if status := c.Query("status"); status != "" {
		s := entity.BOMStatus(status)
		filter.Status = &s
	}
	if productID := c.Query("product_id"); productID != "" {
		if id, err := uuid.Parse(productID); err == nil {
			filter.ProductID = &id
		}
	}

	boms, total, err := h.listBOMsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	var items []dto.BOMResponse
	for _, b := range boms {
		items = append(items, toBOMResponse(b, nil, false))
	}

	successWithMeta(c, items, newMeta(filter.Page, filter.PageSize, total))
}

// ApproveBOM approves a BOM
func (h *BOMHandler) ApproveBOM(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid BOM ID")
		return
	}

	userID := getUserIDFromContext(c)

	result, err := h.approveBOMUC.Execute(c.Request.Context(), id, userID)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	success(c, toBOMResponse(result, nil, false))
}

// Helper functions
func toBOMResponse(b *entity.BOM, formula *entity.FormulaDetails, canViewFormula bool) dto.BOMResponse {
	resp := dto.BOMResponse{
		ID:                   b.ID,
		BOMNumber:            b.BOMNumber,
		ProductID:            b.ProductID,
		Version:              b.Version,
		Name:                 b.Name,
		Description:          b.Description,
		Status:               string(b.Status),
		BatchSize:            b.BatchSize,
		ConfidentialityLevel: string(b.ConfidentialityLevel),
		MaterialCost:         b.MaterialCost,
		LaborCost:            b.LaborCost,
		OverheadCost:         b.OverheadCost,
		TotalCost:            b.TotalCost,
		EffectiveFrom:        b.EffectiveFrom,
		EffectiveTo:          b.EffectiveTo,
		CreatedAt:            b.CreatedAt,
	}

	for _, item := range b.Items {
		resp.Items = append(resp.Items, dto.BOMItemResponse{
			ID:         item.ID,
			LineNumber: item.LineNumber,
			MaterialID: item.MaterialID,
			ItemType:   string(item.ItemType),
			Quantity:   item.Quantity,
			UOMID:      item.UOMID,
			IsCritical: item.IsCritical,
			UnitCost:   item.UnitCost,
			TotalCost:  item.TotalCost,
			Notes:      item.Notes,
		})
	}

	if canViewFormula && formula != nil {
		resp.FormulaDetails = formula
	} else if !canViewFormula && b.ConfidentialityLevel == entity.ConfidentialityRestricted {
		resp.Message = "Full BOM details restricted. Contact R&D Manager."
	}

	return resp
}

func getUserIDFromContext(c *gin.Context) uuid.UUID {
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if id, err := uuid.Parse(userIDStr); err == nil {
			return id
		}
	}
	return uuid.New() // fallback for development
}

func getPageFromQuery(c *gin.Context) int {
	page := 1
	if p := c.Query("page"); p != "" {
		if n, err := parseInt(p); err == nil && n > 0 {
			page = n
		}
	}
	return page
}

func getPageSizeFromQuery(c *gin.Context) int {
	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if n, err := parseInt(ps); err == nil && n > 0 && n <= 100 {
			pageSize = n
		}
	}
	return pageSize
}

func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
