package handler

import (
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	categoryUC "github.com/erp-cosmetics/master-data-service/internal/usecase/category"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	uc *categoryUC.UseCase
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(uc *categoryUC.UseCase) *CategoryHandler {
	return &CategoryHandler{uc: uc}
}

// Create creates a new category
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	var parentID *uuid.UUID
	if req.ParentID != "" {
		id, err := uuid.Parse(req.ParentID)
		if err != nil {
			response.Error(c, errors.BadRequest("Invalid parent_id"))
			return
		}
		parentID = &id
	}

	category, err := h.uc.Create(c.Request.Context(), &categoryUC.CreateRequest{
		Code:         req.Code,
		Name:         req.Name,
		NameEN:       req.NameEN,
		Description:  req.Description,
		CategoryType: entity.CategoryType(req.CategoryType),
		ParentID:     parentID,
		SortOrder:    req.SortOrder,
	})
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, dto.ToCategoryResponse(category))
}

// GetByID retrieves a category by ID
func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid category ID"))
		return
	}

	category, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Category"))
		return
	}

	response.Success(c, dto.ToCategoryResponse(category))
}

// Update updates an existing category
func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid category ID"))
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	category, err := h.uc.Update(c.Request.Context(), &categoryUC.UpdateRequest{
		ID:          id,
		Name:        req.Name,
		NameEN:      req.NameEN,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
	})
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToCategoryResponse(category))
}

// Delete soft deletes a category
func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid category ID"))
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, gin.H{"message": "Category deleted successfully"})
}

// List lists categories
func (h *CategoryHandler) List(c *gin.Context) {
	var query dto.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	categoryType := entity.CategoryType(c.Query("category_type"))
	
	var parentID *uuid.UUID
	if pid := c.Query("parent_id"); pid != "" {
		id, err := uuid.Parse(pid)
		if err == nil {
			parentID = &id
		}
	}

	filter := &repository.CategoryFilter{
		CategoryType: categoryType,
		ParentID:     parentID,
		Status:       query.Status,
		Search:       query.Search,
		Page:         query.Page,
		PageSize:     query.PageSize,
	}

	categories, total, err := h.uc.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.CategoryResponse
	for _, cat := range categories {
		resp = append(resp, dto.ToCategoryResponse(&cat))
	}

	meta := response.NewMeta(query.Page, query.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// GetTree retrieves categories as a hierarchical tree
func (h *CategoryHandler) GetTree(c *gin.Context) {
	categoryType := entity.CategoryType(c.Query("category_type"))

	categories, err := h.uc.GetTree(c.Request.Context(), categoryType)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.CategoryResponse
	for _, cat := range categories {
		resp = append(resp, dto.ToCategoryResponse(&cat))
	}

	response.Success(c, resp)
}
