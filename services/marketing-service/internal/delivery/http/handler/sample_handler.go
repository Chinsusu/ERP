package handler

import (
	"strconv"
	"time"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	sampleuc "github.com/erp-cosmetics/marketing-service/internal/usecase/sample"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SampleHandler handles sample request HTTP requests
type SampleHandler struct {
	createSampleRequest  *sampleuc.CreateSampleRequestUseCase
	getSampleRequest     *sampleuc.GetSampleRequestUseCase
	listSampleRequests   *sampleuc.ListSampleRequestsUseCase
	approveSampleRequest *sampleuc.ApproveSampleRequestUseCase
	rejectSampleRequest  *sampleuc.RejectSampleRequestUseCase
	shipSample           *sampleuc.ShipSampleUseCase
	shipmentRepo         repository.SampleShipmentRepository
}

// NewSampleHandler creates a new sample handler
func NewSampleHandler(
	createSampleRequest *sampleuc.CreateSampleRequestUseCase,
	getSampleRequest *sampleuc.GetSampleRequestUseCase,
	listSampleRequests *sampleuc.ListSampleRequestsUseCase,
	approveSampleRequest *sampleuc.ApproveSampleRequestUseCase,
	rejectSampleRequest *sampleuc.RejectSampleRequestUseCase,
	shipSample *sampleuc.ShipSampleUseCase,
	shipmentRepo repository.SampleShipmentRepository,
) *SampleHandler {
	return &SampleHandler{
		createSampleRequest:  createSampleRequest,
		getSampleRequest:     getSampleRequest,
		listSampleRequests:   listSampleRequests,
		approveSampleRequest: approveSampleRequest,
		rejectSampleRequest:  rejectSampleRequest,
		shipSample:           shipSample,
		shipmentRepo:         shipmentRepo,
	}
}

// SampleItemRequest represents sample item in request
type SampleItemRequest struct {
	ProductID   uuid.UUID `json:"product_id" binding:"required"`
	ProductCode string    `json:"product_code"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity" binding:"required,gt=0"`
	UnitValue   float64   `json:"unit_value"`
	Notes       string    `json:"notes"`
}

// CreateSampleRequestRequest represents create sample request
type CreateSampleRequestRequest struct {
	KOLID            uuid.UUID           `json:"kol_id" binding:"required"`
	CampaignID       *uuid.UUID          `json:"campaign_id"`
	RequestReason    string              `json:"request_reason"`
	DeliveryAddress  string              `json:"delivery_address"`
	RecipientName    string              `json:"recipient_name"`
	RecipientPhone   string              `json:"recipient_phone"`
	ExpectedPostDate string              `json:"expected_post_date"`
	ExpectedReach    int                 `json:"expected_reach"`
	Items            []SampleItemRequest `json:"items" binding:"required,dive"`
	Notes            string              `json:"notes"`
}

// CreateSampleRequest handles POST /samples/requests
func (h *SampleHandler) CreateSampleRequest(c *gin.Context) {
	var req CreateSampleRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	var expectedPostDate *time.Time
	if req.ExpectedPostDate != "" {
		t, _ := time.Parse("2006-01-02", req.ExpectedPostDate)
		expectedPostDate = &t
	}

	items := make([]sampleuc.SampleItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = sampleuc.SampleItemInput{
			ProductID:   item.ProductID,
			ProductCode: item.ProductCode,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			UnitValue:   item.UnitValue,
			Notes:       item.Notes,
		}
	}

	input := &sampleuc.CreateSampleRequestInput{
		KOLID:            req.KOLID,
		CampaignID:       req.CampaignID,
		RequestReason:    req.RequestReason,
		DeliveryAddress:  req.DeliveryAddress,
		RecipientName:    req.RecipientName,
		RecipientPhone:   req.RecipientPhone,
		ExpectedPostDate: expectedPostDate,
		ExpectedReach:    req.ExpectedReach,
		Items:            items,
		Notes:            req.Notes,
	}

	result, err := h.createSampleRequest.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, result)
}

// GetSampleRequest handles GET /samples/requests/:id
func (h *SampleHandler) GetSampleRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid request ID"))
		return
	}

	result, err := h.getSampleRequest.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("sample request"))
		return
	}

	response.Success(c, result)
}

// ListSampleRequests handles GET /samples/requests
func (h *SampleHandler) ListSampleRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.SampleRequestFilter{
		Status:   entity.SampleRequestStatus(c.Query("status")),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
	}

	if kolID := c.Query("kol_id"); kolID != "" {
		if id, err := uuid.Parse(kolID); err == nil {
			filter.KOLID = &id
		}
	}
	if campaignID := c.Query("campaign_id"); campaignID != "" {
		if id, err := uuid.Parse(campaignID); err == nil {
			filter.CampaignID = &id
		}
	}

	results, total, err := h.listSampleRequests.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// ApproveRequest represents approve request body
type ApproveRequest struct {
	Approved bool   `json:"approved"`
	Notes    string `json:"notes"`
}

// ApproveSampleRequest handles PATCH /samples/requests/:id/approve
func (h *SampleHandler) ApproveSampleRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid request ID"))
		return
	}

	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	// TODO: Get user ID from JWT token
	userID := uuid.New()

	if req.Approved {
		result, err := h.approveSampleRequest.Execute(c.Request.Context(), id, userID)
		if err != nil {
			response.Error(c, errors.BadRequest(err.Error()))
			return
		}
		response.Success(c, result)
	} else {
		result, err := h.rejectSampleRequest.Execute(c.Request.Context(), id, req.Notes)
		if err != nil {
			response.Error(c, errors.BadRequest(err.Error()))
			return
		}
		response.Success(c, result)
	}
}

// ShipRequest represents ship request body
type ShipRequest struct {
	Courier           string `json:"courier"`
	TrackingNumber    string `json:"tracking_number"`
	EstimatedDelivery string `json:"estimated_delivery"`
}

// ShipSampleRequest handles PATCH /samples/requests/:id/ship
func (h *SampleHandler) ShipSampleRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid request ID"))
		return
	}

	var req ShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	var estimatedDelivery *time.Time
	if req.EstimatedDelivery != "" {
		t, _ := time.Parse("2006-01-02", req.EstimatedDelivery)
		estimatedDelivery = &t
	}

	input := &sampleuc.ShipSampleInput{
		RequestID:         id,
		Courier:           req.Courier,
		TrackingNumber:    req.TrackingNumber,
		EstimatedDelivery: estimatedDelivery,
	}

	result, err := h.shipSample.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// ListShipments handles GET /samples/shipments
func (h *SampleHandler) ListShipments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.SampleShipmentFilter{
		Status:  entity.ShipmentStatus(c.Query("status")),
		Courier: c.Query("courier"),
		Page:    page,
		Limit:   limit,
	}

	shipments, total, err := h.shipmentRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, shipments, meta)
}
