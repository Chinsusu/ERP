package handler

import (
	"strconv"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/usecase/quotation"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// QuotationHandler handles quotation HTTP requests
type QuotationHandler struct {
	createQuotation *quotation.CreateQuotationUseCase
	getQuotation    *quotation.GetQuotationUseCase
	listQuotations  *quotation.ListQuotationsUseCase
	sendQuotation   *quotation.SendQuotationUseCase
	convertToOrder  *quotation.ConvertToOrderUseCase
	quotationRepo   repository.QuotationRepository
}

// NewQuotationHandler creates a new quotation handler
func NewQuotationHandler(
	createQuotation *quotation.CreateQuotationUseCase,
	getQuotation *quotation.GetQuotationUseCase,
	listQuotations *quotation.ListQuotationsUseCase,
	sendQuotation *quotation.SendQuotationUseCase,
	convertToOrder *quotation.ConvertToOrderUseCase,
	quotationRepo repository.QuotationRepository,
) *QuotationHandler {
	return &QuotationHandler{
		createQuotation: createQuotation,
		getQuotation:    getQuotation,
		listQuotations:  listQuotations,
		sendQuotation:   sendQuotation,
		convertToOrder:  convertToOrder,
		quotationRepo:   quotationRepo,
	}
}

// QuotationItemRequest represents a line item request
type QuotationItemRequest struct {
	ProductID       uuid.UUID  `json:"product_id" binding:"required"`
	ProductCode     string     `json:"product_code"`
	ProductName     string     `json:"product_name"`
	Quantity        float64    `json:"quantity" binding:"required,gt=0"`
	UomID           *uuid.UUID `json:"uom_id"`
	UnitPrice       float64    `json:"unit_price" binding:"required,gte=0"`
	DiscountPercent float64    `json:"discount_percent"`
	TaxPercent      float64    `json:"tax_percent"`
	Notes           string     `json:"notes"`
}

// CreateQuotationRequest represents create quotation request
type CreateQuotationRequest struct {
	CustomerID         uuid.UUID              `json:"customer_id" binding:"required"`
	QuotationDate      string                 `json:"quotation_date" binding:"required"`
	ValidUntil         string                 `json:"valid_until" binding:"required"`
	DiscountPercent    float64                `json:"discount_percent"`
	DiscountAmount     float64                `json:"discount_amount"`
	TaxPercent         float64                `json:"tax_percent"`
	Notes              string                 `json:"notes"`
	TermsAndConditions string                 `json:"terms_and_conditions"`
	Items              []QuotationItemRequest `json:"items" binding:"required,dive"`
}

// CreateQuotation handles POST /quotations
func (h *QuotationHandler) CreateQuotation(c *gin.Context) {
	var req CreateQuotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	quotationDate, _ := time.Parse("2006-01-02", req.QuotationDate)
	validUntil, _ := time.Parse("2006-01-02", req.ValidUntil)

	items := make([]quotation.QuotationItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = quotation.QuotationItemInput{
			ProductID:       item.ProductID,
			ProductCode:     item.ProductCode,
			ProductName:     item.ProductName,
			Quantity:        item.Quantity,
			UomID:           item.UomID,
			UnitPrice:       item.UnitPrice,
			DiscountPercent: item.DiscountPercent,
			TaxPercent:      item.TaxPercent,
			Notes:           item.Notes,
		}
	}

	input := &quotation.CreateQuotationInput{
		CustomerID:         req.CustomerID,
		QuotationDate:      quotationDate,
		ValidUntil:         validUntil,
		DiscountPercent:    req.DiscountPercent,
		DiscountAmount:     req.DiscountAmount,
		TaxPercent:         req.TaxPercent,
		Notes:              req.Notes,
		TermsAndConditions: req.TermsAndConditions,
		Items:              items,
	}

	result, err := h.createQuotation.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, result)
}

// GetQuotation handles GET /quotations/:id
func (h *QuotationHandler) GetQuotation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid quotation ID"))
		return
	}

	result, err := h.getQuotation.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("quotation"))
		return
	}

	response.Success(c, result)
}

// ListQuotations handles GET /quotations
func (h *QuotationHandler) ListQuotations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.QuotationFilter{
		Status:   entity.QuotationStatus(c.Query("status")),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
	}

	if customerID := c.Query("customer_id"); customerID != "" {
		if id, err := uuid.Parse(customerID); err == nil {
			filter.CustomerID = &id
		}
	}

	results, total, err := h.listQuotations.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// UpdateQuotation handles PUT /quotations/:id
func (h *QuotationHandler) UpdateQuotation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid quotation ID"))
		return
	}

	quot, err := h.quotationRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("quotation"))
		return
	}

	var req struct {
		Notes              string `json:"notes"`
		TermsAndConditions string `json:"terms_and_conditions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	quot.Notes = req.Notes
	quot.TermsAndConditions = req.TermsAndConditions

	if err := h.quotationRepo.Update(c.Request.Context(), quot); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, quot)
}

// SendQuotation handles PATCH /quotations/:id/send
func (h *QuotationHandler) SendQuotation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid quotation ID"))
		return
	}

	result, err := h.sendQuotation.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// ConvertToOrderRequest represents convert to order request
type ConvertToOrderRequest struct {
	DeliveryDate    string `json:"delivery_date"`
	DeliveryAddress string `json:"delivery_address"`
	PaymentMethod   string `json:"payment_method"`
	Notes           string `json:"notes"`
}

// ConvertToOrder handles POST /quotations/:id/convert-to-order
func (h *QuotationHandler) ConvertToOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid quotation ID"))
		return
	}

	var req ConvertToOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	var deliveryDate *time.Time
	if req.DeliveryDate != "" {
		t, _ := time.Parse("2006-01-02", req.DeliveryDate)
		deliveryDate = &t
	}

	paymentMethod := entity.PaymentMethodBankTransfer
	if req.PaymentMethod != "" {
		paymentMethod = entity.PaymentMethod(req.PaymentMethod)
	}

	input := &quotation.ConvertToOrderInput{
		QuotationID:     id,
		DeliveryDate:    deliveryDate,
		DeliveryAddress: req.DeliveryAddress,
		PaymentMethod:   paymentMethod,
		Notes:           req.Notes,
	}

	result, err := h.convertToOrder.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, result)
}
