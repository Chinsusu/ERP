package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/usecase/customer"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CustomerHandler handles customer HTTP requests
type CustomerHandler struct {
	createCustomer  *customer.CreateCustomerUseCase
	getCustomer     *customer.GetCustomerUseCase
	listCustomers   *customer.ListCustomersUseCase
	updateCustomer  *customer.UpdateCustomerUseCase
	deleteCustomer  *customer.DeleteCustomerUseCase
	checkCredit     *customer.CheckCreditUseCase
	customerRepo    repository.CustomerRepository
	groupRepo       repository.CustomerGroupRepository
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(
	createCustomer *customer.CreateCustomerUseCase,
	getCustomer *customer.GetCustomerUseCase,
	listCustomers *customer.ListCustomersUseCase,
	updateCustomer *customer.UpdateCustomerUseCase,
	deleteCustomer *customer.DeleteCustomerUseCase,
	checkCredit *customer.CheckCreditUseCase,
	customerRepo repository.CustomerRepository,
	groupRepo repository.CustomerGroupRepository,
) *CustomerHandler {
	return &CustomerHandler{
		createCustomer:  createCustomer,
		getCustomer:     getCustomer,
		listCustomers:   listCustomers,
		updateCustomer:  updateCustomer,
		deleteCustomer:  deleteCustomer,
		checkCredit:     checkCredit,
		customerRepo:    customerRepo,
		groupRepo:       groupRepo,
	}
}

// CreateCustomerRequest represents create customer request
type CreateCustomerRequest struct {
	Name            string     `json:"name" binding:"required"`
	TaxCode         string     `json:"tax_code"`
	CustomerType    string     `json:"customer_type"`
	CustomerGroupID *uuid.UUID `json:"customer_group_id"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Website         string     `json:"website"`
	PaymentTerms    string     `json:"payment_terms"`
	CreditLimit     float64    `json:"credit_limit"`
	Currency        string     `json:"currency"`
	Notes           string     `json:"notes"`
}

// CreateCustomer handles POST /customers
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &customer.CreateCustomerInput{
		Name:            req.Name,
		TaxCode:         req.TaxCode,
		CustomerType:    entity.CustomerType(req.CustomerType),
		CustomerGroupID: req.CustomerGroupID,
		Email:           req.Email,
		Phone:           req.Phone,
		Website:         req.Website,
		PaymentTerms:    req.PaymentTerms,
		CreditLimit:     req.CreditLimit,
		Currency:        req.Currency,
		Notes:           req.Notes,
	}

	result, err := h.createCustomer.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, result)
}

// GetCustomer handles GET /customers/:id
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	result, err := h.getCustomer.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("customer"))
		return
	}

	response.Success(c, result)
}

// ListCustomers handles GET /customers
func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.CustomerFilter{
		Search:       c.Query("search"),
		CustomerType: entity.CustomerType(c.Query("customer_type")),
		Status:       entity.CustomerStatus(c.Query("status")),
		Page:         page,
		Limit:        limit,
	}

	if groupID := c.Query("customer_group_id"); groupID != "" {
		if id, err := uuid.Parse(groupID); err == nil {
			filter.CustomerGroupID = &id
		}
	}

	results, total, err := h.listCustomers.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// UpdateCustomerRequest represents update customer request
type UpdateCustomerRequest struct {
	Name            string     `json:"name" binding:"required"`
	TaxCode         string     `json:"tax_code"`
	CustomerType    string     `json:"customer_type"`
	CustomerGroupID *uuid.UUID `json:"customer_group_id"`
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Website         string     `json:"website"`
	PaymentTerms    string     `json:"payment_terms"`
	CreditLimit     float64    `json:"credit_limit"`
	Currency        string     `json:"currency"`
	Status          string     `json:"status"`
	Notes           string     `json:"notes"`
}

// UpdateCustomer handles PUT /customers/:id
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &customer.UpdateCustomerInput{
		Name:            req.Name,
		TaxCode:         req.TaxCode,
		CustomerType:    entity.CustomerType(req.CustomerType),
		CustomerGroupID: req.CustomerGroupID,
		Email:           req.Email,
		Phone:           req.Phone,
		Website:         req.Website,
		PaymentTerms:    req.PaymentTerms,
		CreditLimit:     req.CreditLimit,
		Currency:        req.Currency,
		Status:          entity.CustomerStatus(req.Status),
		Notes:           req.Notes,
	}

	result, err := h.updateCustomer.Execute(c.Request.Context(), id, input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, result)
}

// DeleteCustomer handles DELETE /customers/:id
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	if err := h.deleteCustomer.Execute(c.Request.Context(), id); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "customer deleted"})
}

// ListGroups handles GET /customer-groups
func (h *CustomerHandler) ListGroups(c *gin.Context) {
	groups, err := h.groupRepo.List(c.Request.Context(), true)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}
	response.Success(c, groups)
}

// GetGroup handles GET /customer-groups/:id
func (h *CustomerHandler) GetGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid group ID"))
		return
	}

	group, err := h.groupRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("group"))
		return
	}

	response.Success(c, group)
}

// GetAddresses handles GET /customers/:id/addresses
func (h *CustomerHandler) GetAddresses(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	addresses, err := h.customerRepo.GetAddresses(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, addresses)
}

// CreateAddressRequest represents create address request
type CreateAddressRequest struct {
	AddressType  string `json:"address_type" binding:"required"`
	AddressLine1 string `json:"address_line1" binding:"required"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	IsDefault    bool   `json:"is_default"`
}

// CreateAddress handles POST /customers/:id/addresses
func (h *CustomerHandler) CreateAddress(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	var req CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	address := &entity.CustomerAddress{
		CustomerID:   customerID,
		AddressType:  entity.AddressType(req.AddressType),
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		PostalCode:   req.PostalCode,
		Country:      req.Country,
		IsDefault:    req.IsDefault,
	}

	if err := h.customerRepo.CreateAddress(c.Request.Context(), address); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, address)
}

// GetContacts handles GET /customers/:id/contacts
func (h *CustomerHandler) GetContacts(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	contacts, err := h.customerRepo.GetContacts(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, contacts)
}

// CreateContactRequest represents create contact request
type CreateContactRequest struct {
	ContactName string `json:"contact_name" binding:"required"`
	Position    string `json:"position"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Mobile      string `json:"mobile"`
	IsPrimary   bool   `json:"is_primary"`
	Notes       string `json:"notes"`
}

// CreateContact handles POST /customers/:id/contacts
func (h *CustomerHandler) CreateContact(c *gin.Context) {
	customerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	var req CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	contact := &entity.CustomerContact{
		CustomerID:  customerID,
		ContactName: req.ContactName,
		Position:    req.Position,
		Email:       req.Email,
		Phone:       req.Phone,
		Mobile:      req.Mobile,
		IsPrimary:   req.IsPrimary,
		Notes:       req.Notes,
	}

	if err := h.customerRepo.CreateContact(c.Request.Context(), contact); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, contact)
}

// CheckCredit handles GET /customers/:id/credit-check
func (h *CustomerHandler) CheckCredit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid customer ID"))
		return
	}

	amount, _ := strconv.ParseFloat(c.DefaultQuery("amount", "0"), 64)

	result, err := h.checkCredit.Execute(c.Request.Context(), id, amount)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, result)
}

// Unused import placeholder
var _ = fmt.Sprint
