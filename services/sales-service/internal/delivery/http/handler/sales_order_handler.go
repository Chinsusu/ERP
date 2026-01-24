package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	salesorder "github.com/erp-cosmetics/sales-service/internal/usecase/sales_order"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SalesOrderHandler handles sales order HTTP requests
type SalesOrderHandler struct {
	createOrder  *salesorder.CreateOrderUseCase
	getOrder     *salesorder.GetOrderUseCase
	listOrders   *salesorder.ListOrdersUseCase
	confirmOrder *salesorder.ConfirmOrderUseCase
	cancelOrder  *salesorder.CancelOrderUseCase
	shipOrder    *salesorder.ShipOrderUseCase
	deliverOrder *salesorder.DeliverOrderUseCase
	orderRepo    repository.SalesOrderRepository
}

// NewSalesOrderHandler creates a new sales order handler
func NewSalesOrderHandler(
	createOrder *salesorder.CreateOrderUseCase,
	getOrder *salesorder.GetOrderUseCase,
	listOrders *salesorder.ListOrdersUseCase,
	confirmOrder *salesorder.ConfirmOrderUseCase,
	cancelOrder *salesorder.CancelOrderUseCase,
	shipOrder *salesorder.ShipOrderUseCase,
	deliverOrder *salesorder.DeliverOrderUseCase,
	orderRepo repository.SalesOrderRepository,
) *SalesOrderHandler {
	return &SalesOrderHandler{
		createOrder:  createOrder,
		getOrder:     getOrder,
		listOrders:   listOrders,
		confirmOrder: confirmOrder,
		cancelOrder:  cancelOrder,
		shipOrder:    shipOrder,
		deliverOrder: deliverOrder,
		orderRepo:    orderRepo,
	}
}

// OrderItemRequest represents a line item request
type OrderItemRequest struct {
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

// CreateOrderRequest represents create sales order request
type CreateOrderRequest struct {
	CustomerID      uuid.UUID          `json:"customer_id" binding:"required"`
	QuotationID     *uuid.UUID         `json:"quotation_id"`
	SODate          string             `json:"so_date"`
	DeliveryDate    string             `json:"delivery_date"`
	DeliveryAddress string             `json:"delivery_address"`
	BillingAddress  string             `json:"billing_address"`
	DiscountPercent float64            `json:"discount_percent"`
	TaxPercent      float64            `json:"tax_percent"`
	PaymentMethod   string             `json:"payment_method"`
	Notes           string             `json:"notes"`
	Items           []OrderItemRequest `json:"items" binding:"required,dive"`
}

// CreateOrder handles POST /sales-orders
func (h *SalesOrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	soDate := time.Now()
	if req.SODate != "" {
		soDate, _ = time.Parse("2006-01-02", req.SODate)
	}

	var deliveryDate *time.Time
	if req.DeliveryDate != "" {
		t, _ := time.Parse("2006-01-02", req.DeliveryDate)
		deliveryDate = &t
	}

	items := make([]salesorder.OrderItemInput, len(req.Items))
	for i, item := range req.Items {
		items[i] = salesorder.OrderItemInput{
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

	paymentMethod := entity.PaymentMethodBankTransfer
	if req.PaymentMethod != "" {
		paymentMethod = entity.PaymentMethod(req.PaymentMethod)
	}

	input := &salesorder.CreateOrderInput{
		CustomerID:      req.CustomerID,
		QuotationID:     req.QuotationID,
		SODate:          soDate,
		DeliveryDate:    deliveryDate,
		DeliveryAddress: req.DeliveryAddress,
		BillingAddress:  req.BillingAddress,
		DiscountPercent: req.DiscountPercent,
		TaxPercent:      req.TaxPercent,
		PaymentMethod:   paymentMethod,
		Notes:           req.Notes,
		Items:           items,
	}

	result, err := h.createOrder.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Created(c, result)
}

// GetOrder handles GET /sales-orders/:id
func (h *SalesOrderHandler) GetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid order ID"))
		return
	}

	result, err := h.getOrder.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("order"))
		return
	}

	response.Success(c, result)
}

// ListOrders handles GET /sales-orders
func (h *SalesOrderHandler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.SalesOrderFilter{
		Status:        entity.SOStatus(c.Query("status")),
		PaymentStatus: entity.PaymentStatus(c.Query("payment_status")),
		DateFrom:      c.Query("date_from"),
		DateTo:        c.Query("date_to"),
		Page:          page,
		Limit:         limit,
	}

	if customerID := c.Query("customer_id"); customerID != "" {
		if id, err := uuid.Parse(customerID); err == nil {
			filter.CustomerID = &id
		}
	}

	results, total, err := h.listOrders.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// UpdateOrder handles PUT /sales-orders/:id
func (h *SalesOrderHandler) UpdateOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid order ID"))
		return
	}

	order, err := h.orderRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("order"))
		return
	}

	var req struct {
		DeliveryAddress string `json:"delivery_address"`
		Notes           string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	order.DeliveryAddress = req.DeliveryAddress
	order.Notes = req.Notes

	if err := h.orderRepo.Update(c.Request.Context(), order); err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, order)
}

// ConfirmOrder handles PATCH /sales-orders/:id/confirm
func (h *SalesOrderHandler) ConfirmOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid order ID"))
		return
	}

	// TODO: Get user ID from JWT token
	userID := uuid.New()

	result, err := h.confirmOrder.Execute(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// CancelOrderRequest represents cancel order request
type CancelOrderRequest struct {
	Reason string `json:"reason"`
}

// CancelOrder handles PATCH /sales-orders/:id/cancel
func (h *SalesOrderHandler) CancelOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid order ID"))
		return
	}

	var req CancelOrderRequest
	c.ShouldBindJSON(&req)

	// TODO: Get user ID from JWT token
	userID := uuid.New()

	result, err := h.cancelOrder.Execute(c.Request.Context(), id, userID, req.Reason)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// ShipOrder handles PATCH /sales-orders/:id/ship
func (h *SalesOrderHandler) ShipOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid order ID"))
		return
	}

	result, err := h.shipOrder.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// DeliverOrder handles PATCH /sales-orders/:id/deliver
func (h *SalesOrderHandler) DeliverOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid order ID"))
		return
	}

	result, err := h.deliverOrder.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
