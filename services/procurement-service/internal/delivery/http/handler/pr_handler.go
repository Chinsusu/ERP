package handler

import (
	"strconv"

	"github.com/erp-cosmetics/procurement-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/procurement-service/internal/domain/repository"
	"github.com/erp-cosmetics/procurement-service/internal/usecase/pr"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PRHandler handles PR-related requests
type PRHandler struct {
	createPR  *pr.CreatePRUseCase
	getPR     *pr.GetPRUseCase
	listPRs   *pr.ListPRsUseCase
	submitPR  *pr.SubmitPRUseCase
	approvePR *pr.ApprovePRUseCase
	rejectPR  *pr.RejectPRUseCase
}

// NewPRHandler creates a new PR handler
func NewPRHandler(
	createPR *pr.CreatePRUseCase,
	getPR *pr.GetPRUseCase,
	listPRs *pr.ListPRsUseCase,
	submitPR *pr.SubmitPRUseCase,
	approvePR *pr.ApprovePRUseCase,
	rejectPR *pr.RejectPRUseCase,
) *PRHandler {
	return &PRHandler{
		createPR:  createPR,
		getPR:     getPR,
		listPRs:   listPRs,
		submitPR:  submitPR,
		approvePR: approvePR,
		rejectPR:  rejectPR,
	}
}

// Create handles POST /api/v1/purchase-requisitions
func (h *PRHandler) Create(c *gin.Context) {
	var req dto.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	// Get requester ID from header (set by API Gateway)
	userID := c.GetHeader("X-User-ID")
	requesterID, err := uuid.Parse(userID)
	if err != nil {
		requesterID = uuid.New() // For testing without gateway
	}

	// Convert to use case request
	ucReq := &pr.CreatePRRequest{
		RequiredDate:  req.RequiredDate,
		Priority:      req.Priority,
		Justification: req.Justification,
		Notes:         req.Notes,
		RequesterID:   requesterID,
	}

	for _, item := range req.Items {
		materialID, _ := uuid.Parse(item.MaterialID)
		ucReq.Items = append(ucReq.Items, pr.PRLineItemRequest{
			MaterialID:     materialID,
			MaterialCode:   item.MaterialCode,
			MaterialName:   item.MaterialName,
			Quantity:       item.Quantity,
			UOMCode:        item.UOMCode,
			UnitPrice:      item.UnitPrice,
			Specifications: item.Specifications,
		})
	}

	result, err := h.createPR.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto.ToPRResponse(result))
}

// Get handles GET /api/v1/purchase-requisitions/:id
func (h *PRHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PR ID"))
		return
	}

	result, err := h.getPR.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("PR not found"))
		return
	}

	response.Success(c, dto.ToPRResponse(result))
}

// List handles GET /api/v1/purchase-requisitions
func (h *PRHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.PRFilter{
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
		Search:   c.Query("search"),
		Page:     page,
		Limit:    limit,
	}

	if requesterID := c.Query("requester_id"); requesterID != "" {
		if uid, err := uuid.Parse(requesterID); err == nil {
			filter.RequesterID = &uid
		}
	}

	results, total, err := h.listPRs.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMeta(c, dto.ToPRListResponse(results), &response.Meta{
		Page:       page,
		PageSize:   limit,
		TotalItems: total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	})
}

// Submit handles POST /api/v1/purchase-requisitions/:id/submit
func (h *PRHandler) Submit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PR ID"))
		return
	}

	result, err := h.submitPR.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToPRResponse(result))
}

// Approve handles POST /api/v1/purchase-requisitions/:id/approve
func (h *PRHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PR ID"))
		return
	}

	var req dto.ApprovePRRequest
	c.ShouldBindJSON(&req)

	userID := c.GetHeader("X-User-ID")
	approverID, err := uuid.Parse(userID)
	if err != nil {
		approverID = uuid.New()
	}

	result, err := h.approvePR.Execute(c.Request.Context(), id, approverID, req.Notes)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToPRResponse(result))
}

// Reject handles POST /api/v1/purchase-requisitions/:id/reject
func (h *PRHandler) Reject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid PR ID"))
		return
	}

	var req dto.ApprovePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest("reason is required"))
		return
	}

	userID := c.GetHeader("X-User-ID")
	rejectedBy, err := uuid.Parse(userID)
	if err != nil {
		rejectedBy = uuid.New()
	}

	result, err := h.rejectPR.Execute(c.Request.Context(), id, rejectedBy, req.Reason)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToPRResponse(result))
}
