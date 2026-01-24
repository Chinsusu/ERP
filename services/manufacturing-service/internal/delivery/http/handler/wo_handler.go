package handler

import (
	"time"

	"github.com/erp-cosmetics/manufacturing-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/workorder"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WOHandler handles work order-related requests
type WOHandler struct {
	createWOUC   *workorder.CreateWOUseCase
	getWOUC      *workorder.GetWOUseCase
	listWOsUC    *workorder.ListWOsUseCase
	releaseWOUC  *workorder.ReleaseWOUseCase
	startWOUC    *workorder.StartWOUseCase
	completeWOUC *workorder.CompleteWOUseCase
}

// NewWOHandler creates a new WOHandler
func NewWOHandler(
	createWOUC *workorder.CreateWOUseCase,
	getWOUC *workorder.GetWOUseCase,
	listWOsUC *workorder.ListWOsUseCase,
	releaseWOUC *workorder.ReleaseWOUseCase,
	startWOUC *workorder.StartWOUseCase,
	completeWOUC *workorder.CompleteWOUseCase,
) *WOHandler {
	return &WOHandler{
		createWOUC:   createWOUC,
		getWOUC:      getWOUC,
		listWOsUC:    listWOsUC,
		releaseWOUC:  releaseWOUC,
		startWOUC:    startWOUC,
		completeWOUC: completeWOUC,
	}
}

// CreateWO creates a new work order
func (h *WOHandler) CreateWO(c *gin.Context) {
	var req dto.CreateWORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	priority := entity.WOPriorityNormal
	if req.Priority != "" {
		priority = entity.WOPriority(req.Priority)
	}

	input := workorder.CreateWOInput{
		ProductID:        req.ProductID,
		BOMID:            req.BOMID,
		PlannedQuantity:  req.PlannedQuantity,
		UOMID:            req.UOMID,
		PlannedStartDate: req.PlannedStartDate,
		PlannedEndDate:   req.PlannedEndDate,
		BatchNumber:      req.BatchNumber,
		SalesOrderID:     req.SalesOrderID,
		ProductionLine:   req.ProductionLine,
		Shift:            req.Shift,
		Priority:         priority,
		Notes:            req.Notes,
		CreatedBy:        userID,
	}

	result, err := h.createWOUC.Execute(c.Request.Context(), input)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	created(c, toWOResponse(result))
}

// GetWO gets a work order by ID
func (h *WOHandler) GetWO(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid work order ID")
		return
	}

	result, err := h.getWOUC.Execute(c.Request.Context(), id)
	if err != nil {
		notFound(c, "Work order not found")
		return
	}

	success(c, toWOResponse(result))
}

// ListWOs lists work orders
func (h *WOHandler) ListWOs(c *gin.Context) {
	filter := repository.WOFilter{
		Page:     getPageFromQuery(c),
		PageSize: getPageSizeFromQuery(c),
		Search:   c.Query("search"),
	}

	if status := c.Query("status"); status != "" {
		s := entity.WOStatus(status)
		filter.Status = &s
	}
	if priority := c.Query("priority"); priority != "" {
		p := entity.WOPriority(priority)
		filter.Priority = &p
	}

	wos, total, err := h.listWOsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		internalError(c, err.Error())
		return
	}

	var items []dto.WOResponse
	for _, wo := range wos {
		items = append(items, toWOResponse(wo))
	}

	successWithMeta(c, items, newMeta(filter.Page, filter.PageSize, total))
}

// ReleaseWO releases a work order
func (h *WOHandler) ReleaseWO(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid work order ID")
		return
	}

	userID := getUserIDFromContext(c)

	result, err := h.releaseWOUC.Execute(c.Request.Context(), id, userID)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	success(c, toWOResponse(result))
}

// StartWO starts a work order
func (h *WOHandler) StartWO(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid work order ID")
		return
	}

	var req dto.StartWORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	result, err := h.startWOUC.Execute(c.Request.Context(), id, req.SupervisorID)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	success(c, toWOResponse(result))
}

// CompleteWO completes a work order
func (h *WOHandler) CompleteWO(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		badRequest(c, "Invalid work order ID")
		return
	}

	var req dto.CompleteWORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c, err.Error())
		return
	}

	userID := getUserIDFromContext(c)

	input := workorder.CompleteWOInput{
		WOID:             id,
		ActualQuantity:   req.ActualQuantity,
		GoodQuantity:     req.GoodQuantity,
		RejectedQuantity: req.RejectedQuantity,
		Notes:            req.Notes,
		UpdatedBy:        userID,
	}

	result, err := h.completeWOUC.Execute(c.Request.Context(), input)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	success(c, toWOResponse(result))
}

func toWOResponse(wo *entity.WorkOrder) dto.WOResponse {
	return dto.WOResponse{
		ID:               wo.ID,
		WONumber:         wo.WONumber,
		WODate:           wo.WODate,
		ProductID:        wo.ProductID,
		BOMID:            wo.BOMID,
		Status:           string(wo.Status),
		Priority:         string(wo.Priority),
		PlannedQuantity:  wo.PlannedQuantity,
		ActualQuantity:   wo.ActualQuantity,
		GoodQuantity:     wo.GoodQuantity,
		RejectedQuantity: wo.RejectedQuantity,
		YieldPercentage:  wo.YieldPercentage,
		BatchNumber:      wo.BatchNumber,
		PlannedStartDate: wo.PlannedStartDate,
		PlannedEndDate:   wo.PlannedEndDate,
		ActualStartDate:  wo.ActualStartDate,
		ActualEndDate:    wo.ActualEndDate,
		ProductionLine:   wo.ProductionLine,
		Notes:            wo.Notes,
		CreatedAt:        wo.CreatedAt,
	}
}

// Dummy time for compilation
var _ = time.Now
