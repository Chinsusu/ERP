package handler

import (
	"strconv"
	"time"

	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/erp-cosmetics/supplier-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/erp-cosmetics/supplier-service/internal/usecase/supplier"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SupplierHandler handles supplier HTTP requests
type SupplierHandler struct {
	createUC    *supplier.CreateSupplierUseCase
	getUC       *supplier.GetSupplierUseCase
	listUC      *supplier.ListSuppliersUseCase
	updateUC    *supplier.UpdateSupplierUseCase
	approveUC   *supplier.ApproveSupplierUseCase
	blockUC     *supplier.BlockSupplierUseCase
	addressRepo repository.AddressRepository
	contactRepo repository.ContactRepository
}

// NewSupplierHandler creates a new SupplierHandler
func NewSupplierHandler(
	createUC *supplier.CreateSupplierUseCase,
	getUC *supplier.GetSupplierUseCase,
	listUC *supplier.ListSuppliersUseCase,
	updateUC *supplier.UpdateSupplierUseCase,
	approveUC *supplier.ApproveSupplierUseCase,
	blockUC *supplier.BlockSupplierUseCase,
	addressRepo repository.AddressRepository,
	contactRepo repository.ContactRepository,
) *SupplierHandler {
	return &SupplierHandler{
		createUC:    createUC,
		getUC:       getUC,
		listUC:      listUC,
		updateUC:    updateUC,
		approveUC:   approveUC,
		blockUC:     blockUC,
		addressRepo: addressRepo,
		contactRepo: contactRepo,
	}
}

// Create handles POST /api/v1/suppliers
func (h *SupplierHandler) Create(c *gin.Context) {
	var req dto.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	ucReq := &supplier.CreateSupplierRequest{
		Name:         req.Name,
		LegalName:    req.LegalName,
		TaxCode:      req.TaxCode,
		SupplierType: req.SupplierType,
		BusinessType: req.BusinessType,
		Email:        req.Email,
		Phone:        req.Phone,
		Website:      req.Website,
		PaymentTerms: req.PaymentTerms,
		Currency:     req.Currency,
		CreditLimit:  req.CreditLimit,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankBranch:   req.BankBranch,
		Notes:        req.Notes,
	}

	result, err := h.createUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, dto.ToSupplierResponse(result))
}

// Get handles GET /api/v1/suppliers/:id
func (h *SupplierHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	result, err := h.getUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Supplier"))
		return
	}

	response.Success(c, dto.ToSupplierResponse(result))
}

// List handles GET /api/v1/suppliers
func (h *SupplierHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.SupplierFilter{
		SupplierType: c.Query("supplier_type"),
		BusinessType: c.Query("business_type"),
		Status:       c.Query("status"),
		Search:       c.Query("search"),
		Page:         page,
		Limit:        limit,
	}

	if minRating := c.Query("min_rating"); minRating != "" {
		if rating, err := strconv.ParseFloat(minRating, 64); err == nil {
			filter.MinRating = &rating
		}
	}

	if hasGMP := c.Query("has_gmp"); hasGMP == "true" {
		gmp := true
		filter.HasGMP = &gmp
	}

	results, total, err := h.listUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.SupplierListResponse
	for _, s := range results {
		items = append(items, *dto.ToSupplierListResponse(s))
	}

	response.SuccessWithMeta(c, items, response.NewMeta(page, limit, total))
}

// Update handles PUT /api/v1/suppliers/:id
func (h *SupplierHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	ucReq := &supplier.UpdateSupplierRequest{
		Name:         req.Name,
		LegalName:    req.LegalName,
		TaxCode:      req.TaxCode,
		Email:        req.Email,
		Phone:        req.Phone,
		Fax:          req.Fax,
		Website:      req.Website,
		PaymentTerms: req.PaymentTerms,
		Currency:     req.Currency,
		CreditLimit:  req.CreditLimit,
		BankName:     req.BankName,
		BankAccount:  req.BankAccount,
		BankBranch:   req.BankBranch,
		Notes:        req.Notes,
	}

	result, err := h.updateUC.Execute(c.Request.Context(), id, ucReq)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, dto.ToSupplierResponse(result))
}

// Approve handles PATCH /api/v1/suppliers/:id/approve
func (h *SupplierHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.ApproveRequest
	c.ShouldBindJSON(&req)

	// Get user ID from header (injected by API Gateway)
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		userID = uuid.New() // Fallback for testing
	}

	result, err := h.approveUC.Execute(c.Request.Context(), id, userID, req.Notes)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, dto.ToSupplierResponse(result))
}

// Block handles PATCH /api/v1/suppliers/:id/block
func (h *SupplierHandler) Block(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.BlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		userID = uuid.New()
	}

	result, err := h.blockUC.Execute(c.Request.Context(), id, userID, req.Reason)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, dto.ToSupplierResponse(result))
}

// ListAddresses handles GET /api/v1/suppliers/:id/addresses
func (h *SupplierHandler) ListAddresses(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	addresses, err := h.addressRepo.GetBySupplierID(c.Request.Context(), supplierID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.AddressResponse
	for _, a := range addresses {
		items = append(items, *dto.ToAddressResponse(a))
	}

	response.Success(c, items)
}

// CreateAddress handles POST /api/v1/suppliers/:id/addresses
func (h *SupplierHandler) CreateAddress(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	country := req.Country
	if country == "" {
		country = "Vietnam"
	}

	address := &entity.Address{
		ID:           uuid.New(),
		SupplierID:   supplierID,
		AddressType:  entity.AddressType(req.AddressType),
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		Ward:         req.Ward,
		District:     req.District,
		City:         req.City,
		Province:     req.Province,
		Country:      country,
		PostalCode:   req.PostalCode,
		IsPrimary:    req.IsPrimary,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.addressRepo.Create(c.Request.Context(), address); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, dto.ToAddressResponse(address))
}

// ListContacts handles GET /api/v1/suppliers/:id/contacts
func (h *SupplierHandler) ListContacts(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	contacts, err := h.contactRepo.GetBySupplierID(c.Request.Context(), supplierID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	var items []dto.ContactResponse
	for _, ct := range contacts {
		items = append(items, *dto.ToContactResponse(ct))
	}

	response.Success(c, items)
}

// CreateContact handles POST /api/v1/suppliers/:id/contacts
func (h *SupplierHandler) CreateContact(c *gin.Context) {
	supplierID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid supplier ID"))
		return
	}

	var req dto.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	contact := &entity.Contact{
		ID:          uuid.New(),
		SupplierID:  supplierID,
		ContactType: entity.ContactType(req.ContactType),
		FullName:    req.FullName,
		Position:    req.Position,
		Department:  req.Department,
		Email:       req.Email,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		IsPrimary:   req.IsPrimary,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.contactRepo.Create(c.Request.Context(), contact); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, dto.ToContactResponse(contact))
}
