package handler

import (
	"time"

	"github.com/erp-cosmetics/wms-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/usecase/grn"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GRNHandler handles GRN endpoints
type GRNHandler struct {
	createGRNUC   *grn.CreateGRNUseCase
	completeGRNUC *grn.CompleteGRNUseCase
	getGRNUC      *grn.GetGRNUseCase
	listGRNsUC    *grn.ListGRNsUseCase
}

// NewGRNHandler creates a new GRN handler
func NewGRNHandler(
	createGRNUC *grn.CreateGRNUseCase,
	completeGRNUC *grn.CompleteGRNUseCase,
	getGRNUC *grn.GetGRNUseCase,
	listGRNsUC *grn.ListGRNsUseCase,
) *GRNHandler {
	return &GRNHandler{
		createGRNUC:   createGRNUC,
		completeGRNUC: completeGRNUC,
		getGRNUC:      getGRNUC,
		listGRNsUC:    listGRNsUC,
	}
}

// CreateGRN handles POST /grn
func (h *GRNHandler) CreateGRN(c *gin.Context) {
	var req dto.CreateGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	grnDate, err := time.Parse("2006-01-02", req.GRNDate)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid GRN date format"))
		return
	}

	// Get user ID from context (for now use a placeholder)
	userID := uuid.New()

	input := &grn.CreateGRNInput{
		GRNDate:            grnDate,
		POID:               req.POID,
		PONumber:           req.PONumber,
		SupplierID:         req.SupplierID,
		WarehouseID:        req.WarehouseID,
		DeliveryNoteNumber: req.DeliveryNoteNumber,
		VehicleNumber:      req.VehicleNumber,
		Notes:              req.Notes,
		ReceivedBy:         userID,
		Items:              make([]grn.CreateGRNItemInput, len(req.Items)),
	}

	for i, item := range req.Items {
		expiryDate, err := time.Parse("2006-01-02", item.ExpiryDate)
		if err != nil {
			response.Error(c, errors.BadRequest("Invalid expiry date format"))
			return
		}

		var manufacturedDate *time.Time
		if item.ManufacturedDate != nil {
			d, _ := time.Parse("2006-01-02", *item.ManufacturedDate)
			manufacturedDate = &d
		}

		input.Items[i] = grn.CreateGRNItemInput{
			POLineItemID:      item.POLineItemID,
			MaterialID:        item.MaterialID,
			ExpectedQty:       item.ExpectedQty,
			ReceivedQty:       item.ReceivedQty,
			UnitID:            item.UnitID,
			SupplierLotNumber: item.SupplierLotNumber,
			ManufacturedDate:  manufacturedDate,
			ExpiryDate:        expiryDate,
			LocationID:        item.LocationID,
		}
	}

	result, err := h.createGRNUC.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, gin.H{
		"id":         result.ID,
		"grn_number": result.GRNNumber,
		"status":     result.Status,
	})
}

// GetGRN handles GET /grn/:id
func (h *GRNHandler) GetGRN(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid GRN ID"))
		return
	}

	result, err := h.getGRNUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("GRN"))
		return
	}

	resp := toGRNResponse(result)
	response.Success(c, resp)
}

// ListGRNs handles GET /grn
func (h *GRNHandler) ListGRNs(c *gin.Context) {
	filter := &repository.GRNFilter{
		Status:   c.Query("status"),
		QCStatus: c.Query("qc_status"),
		Search:   c.Query("search"),
		Page:     getPageParam(c),
		Limit:    getLimitParam(c),
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		id, _ := uuid.Parse(warehouseID)
		filter.WarehouseID = &id
	}
	if poID := c.Query("po_id"); poID != "" {
		id, _ := uuid.Parse(poID)
		filter.POID = &id
	}

	grns, total, err := h.listGRNsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.GRNResponse
	for _, g := range grns {
		items = append(items, toGRNResponse(g))
	}

	response.SuccessWithMeta(c, items, response.NewMeta(filter.Page, filter.Limit, int64(total)))
}

// CompleteGRN handles PATCH /grn/:id/complete
func (h *GRNHandler) CompleteGRN(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid GRN ID"))
		return
	}

	var req dto.CompleteGRNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &grn.CompleteGRNInput{
		GRNID:    id,
		QCStatus: entity.QCStatus(req.QCStatus),
		QCNotes:  req.QCNotes,
	}

	result, err := h.completeGRNUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err == entity.ErrAlreadyCompleted {
			response.Error(c, errors.BadRequest("GRN already completed"))
			return
		}
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, toGRNResponse(result))
}

func toGRNResponse(g *entity.GRN) dto.GRNResponse {
	resp := dto.GRNResponse{
		ID:                 g.ID,
		GRNNumber:          g.GRNNumber,
		GRNDate:            g.GRNDate.Format("2006-01-02"),
		PONumber:           g.PONumber,
		WarehouseID:        g.WarehouseID,
		DeliveryNoteNumber: g.DeliveryNoteNumber,
		VehicleNumber:      g.VehicleNumber,
		Status:             string(g.Status),
		QCStatus:           string(g.QCStatus),
		QCNotes:            g.QCNotes,
		Notes:              g.Notes,
		CompletedAt:        g.CompletedAt,
		CreatedAt:          g.CreatedAt,
	}

	if g.Warehouse != nil {
		resp.WarehouseName = g.Warehouse.Name
	}

	for _, item := range g.LineItems {
		lineItem := dto.GRNLineItemResponse{
			ID:                item.ID,
			LineNumber:        item.LineNumber,
			MaterialID:        item.MaterialID,
			ExpectedQty:       item.ExpectedQty,
			ReceivedQty:       item.ReceivedQty,
			AcceptedQty:       item.AcceptedQty,
			RejectedQty:       item.RejectedQty,
			SupplierLotNumber: item.SupplierLotNumber,
			ExpiryDate:        item.ExpiryDate.Format("2006-01-02"),
			QCStatus:          string(item.QCStatus),
		}

		if item.Lot != nil {
			lineItem.LotNumber = item.Lot.LotNumber
		}
		if item.Location != nil {
			lineItem.LocationCode = item.Location.Code
		}

		resp.LineItems = append(resp.LineItems, lineItem)
	}

	return resp
}
