package handler

import (
	"strconv"

	"github.com/erp-cosmetics/sales-service/internal/domain/entity"
	"github.com/erp-cosmetics/sales-service/internal/domain/repository"
	"github.com/erp-cosmetics/sales-service/internal/usecase/shipment"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShipmentHandler handles shipment HTTP requests
type ShipmentHandler struct {
	createShipment  *shipment.CreateShipmentUseCase
	getShipment     *shipment.GetShipmentUseCase
	listShipments   *shipment.ListShipmentsUseCase
	shipShipment    *shipment.ShipShipmentUseCase
	deliverShipment *shipment.DeliverShipmentUseCase
}

// NewShipmentHandler creates a new shipment handler
func NewShipmentHandler(
	createShipment *shipment.CreateShipmentUseCase,
	getShipment *shipment.GetShipmentUseCase,
	listShipments *shipment.ListShipmentsUseCase,
	shipShipment *shipment.ShipShipmentUseCase,
	deliverShipment *shipment.DeliverShipmentUseCase,
) *ShipmentHandler {
	return &ShipmentHandler{
		createShipment:  createShipment,
		getShipment:     getShipment,
		listShipments:   listShipments,
		shipShipment:    shipShipment,
		deliverShipment: deliverShipment,
	}
}

// CreateShipmentRequest represents create shipment request
type CreateShipmentRequest struct {
	SalesOrderID    uuid.UUID `json:"sales_order_id" binding:"required"`
	Carrier         string    `json:"carrier"`
	TrackingNumber  string    `json:"tracking_number"`
	ShippingMethod  string    `json:"shipping_method"`
	ShippingCost    float64   `json:"shipping_cost"`
	RecipientName   string    `json:"recipient_name"`
	RecipientPhone  string    `json:"recipient_phone"`
	DeliveryAddress string    `json:"delivery_address"`
	Notes           string    `json:"notes"`
}

// CreateShipment handles POST /shipments
func (h *ShipmentHandler) CreateShipment(c *gin.Context) {
	var req CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	input := &shipment.CreateShipmentInput{
		SalesOrderID:    req.SalesOrderID,
		Carrier:         req.Carrier,
		TrackingNumber:  req.TrackingNumber,
		ShippingMethod:  req.ShippingMethod,
		ShippingCost:    req.ShippingCost,
		RecipientName:   req.RecipientName,
		RecipientPhone:  req.RecipientPhone,
		DeliveryAddress: req.DeliveryAddress,
		Notes:           req.Notes,
	}

	result, err := h.createShipment.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Created(c, result)
}

// GetShipment handles GET /shipments/:id
func (h *ShipmentHandler) GetShipment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid shipment ID"))
		return
	}

	result, err := h.getShipment.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("shipment"))
		return
	}

	response.Success(c, result)
}

// ListShipments handles GET /shipments
func (h *ShipmentHandler) ListShipments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := &repository.ShipmentFilter{
		Status:   entity.ShipmentStatus(c.Query("status")),
		Carrier:  c.Query("carrier"),
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Page:     page,
		Limit:    limit,
	}

	if orderID := c.Query("sales_order_id"); orderID != "" {
		if id, err := uuid.Parse(orderID); err == nil {
			filter.SalesOrderID = &id
		}
	}

	results, total, err := h.listShipments.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	meta := response.NewMeta(page, limit, total)
	response.SuccessWithMeta(c, results, meta)
}

// ShipRequest represents ship request
type ShipRequest struct {
	Carrier        string `json:"carrier"`
	TrackingNumber string `json:"tracking_number"`
}

// ShipShipment handles PATCH /shipments/:id/ship
func (h *ShipmentHandler) ShipShipment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid shipment ID"))
		return
	}

	var req ShipRequest
	c.ShouldBindJSON(&req)

	input := &shipment.ShipInput{
		ShipmentID:     id,
		Carrier:        req.Carrier,
		TrackingNumber: req.TrackingNumber,
	}

	result, err := h.shipShipment.Execute(c.Request.Context(), input)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}

// DeliverShipment handles PATCH /shipments/:id/deliver
func (h *ShipmentHandler) DeliverShipment(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("invalid shipment ID"))
		return
	}

	result, err := h.deliverShipment.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.BadRequest(err.Error()))
		return
	}

	response.Success(c, result)
}
