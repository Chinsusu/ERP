package handler

import (
	"github.com/erp-cosmetics/master-data-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/master-data-service/internal/domain/entity"
	"github.com/erp-cosmetics/master-data-service/internal/domain/repository"
	productUC "github.com/erp-cosmetics/master-data-service/internal/usecase/product"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	uc *productUC.UseCase
}

// NewProductHandler creates a new product handler
func NewProductHandler(uc *productUC.UseCase) *ProductHandler {
	return &ProductHandler{uc: uc}
}

// Create creates a new product
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	baseUnitID, err := uuid.Parse(req.BaseUnitID)
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid base_unit_id"))
		return
	}

	categoryID, _ := dto.ParseUUID(req.CategoryID)

	product, err := h.uc.Create(c.Request.Context(), &productUC.CreateRequest{
		Code:                  req.Code,
		SKU:                   req.SKU,
		Barcode:               req.Barcode,
		Name:                  req.Name,
		NameEN:                req.NameEN,
		Description:           req.Description,
		CategoryID:            categoryID,
		ProductLine:           req.ProductLine,
		Brand:                 req.Brand,
		Volume:                req.Volume,
		VolumeUnit:            req.VolumeUnit,
		CosmeticLicenseNumber: req.CosmeticLicenseNumber,
		IngredientsSummary:    req.IngredientsSummary,
		TargetSkinType:        req.TargetSkinType,
		PackagingType:         req.PackagingType,
		StandardCost:          req.StandardCost,
		StandardPrice:         req.StandardPrice,
		Currency:              req.Currency,
		BaseUnitID:            baseUnitID,
		ShelfLifeMonths:       req.ShelfLifeMonths,
	})
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, dto.ToProductResponse(product))
}

// GetByID retrieves a product by ID
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid product ID"))
		return
	}

	product, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Product"))
		return
	}

	response.Success(c, dto.ToProductResponse(product))
}

// Update updates an existing product
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid product ID"))
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	product, err := h.uc.Update(c.Request.Context(), id, updates)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, dto.ToProductResponse(product))
}

// Delete soft deletes a product
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid product ID"))
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, gin.H{"message": "Product deleted successfully"})
}

// List lists products
func (h *ProductHandler) List(c *gin.Context) {
	var query dto.ProductListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	categoryID, _ := dto.ParseUUID(query.CategoryID)

	filter := &repository.ProductFilter{
		CategoryID:  categoryID,
		ProductLine: query.ProductLine,
		Brand:       query.Brand,
		Status:      query.Status,
		Search:      query.Search,
		Page:        query.Page,
		PageSize:    query.PageSize,
	}

	products, total, err := h.uc.List(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.ProductResponse
	for _, prod := range products {
		resp = append(resp, dto.ToProductResponse(&prod))
	}

	meta := response.NewMeta(query.Page, query.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// Search searches products
func (h *ProductHandler) Search(c *gin.Context) {
	searchQuery := c.Query("q")
	if searchQuery == "" {
		h.List(c)
		return
	}

	var listQuery dto.ProductListQuery
	if err := c.ShouldBindQuery(&listQuery); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	categoryID, _ := dto.ParseUUID(listQuery.CategoryID)

	filter := &repository.ProductFilter{
		CategoryID:  categoryID,
		ProductLine: listQuery.ProductLine,
		Brand:       listQuery.Brand,
		Status:      listQuery.Status,
		Page:        listQuery.Page,
		PageSize:    listQuery.PageSize,
	}

	products, total, err := h.uc.Search(c.Request.Context(), searchQuery, filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.ProductResponse
	for _, prod := range products {
		resp = append(resp, dto.ToProductResponse(&prod))
	}

	meta := response.NewMeta(listQuery.Page, listQuery.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// GetByCategory gets products by category
func (h *ProductHandler) GetByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("category_id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid category ID"))
		return
	}

	var query dto.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	products, total, err := h.uc.GetByCategory(c.Request.Context(), categoryID, query.Page, query.PageSize)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var resp []dto.ProductResponse
	for _, prod := range products {
		resp = append(resp, dto.ToProductResponse(&prod))
	}

	meta := response.NewMeta(query.Page, query.PageSize, total)
	response.SuccessWithMeta(c, resp, meta)
}

// AddImage adds an image to a product
func (h *ProductHandler) AddImage(c *gin.Context) {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid product ID"))
		return
	}

	var image entity.ProductImage
	if err := c.ShouldBindJSON(&image); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	if err := h.uc.AddImage(c.Request.Context(), productID, &image); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, image)
}
