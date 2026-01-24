package handler

import (
	"time"

	"github.com/erp-cosmetics/wms-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/usecase/issue"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GoodsIssueHandler handles goods issue endpoints
type GoodsIssueHandler struct {
	createIssueUC *issue.CreateGoodsIssueUseCase
	getIssueUC    *issue.GetGoodsIssueUseCase
	listIssuesUC  *issue.ListGoodsIssuesUseCase
}

// NewGoodsIssueHandler creates a new goods issue handler
func NewGoodsIssueHandler(
	createIssueUC *issue.CreateGoodsIssueUseCase,
	getIssueUC *issue.GetGoodsIssueUseCase,
	listIssuesUC *issue.ListGoodsIssuesUseCase,
) *GoodsIssueHandler {
	return &GoodsIssueHandler{
		createIssueUC: createIssueUC,
		getIssueUC:    getIssueUC,
		listIssuesUC:  listIssuesUC,
	}
}

// CreateGoodsIssue handles POST /goods-issue
func (h *GoodsIssueHandler) CreateGoodsIssue(c *gin.Context) {
	var req dto.IssueStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid issue date format"))
		return
	}

	// Get user ID (placeholder)
	userID := uuid.New()

	input := &issue.CreateGoodsIssueInput{
		IssueDate:       issueDate,
		IssueType:       entity.IssueType(req.IssueType),
		ReferenceID:     req.ReferenceID,
		ReferenceNumber: req.ReferenceNumber,
		WarehouseID:     req.WarehouseID,
		Notes:           req.Notes,
		IssuedBy:        userID,
		Items:           make([]issue.CreateGoodsIssueItemInput, len(req.Items)),
	}

	for i, item := range req.Items {
		input.Items[i] = issue.CreateGoodsIssueItemInput{
			MaterialID: item.MaterialID,
			Quantity:   item.Quantity,
			UnitID:     item.UnitID,
		}
	}

	result, err := h.createIssueUC.Execute(c.Request.Context(), input)
	if err != nil {
		if err == entity.ErrInsufficientStock {
			response.Error(c, errors.BadRequest("Insufficient stock available"))
			return
		}
		response.Error(c, errors.Internal(err))
		return
	}

	// Convert lots issued to response format
	lineItems := make([]dto.IssueStockResponse, len(result.LineItems))
	for i, li := range result.LineItems {
		lotsUsed := make([]dto.LotIssuedResponse, len(li.LotsUsed))
		for j, lot := range li.LotsUsed {
			lotsUsed[j] = dto.LotIssuedResponse{
				LotNumber:  lot.LotNumber,
				Quantity:   lot.Quantity,
				ExpiryDate: lot.ExpiryDate.Format("2006-01-02"),
			}
		}
		lineItems[i] = dto.IssueStockResponse{
			IssueNumber: result.IssueNumber,
			LotsIssued:  lotsUsed,
		}
	}

	response.Created(c, gin.H{
		"id":           result.ID,
		"issue_number": result.IssueNumber,
		"status":       result.Status,
		"line_items":   lineItems,
	})
}

// GetGoodsIssue handles GET /goods-issue/:id
func (h *GoodsIssueHandler) GetGoodsIssue(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid goods issue ID"))
		return
	}

	result, err := h.getIssueUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Goods Issue"))
		return
	}

	response.Success(c, result)
}

// ListGoodsIssues handles GET /goods-issue
func (h *GoodsIssueHandler) ListGoodsIssues(c *gin.Context) {
	filter := &repository.GoodsIssueFilter{
		IssueType: c.Query("issue_type"),
		Status:    c.Query("status"),
		Search:    c.Query("search"),
		Page:      getPageParam(c),
		Limit:     getLimitParam(c),
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		id, _ := uuid.Parse(warehouseID)
		filter.WarehouseID = &id
	}

	issues, total, err := h.listIssuesUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.SuccessWithMeta(c, issues, response.NewMeta(filter.Page, filter.Limit, int64(total)))
}
