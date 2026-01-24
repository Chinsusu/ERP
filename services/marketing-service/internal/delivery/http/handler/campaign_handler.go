package handler

import (
	"strconv"
	"time"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	"github.com/erp-cosmetics/marketing-service/internal/usecase/campaign"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CampaignHandler handles campaign HTTP requests
type CampaignHandler struct {
	createCampaign *campaign.CreateCampaignUseCase
	getCampaign    *campaign.GetCampaignUseCase
	listCampaigns  *campaign.ListCampaignsUseCase
	launchCampaign *campaign.LaunchCampaignUseCase
	updateCampaign *campaign.UpdateCampaignUseCase
	collabRepo     repository.KOLCollaborationRepository
}

// NewCampaignHandler creates a new campaign handler
func NewCampaignHandler(
	createCampaign *campaign.CreateCampaignUseCase,
	getCampaign *campaign.GetCampaignUseCase,
	listCampaigns *campaign.ListCampaignsUseCase,
	launchCampaign *campaign.LaunchCampaignUseCase,
	updateCampaign *campaign.UpdateCampaignUseCase,
	collabRepo repository.KOLCollaborationRepository,
) *CampaignHandler {
	return &CampaignHandler{
		createCampaign: createCampaign,
		getCampaign:    getCampaign,
		listCampaigns:  listCampaigns,
		launchCampaign: launchCampaign,
		updateCampaign: updateCampaign,
		collabRepo:     collabRepo,
	}
}

// CreateCampaignRequest represents create campaign request
type CreateCampaignRequest struct {
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	CampaignType   string   `json:"campaign_type" binding:"required"`
	StartDate      string   `json:"start_date" binding:"required"`
	EndDate        string   `json:"end_date" binding:"required"`
	TargetAudience string   `json:"target_audience"`
	Channels       []string `json:"channels"`
	Budget         float64  `json:"budget"`
	Currency       string   `json:"currency"`
	Notes          string   `json:"notes"`
}

// CreateCampaign handles POST /campaigns
func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var req CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	input := &campaign.CreateCampaignInput{
		Name:           req.Name,
		Description:    req.Description,
		CampaignType:   entity.CampaignType(req.CampaignType),
		StartDate:      startDate,
		EndDate:        endDate,
		TargetAudience: req.TargetAudience,
		Budget:         req.Budget,
		Currency:       req.Currency,
		Notes:          req.Notes,
	}

	result, err := h.createCampaign.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, result)
}

// GetCampaign handles GET /campaigns/:id
func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid campaign ID"))
		return
	}

	result, err := h.getCampaign.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("campaign"))
		return
	}

	response.Success(c, result)
}

// ListCampaigns handles GET /campaigns
func (h *CampaignHandler) ListCampaigns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.CampaignFilter{
		Search:       c.Query("search"),
		CampaignType: entity.CampaignType(c.Query("type")),
		Status:       entity.CampaignStatus(c.Query("status")),
		DateFrom:     c.Query("date_from"),
		DateTo:       c.Query("date_to"),
		Page:         page,
		Limit:        limit,
	}

	results, total, err := h.listCampaigns.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// UpdateCampaign handles PUT /campaigns/:id
func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid campaign ID"))
		return
	}

	var req struct {
		Name           string  `json:"name"`
		Description    string  `json:"description"`
		TargetAudience string  `json:"target_audience"`
		Budget         float64 `json:"budget"`
		Notes          string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &campaign.UpdateCampaignInput{
		Name:           req.Name,
		Description:    req.Description,
		TargetAudience: req.TargetAudience,
		Budget:         req.Budget,
		Notes:          req.Notes,
	}

	result, err := h.updateCampaign.Execute(c.Request.Context(), id, input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, result)
}

// LaunchCampaign handles PATCH /campaigns/:id/launch
func (h *CampaignHandler) LaunchCampaign(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid campaign ID"))
		return
	}

	result, err := h.launchCampaign.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// GetCampaignCollaborations handles GET /campaigns/:id/collaborations
func (h *CampaignHandler) GetCampaignCollaborations(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid campaign ID"))
		return
	}

	filter := &repository.KOLCollaborationFilter{
		CampaignID: &id,
		Page:       1,
		Limit:      100,
	}

	collabs, _, err := h.collabRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, collabs)
}

// GetCampaignPerformance handles GET /campaigns/:id/performance
func (h *CampaignHandler) GetCampaignPerformance(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid campaign ID"))
		return
	}

	camp, err := h.getCampaign.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("campaign"))
		return
	}

	performance := gin.H{
		"campaign_id":        camp.ID,
		"campaign_code":      camp.CampaignCode,
		"name":               camp.Name,
		"budget":             camp.Budget,
		"spent":              camp.Spent,
		"budget_utilization": camp.GetBudgetUtilization(),
		"impressions":        camp.Impressions,
		"reach":              camp.Reach,
		"engagement":         camp.Engagement,
		"conversions":        camp.Conversions,
		"revenue_generated":  camp.RevenueGenerated,
		"roi":                camp.GetROI(),
	}

	response.Success(c, performance)
}
