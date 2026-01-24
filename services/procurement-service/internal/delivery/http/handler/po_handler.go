package handler

import (
	"strconv"

	"github.com/erp-cosmetics/procurement-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/usecase/po"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// POHandler handles PO-related requests
type POHandler struct {
	createPOFromPR  *po.CreatePOFromPRUseCase
	getPO           *po.GetPOUseCase
	listPOs         *po.ListPOsUseCase
	confirmPO       *po.ConfirmPOUseCase
	cancelPO        *po.CancelPOUseCase
	closePO         *po.ClosePOUseCase
	getPOReceipts   *po.GetPOReceiptsUseCase
}

// NewPOHandler creates a new PO handler
func NewPOHandler(
	createPOFromPR *po.CreatePOFromPRUseCase,
	getPO *po.GetPOUseCase,
	listPOs *po.ListPOsUseCase,
	confirmPO *po.ConfirmPOUseCase,
	cancelPO *po.CancelPOUseCase,
	closePO *po.ClosePOUseCase,
	getPOReceipts *po.GetPOReceiptsUseCase,
) *POHandler {
	return &POHandler{
		createPOFromPR: createPOFromPR,
		getPO:          getPO,
		listPOs:        listPOs,
		confirmPO:      confirmPO,
		cancelPO:       cancelPO,
		closePO:        closePO,
		getPOReceipts:  getPOReceipts,
	}
}

// ConvertPRToPO handles POST /api/v1/purchase-requisitions/:id/convert-to-po
func (h *POHandler) ConvertPRToPO(c *gin.Context) {
	prID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PR ID"))
		return
	}

	var req dto.ConvertPRToPORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	supplierID, err := uuid.Parse(req.SupplierID)
	if err != nil {
		response.Error(c, errors.BadRequest("invalid supplier ID"))
		return
	}

	userID := c.GetHeader("X-User-ID")
	createdBy, _ := uuid.Parse(userID)
	if createdBy == uuid.Nil {
		createdBy = uuid.New()
	}

	ucReq := &po.CreatePOFromPRRequest{
		PRID:                 prID,
		SupplierID:           supplierID,
		SupplierCode:         req.SupplierCode,
		SupplierName:         req.SupplierName,
		DeliveryAddress:      req.DeliveryAddress,
		DeliveryTerms:        req.DeliveryTerms,
		PaymentTerms:         req.PaymentTerms,
		ExpectedDeliveryDate: req.ExpectedDeliveryDate,
		Notes:                req.Notes,
		CreatedBy:            createdBy,
	}

	result, err := h.createPOFromPR.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, dto.ToPOResponse(result))
}

// Get handles GET /api/v1/purchase-orders/:id
func (h *POHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PO ID"))
		return
	}

	result, err := h.getPO.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("PO not found"))
		return
	}

	response.Success(c, dto.ToPOResponse(result))
}

// List handles GET /api/v1/purchase-orders
func (h *POHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.POFilter{
		Status: c.Query("status"),
		Search: c.Query("search"),
		Page:   page,
		Limit:  limit,
	}

	if supplierID := c.Query("supplier_id"); supplierID != "" {
		if uid, err := uuid.Parse(supplierID); err == nil {
			filter.SupplierID = &uid
		}
	}

	results, total, err := h.listPOs.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMeta(c, dto.ToPOListResponse(results), &response.Meta{
		Page:       page,
		PageSize:   limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	})
}

// Confirm handles POST /api/v1/purchase-orders/:id/confirm
func (h *POHandler) Confirm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PO ID"))
		return
	}

	userID := c.GetHeader("X-User-ID")
	confirmedBy, _ := uuid.Parse(userID)
	if confirmedBy == uuid.Nil {
		confirmedBy = uuid.New()
	}

	result, err := h.confirmPO.Execute(c.Request.Context(), id, confirmedBy)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToPOResponse(result))
}

// Cancel handles POST /api/v1/purchase-orders/:id/cancel
func (h *POHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PO ID"))
		return
	}

	var req dto.CancelPORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("reason is required"))
		return
	}

	userID := c.GetHeader("X-User-ID")
	cancelledBy, _ := uuid.Parse(userID)
	if cancelledBy == uuid.Nil {
		cancelledBy = uuid.New()
	}

	result, err := h.cancelPO.Execute(c.Request.Context(), id, cancelledBy, req.Reason)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToPOResponse(result))
}

// Close handles POST /api/v1/purchase-orders/:id/close
func (h *POHandler) Close(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PO ID"))
		return
	}

	result, err := h.closePO.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToPOResponse(result))
}

// GetReceipts handles GET /api/v1/purchase-orders/:id/receipts
func (h *POHandler) GetReceipts(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PO ID"))
		return
	}

	receipts, err := h.getPOReceipts.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.ToPOReceiptResponses(receipts))
}
