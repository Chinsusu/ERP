package handler

import (
	"strconv"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	"github.com/erp-cosmetics/marketing-service/internal/usecase/kol"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// KOLHandler handles KOL HTTP requests
type KOLHandler struct {
	createKOL  *kol.CreateKOLUseCase
	getKOL     *kol.GetKOLUseCase
	listKOLs   *kol.ListKOLsUseCase
	updateKOL  *kol.UpdateKOLUseCase
	deleteKOL  *kol.DeleteKOLUseCase
	tierRepo   repository.KOLTierRepository
	postRepo   repository.KOLPostRepository
}

// NewKOLHandler creates a new KOL handler
func NewKOLHandler(
	createKOL *kol.CreateKOLUseCase,
	getKOL *kol.GetKOLUseCase,
	listKOLs *kol.ListKOLsUseCase,
	updateKOL *kol.UpdateKOLUseCase,
	deleteKOL *kol.DeleteKOLUseCase,
	tierRepo repository.KOLTierRepository,
	postRepo repository.KOLPostRepository,
) *KOLHandler {
	return &KOLHandler{
		createKOL: createKOL,
		getKOL:    getKOL,
		listKOLs:  listKOLs,
		updateKOL: updateKOL,
		deleteKOL: deleteKOL,
		tierRepo:  tierRepo,
		postRepo:  postRepo,
	}
}

// CreateKOLRequest represents create KOL request
type CreateKOLRequest struct {
	Name               string     `json:"name" binding:"required"`
	Email              string     `json:"email"`
	Phone              string     `json:"phone"`
	TierID             *uuid.UUID `json:"tier_id"`
	Category           string     `json:"category"`
	InstagramHandle    string     `json:"instagram_handle"`
	InstagramFollowers int        `json:"instagram_followers"`
	YouTubeChannel     string     `json:"youtube_channel"`
	YouTubeSubscribers int        `json:"youtube_subscribers"`
	TikTokHandle       string     `json:"tiktok_handle"`
	TikTokFollowers    int        `json:"tiktok_followers"`
	AvgEngagementRate  float64    `json:"avg_engagement_rate"`
	Niche              string     `json:"niche"`
	CollaborationRate  float64    `json:"collaboration_rate"`
	Currency           string     `json:"currency"`
	AddressLine1       string     `json:"address_line1"`
	City               string     `json:"city"`
	Notes              string     `json:"notes"`
}

// CreateKOL handles POST /kols
func (h *KOLHandler) CreateKOL(c *gin.Context) {
	var req CreateKOLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &kol.CreateKOLInput{
		Name:               req.Name,
		Email:              req.Email,
		Phone:              req.Phone,
		TierID:             req.TierID,
		Category:           entity.KOLCategory(req.Category),
		InstagramHandle:    req.InstagramHandle,
		InstagramFollowers: req.InstagramFollowers,
		YouTubeChannel:     req.YouTubeChannel,
		YouTubeSubscribers: req.YouTubeSubscribers,
		TikTokHandle:       req.TikTokHandle,
		TikTokFollowers:    req.TikTokFollowers,
		AvgEngagementRate:  req.AvgEngagementRate,
		Niche:              req.Niche,
		CollaborationRate:  req.CollaborationRate,
		Currency:           req.Currency,
		AddressLine1:       req.AddressLine1,
		City:               req.City,
		Notes:              req.Notes,
	}

	result, err := h.createKOL.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, result)
}

// GetKOL handles GET /kols/:id
func (h *KOLHandler) GetKOL(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid KOL ID"))
		return
	}

	result, err := h.getKOL.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("KOL"))
		return
	}

	response.Success(c, result)
}

// ListKOLs handles GET /kols
func (h *KOLHandler) ListKOLs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.KOLFilter{
		Search:   c.Query("search"),
		Category: entity.KOLCategory(c.Query("category")),
		Niche:    c.Query("niche"),
		Status:   entity.KOLStatus(c.Query("status")),
		Page:     page,
		Limit:    limit,
	}

	if tierID := c.Query("tier_id"); tierID != "" {
		if id, err := uuid.Parse(tierID); err == nil {
			filter.TierID = &id
		}
	}

	results, total, err := h.listKOLs.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// UpdateKOL handles PUT /kols/:id
func (h *KOLHandler) UpdateKOL(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid KOL ID"))
		return
	}

	var req CreateKOLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &kol.UpdateKOLInput{
		Name:               req.Name,
		Email:              req.Email,
		Phone:              req.Phone,
		TierID:             req.TierID,
		Category:           entity.KOLCategory(req.Category),
		InstagramHandle:    req.InstagramHandle,
		InstagramFollowers: req.InstagramFollowers,
		YouTubeChannel:     req.YouTubeChannel,
		YouTubeSubscribers: req.YouTubeSubscribers,
		TikTokHandle:       req.TikTokHandle,
		TikTokFollowers:    req.TikTokFollowers,
		AvgEngagementRate:  req.AvgEngagementRate,
		Niche:              req.Niche,
		CollaborationRate:  req.CollaborationRate,
		AddressLine1:       req.AddressLine1,
		City:               req.City,
		Notes:              req.Notes,
	}

	result, err := h.updateKOL.Execute(c.Request.Context(), id, input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, result)
}

// DeleteKOL handles DELETE /kols/:id
func (h *KOLHandler) DeleteKOL(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid KOL ID"))
		return
	}

	if err := h.deleteKOL.Execute(c.Request.Context(), id); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.NoContent(c)
}

// ListTiers handles GET /kol-tiers
func (h *KOLHandler) ListTiers(c *gin.Context) {
	tiers, err := h.tierRepo.List(c.Request.Context(), true)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}
	response.Success(c, tiers)
}

// GetKOLPosts handles GET /kols/:id/posts
func (h *KOLHandler) GetKOLPosts(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid KOL ID"))
		return
	}

	posts, err := h.postRepo.GetByKOL(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, posts)
}
