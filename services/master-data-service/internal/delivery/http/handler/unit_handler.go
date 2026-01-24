package handler

import (
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	unitUC "github.com/erp-cosmetics/master-data-service/internal/usecase/unit"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UnitHandler handles unit HTTP requests
type UnitHandler struct {
	uc *unitUC.UseCase
}

// NewUnitHandler creates a new unit handler
func NewUnitHandler(uc *unitUC.UseCase) *UnitHandler {
	return &UnitHandler{uc: uc}
}

// Create creates a new unit
func (h *UnitHandler) Create(c *gin.Context) {
	var req dto.CreateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	var baseUnitID *uuid.UUID
	if req.BaseUnitID != "" {
		id, err := uuid.Parse(req.BaseUnitID)
		if err != nil {
			response.Error(c, errors.BadRequest("Invalid base_unit_id"))
			return
		}
		baseUnitID = &id
	}

	unit, err := h.uc.Create(c.Request.Context(), &unitUC.CreateRequest{
		Code:             req.Code,
		Name:             req.Name,
		NameEN:           req.NameEN,
		Symbol:           req.Symbol,
		UoMType:          entity.UoMType(req.UoMType),
		IsBaseUnit:       req.IsBaseUnit,
		BaseUnitID:       baseUnitID,
		ConversionFactor: req.ConversionFactor,
		Description:      req.Description,
	})
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, dto.ToUnitResponse(unit))
}

// GetByID retrieves a unit by ID
func (h *UnitHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid unit ID"))
		return
	}

	unit, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Unit"))
		return
	}

	response.Success(c, dto.ToUnitResponse(unit))
}

// List lists units
func (h *UnitHandler) List(c *gin.Context) {
	var query dto.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	uomType := entity.UoMType(c.Query("uom_type"))

	filter := &repository.UnitFilter{
		UoMType:  uomType,
		Status:   query.Status,
		Search:   query.Search,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	units, total, err := h.uc.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.UnitResponse
	for _, unit := range units {
		resp = append(resp, dto.ToUnitResponse(&unit))
	}

	meta := response.NewMeta(query.Page, query.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// Convert converts a value between units
func (h *UnitHandler) Convert(c *gin.Context) {
	var req dto.ConvertUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	fromUnitID, err := uuid.Parse(req.FromUnitID)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid from_unit_id"))
		return
	}

	toUnitID, err := uuid.Parse(req.ToUnitID)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid to_unit_id"))
		return
	}

	result, err := h.uc.Convert(c.Request.Context(), &unitUC.ConvertRequest{
		Value:      req.Value,
		FromUnitID: fromUnitID,
		ToUnitID:   toUnitID,
	})
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ConvertUnitResponse{
		OriginalValue:  result.OriginalValue,
		OriginalUoM:    result.OriginalUoM,
		ConvertedValue: result.ConvertedValue,
		ConvertedUoM:   result.ConvertedUoM,
	})
}

// GetConversions gets all conversions for a unit
func (h *UnitHandler) GetConversions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid unit ID"))
		return
	}

	conversions, err := h.uc.GetConversions(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, conversions)
}
