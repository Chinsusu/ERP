package handler

import (
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/usecase/inventory"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InventoryCountHandler handles inventory count endpoints
type InventoryCountHandler struct {
	createCountUC   *inventory.CreateInventoryCountUseCase
	startCountUC    *inventory.StartInventoryCountUseCase
	recordCountUC   *inventory.RecordCountUseCase
	completeCountUC *inventory.CompleteInventoryCountUseCase
	getCountUC      *inventory.GetInventoryCountUseCase
	listCountsUC    *inventory.ListInventoryCountsUseCase
}

// NewInventoryCountHandler creates a new handler
func NewInventoryCountHandler(
	createCountUC *inventory.CreateInventoryCountUseCase,
	startCountUC *inventory.StartInventoryCountUseCase,
	recordCountUC *inventory.RecordCountUseCase,
	completeCountUC *inventory.CompleteInventoryCountUseCase,
	getCountUC *inventory.GetInventoryCountUseCase,
	listCountsUC *inventory.ListInventoryCountsUseCase,
) *InventoryCountHandler {
	return &InventoryCountHandler{
		createCountUC:   createCountUC,
		startCountUC:    startCountUC,
		recordCountUC:   recordCountUC,
		completeCountUC: completeCountUC,
		getCountUC:      getCountUC,
		listCountsUC:    listCountsUC,
	}
}

// CreateInventoryCountRequest represents create request
type CreateInventoryCountRequest struct {
	CountDate   string     `json:"count_date" binding:"required"`
	CountType   string     `json:"count_type" binding:"required,oneof=FULL CYCLE SPOT"`
	WarehouseID uuid.UUID  `json:"warehouse_id" binding:"required"`
	ZoneID      *uuid.UUID `json:"zone_id"`
	Notes       string     `json:"notes"`
}

// CreateInventoryCount handles POST /inventory-counts
func (h *InventoryCountHandler) CreateInventoryCount(c *gin.Context) {
	var req CreateInventoryCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	countDate, err := time.Parse("2006-01-02", req.CountDate)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid count date format"))
		return
	}

	userID := uuid.New() // Placeholder

	input := &inventory.CreateInventoryCountInput{
		CountDate:   countDate,
		CountType:   entity.InventoryCountType(req.CountType),
		WarehouseID: req.WarehouseID,
		ZoneID:      req.ZoneID,
		Notes:       req.Notes,
		CreatedBy:   userID,
	}

	result, err := h.createCountUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, gin.H{
		"id":           result.ID,
		"count_number": result.CountNumber,
		"status":       result.Status,
	})
}

// GetInventoryCount handles GET /inventory-counts/:id
func (h *InventoryCountHandler) GetInventoryCount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid count ID"))
		return
	}

	result, err := h.getCountUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Inventory Count"))
		return
	}

	response.Success(c, result)
}

// ListInventoryCounts handles GET /inventory-counts
func (h *InventoryCountHandler) ListInventoryCounts(c *gin.Context) {
	filter := &repository.InventoryCountFilter{
		Status:    c.Query("status"),
		CountType: c.Query("count_type"),
		Search:    c.Query("search"),
		Page:      getPageParam(c),
		Limit:     getLimitParam(c),
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		id, _ := uuid.Parse(warehouseID)
		filter.WarehouseID = &id
	}

	counts, total, err := h.listCountsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.SuccessWithMeta(c, counts, response.NewMeta(filter.Page, filter.Limit, int64(total)))
}

// StartInventoryCount handles PATCH /inventory-counts/:id/start
func (h *InventoryCountHandler) StartInventoryCount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid count ID"))
		return
	}

	result, err := h.startCountUC.Execute(c.Request.Context(), id)
	if err != nil {
		if err == entity.ErrInvalidStatus {
			response.Error(c, errors.BadRequest("Cannot start count in current status"))
			return
		}
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, gin.H{
		"id":     result.ID,
		"status": result.Status,
	})
}

// RecordCountRequest represents record count request
type RecordCountRequest struct {
	LineItemID uuid.UUID `json:"line_item_id" binding:"required"`
	CountedQty float64   `json:"counted_qty" binding:"required"`
	Notes      string    `json:"notes"`
}

// RecordCount handles POST /inventory-counts/:id/record
func (h *InventoryCountHandler) RecordCount(c *gin.Context) {
	var req RecordCountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	userID := uuid.New() // Placeholder

	input := &inventory.RecordCountInput{
		LineItemID: req.LineItemID,
		CountedQty: req.CountedQty,
		CountedBy:  userID,
		Notes:      req.Notes,
	}

	result, err := h.recordCountUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, gin.H{
		"id":        result.ID,
		"variance":  result.Variance,
		"is_counted": result.IsCounted,
	})
}

// CompleteInventoryCountRequest represents complete request
type CompleteInventoryCountRequest struct {
	ApplyVariance bool `json:"apply_variance"`
}

// CompleteInventoryCount handles PATCH /inventory-counts/:id/complete
func (h *InventoryCountHandler) CompleteInventoryCount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid count ID"))
		return
	}

	var req CompleteInventoryCountRequest
	c.ShouldBindJSON(&req)

	userID := uuid.New() // Placeholder

	input := &inventory.CompleteInventoryCountInput{
		CountID:       id,
		ApplyVariance: req.ApplyVariance,
		ApprovedBy:    userID,
	}

	result, err := h.completeCountUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err == entity.ErrInvalidStatus {
			response.Error(c, errors.BadRequest("Cannot complete count in current status"))
			return
		}
		if err == entity.ErrPendingItems {
			response.Error(c, errors.BadRequest("All items must be counted before completion"))
			return
		}
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, gin.H{
		"id":     result.ID,
		"status": result.Status,
	})
}
