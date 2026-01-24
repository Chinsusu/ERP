package handler

import (
	"time"

	"github.com/erp-cosmetics/wms-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/usecase/adjustment"
	"github.com/erp-cosmetics/wms-service/internal/usecase/reservation"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReservationHandler handles reservation endpoints
type ReservationHandler struct {
	createReservationUC   *reservation.CreateReservationUseCase
	releaseReservationUC  *reservation.ReleaseReservationUseCase
	checkAvailabilityUC   *reservation.CheckAvailabilityUseCase
}

// NewReservationHandler creates a new reservation handler
func NewReservationHandler(
	createReservationUC *reservation.CreateReservationUseCase,
	releaseReservationUC *reservation.ReleaseReservationUseCase,
	checkAvailabilityUC *reservation.CheckAvailabilityUseCase,
) *ReservationHandler {
	return &ReservationHandler{
		createReservationUC:  createReservationUC,
		releaseReservationUC: releaseReservationUC,
		checkAvailabilityUC:  checkAvailabilityUC,
	}
}

// CreateReservation handles POST /reservations
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
	var req dto.ReserveStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t, _ := time.Parse(time.RFC3339, *req.ExpiresAt)
		expiresAt = &t
	}

	userID := uuid.New() // Placeholder

	input := &reservation.CreateReservationInput{
		MaterialID:      req.MaterialID,
		Quantity:        req.Quantity,
		UnitID:          req.UnitID,
		ReservationType: entity.ReservationType(req.ReservationType),
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: req.ReferenceNumber,
		ExpiresAt:       expiresAt,
		CreatedBy:       userID,
	}

	result, err := h.createReservationUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err == entity.ErrInsufficientStock {
			response.Error(c, errors.BadRequest("Insufficient stock for reservation"))
			return
		}
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, dto.ReserveStockResponse{
		ReservationID:    result.ID,
		ReservedQuantity: result.Quantity,
	})
}

// ReleaseReservation handles DELETE /reservations/:id
func (h *ReservationHandler) ReleaseReservation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid reservation ID"))
		return
	}

	if err := h.releaseReservationUC.Execute(c.Request.Context(), id); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.NoContent(c)
}

// CheckAvailability handles GET /stock/availability/:material_id
func (h *ReservationHandler) CheckAvailability(c *gin.Context) {
	materialID, err := uuid.Parse(c.Param("material_id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	requestedQty := 0.0
	if qty := c.Query("quantity"); qty != "" {
		var q float64
		if _, err := c.GetQuery("quantity"); err {
			// parse the quantity
		}
		requestedQty = q
	}

	result, err := h.checkAvailabilityUC.Execute(c.Request.Context(), materialID, requestedQty)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, result)
}

// AdjustmentHandler handles stock adjustment endpoints
type AdjustmentHandler struct {
	createAdjustmentUC *adjustment.CreateAdjustmentUseCase
	transferStockUC    *adjustment.TransferStockUseCase
}

// NewAdjustmentHandler creates a new adjustment handler
func NewAdjustmentHandler(
	createAdjustmentUC *adjustment.CreateAdjustmentUseCase,
	transferStockUC *adjustment.TransferStockUseCase,
) *AdjustmentHandler {
	return &AdjustmentHandler{
		createAdjustmentUC: createAdjustmentUC,
		transferStockUC:    transferStockUC,
	}
}

// CreateAdjustment handles POST /adjustments
func (h *AdjustmentHandler) CreateAdjustment(c *gin.Context) {
	var req dto.AdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	adjustmentDate, err := time.Parse("2006-01-02", req.AdjustmentDate)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid adjustment date format"))
		return
	}

	userID := uuid.New() // Placeholder

	input := &adjustment.CreateAdjustmentInput{
		AdjustmentDate: adjustmentDate,
		AdjustmentType: adjustment.AdjustmentType(req.AdjustmentType),
		LocationID:     req.LocationID,
		MaterialID:     req.MaterialID,
		LotID:          req.LotID,
		UnitID:         req.UnitID,
		SystemQty:      req.SystemQty,
		ActualQty:      req.ActualQty,
		Reason:         req.Reason,
		Notes:          req.Notes,
		AdjustedBy:     userID,
	}

	result, err := h.createAdjustmentUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, gin.H{
		"adjustment_number": result.AdjustmentNumber,
		"variance":          result.Variance,
		"movement_number":   result.MovementNumber,
	})
}

// TransferStock handles POST /transfers
func (h *AdjustmentHandler) TransferStock(c *gin.Context) {
	var req dto.TransferStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	userID := uuid.New() // Placeholder

	input := &adjustment.TransferStockInput{
		MaterialID:     req.MaterialID,
		LotID:          req.LotID,
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Quantity:       req.Quantity,
		UnitID:         req.UnitID,
		Reason:         req.Reason,
		TransferredBy:  userID,
	}

	movementNumber, err := h.transferStockUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err == entity.ErrInsufficientStock {
			response.Error(c, errors.BadRequest("Insufficient stock for transfer"))
			return
		}
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, gin.H{
		"movement_number": movementNumber,
	})
}
