package handler

import (
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	materialUC "github.com/erp-cosmetics/master-data-service/internal/usecase/material"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MaterialHandler handles material HTTP requests
type MaterialHandler struct {
	uc *materialUC.UseCase
}

// NewMaterialHandler creates a new material handler
func NewMaterialHandler(uc *materialUC.UseCase) *MaterialHandler {
	return &MaterialHandler{uc: uc}
}

// Create creates a new material
func (h *MaterialHandler) Create(c *gin.Context) {
	var req dto.CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	baseUnitID, err := uuid.Parse(req.BaseUnitID)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid base_unit_id"))
		return
	}

	categoryID, _ := dto.ParseUUID(req.CategoryID)

	material, err := h.uc.Create(c.Request.Context(), &materialUC.CreateRequest{
		Code:             req.Code,
		Name:             req.Name,
		NameEN:           req.NameEN,
		Description:      req.Description,
		MaterialType:     entity.MaterialType(req.MaterialType),
		CategoryID:       categoryID,
		BaseUnitID:       baseUnitID,
		INCIName:         req.INCIName,
		CASNumber:        req.CASNumber,
		IsAllergen:       req.IsAllergen,
		AllergenInfo:     req.AllergenInfo,
		IsOrganic:        req.IsOrganic,
		IsNatural:        req.IsNatural,
		IsVegan:          req.IsVegan,
		OriginCountry:    req.OriginCountry,
		StorageCondition: entity.StorageCondition(req.StorageCondition),
		MinTemp:          req.MinTemp,
		MaxTemp:          req.MaxTemp,
		ShelfLifeDays:    req.ShelfLifeDays,
		IsHazardous:      req.IsHazardous,
		LeadTimeDays:     req.LeadTimeDays,
		MinOrderQty:      req.MinOrderQty,
		ReorderPoint:     req.ReorderPoint,
		SafetyStock:      req.SafetyStock,
		StandardCost:     req.StandardCost,
		Currency:         req.Currency,
	})
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, dto.ToMaterialResponse(material))
}

// GetByID retrieves a material by ID
func (h *MaterialHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	material, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Material"))
		return
	}

	response.Success(c, dto.ToMaterialResponse(material))
}

// Update updates an existing material
func (h *MaterialHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	material, err := h.uc.Update(c.Request.Context(), id, updates)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToMaterialResponse(material))
}

// Delete soft deletes a material
func (h *MaterialHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, gin.H{"message": "Material deleted successfully"})
}

// List lists materials
func (h *MaterialHandler) List(c *gin.Context) {
	var query dto.MaterialListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	categoryID, _ := dto.ParseUUID(query.CategoryID)

	filter := &repository.MaterialFilter{
		MaterialType:     entity.MaterialType(query.MaterialType),
		CategoryID:       categoryID,
		StorageCondition: entity.StorageCondition(query.StorageCondition),
		IsOrganic:        query.IsOrganic,
		IsNatural:        query.IsNatural,
		Status:           query.Status,
		Search:           query.Search,
		Page:             query.Page,
		PageSize:         query.PageSize,
	}

	materials, total, err := h.uc.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.MaterialResponse
	for _, mat := range materials {
		resp = append(resp, dto.ToMaterialResponse(&mat))
	}

	meta := response.NewMeta(query.Page, query.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// Search searches materials
func (h *MaterialHandler) Search(c *gin.Context) {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		h.List(c)
		return
	}

	var listQuery dto.MaterialListQuery
	if err := c.ShouldBindQuery(&listQuery); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	categoryID, _ := dto.ParseUUID(listQuery.CategoryID)

	filter := &repository.MaterialFilter{
		MaterialType:     entity.MaterialType(listQuery.MaterialType),
		CategoryID:       categoryID,
		StorageCondition: entity.StorageCondition(listQuery.StorageCondition),
		IsOrganic:        listQuery.IsOrganic,
		IsNatural:        listQuery.IsNatural,
		Status:           listQuery.Status,
		Page:             listQuery.Page,
		PageSize:         listQuery.PageSize,
	}

	materials, total, err := h.uc.Search(c.Request.Context(), searchQuery, filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.MaterialResponse
	for _, mat := range materials {
		resp = append(resp, dto.ToMaterialResponse(&mat))
	}

	meta := response.NewMeta(listQuery.Page, listQuery.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// AddSpecification adds a specification to a material
func (h *MaterialHandler) AddSpecification(c *gin.Context) {
	materialID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	var spec entity.MaterialSpecification
	if err := c.ShouldBindJSON(&spec); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	if err := h.uc.AddSpecification(c.Request.Context(), materialID, &spec); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, spec)
}

// GetSpecifications gets specifications for a material
func (h *MaterialHandler) GetSpecifications(c *gin.Context) {
	materialID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	specs, err := h.uc.GetSpecifications(c.Request.Context(), materialID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, specs)
}
