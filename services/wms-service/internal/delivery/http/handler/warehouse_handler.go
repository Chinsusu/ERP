package handler

import (
	"strconv"

	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/usecase/lot"
	"github.com/erp-cosmetics/wms-service/internal/usecase/stock"
	"github.com/erp-cosmetics/wms-service/internal/usecase/warehouse"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WarehouseHandler handles warehouse endpoints
type WarehouseHandler struct {
	listWarehousesUC *warehouse.ListWarehousesUseCase
	getWarehouseUC   *warehouse.GetWarehouseUseCase
	getZonesUC       *warehouse.GetZonesUseCase
	getLocationsUC   *warehouse.GetLocationsUseCase
}

// NewWarehouseHandler creates a new warehouse handler
func NewWarehouseHandler(
	listWarehousesUC *warehouse.ListWarehousesUseCase,
	getWarehouseUC *warehouse.GetWarehouseUseCase,
	getZonesUC *warehouse.GetZonesUseCase,
	getLocationsUC *warehouse.GetLocationsUseCase,
) *WarehouseHandler {
	return &WarehouseHandler{
		listWarehousesUC: listWarehousesUC,
		getWarehouseUC:   getWarehouseUC,
		getZonesUC:       getZonesUC,
		getLocationsUC:   getLocationsUC,
	}
}

// ListWarehouses handles GET /warehouses
func (h *WarehouseHandler) ListWarehouses(c *gin.Context) {
	filter := &repository.WarehouseFilter{
		Type:   c.Query("type"),
		Search: c.Query("search"),
		Page:   getPageParam(c),
		Limit:  getLimitParam(c),
	}

	warehouses, total, err := h.listWarehousesUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.SuccessWithMeta(c, warehouses, response.NewMeta(filter.Page, filter.Limit, int64(total)))
}

// GetWarehouse handles GET /warehouses/:id
func (h *WarehouseHandler) GetWarehouse(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid warehouse ID"))
		return
	}

	wh, err := h.getWarehouseUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Warehouse"))
		return
	}

	response.Success(c, wh)
}

// GetZones handles GET /warehouses/:id/zones
func (h *WarehouseHandler) GetZones(c *gin.Context) {
	warehouseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid warehouse ID"))
		return
	}

	zones, err := h.getZonesUC.Execute(c.Request.Context(), warehouseID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, zones)
}

// GetLocations handles GET /zones/:id/locations
func (h *WarehouseHandler) GetLocations(c *gin.Context) {
	zoneID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid zone ID"))
		return
	}

	locations, err := h.getLocationsUC.Execute(c.Request.Context(), zoneID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, locations)
}

// StockHandler handles stock endpoints
type StockHandler struct {
	getStockUC         *stock.GetStockUseCase
	issueStockFEFOUC   *stock.IssueStockFEFOUseCase
	reserveStockUC     *stock.ReserveStockUseCase
	releaseReservationUC *stock.ReleaseReservationUseCase
}

// NewStockHandler creates a new stock handler
func NewStockHandler(
	getStockUC *stock.GetStockUseCase,
	issueStockFEFOUC *stock.IssueStockFEFOUseCase,
	reserveStockUC *stock.ReserveStockUseCase,
	releaseReservationUC *stock.ReleaseReservationUseCase,
) *StockHandler {
	return &StockHandler{
		getStockUC:           getStockUC,
		issueStockFEFOUC:     issueStockFEFOUC,
		reserveStockUC:       reserveStockUC,
		releaseReservationUC: releaseReservationUC,
	}
}

// ListStock handles GET /stock
func (h *StockHandler) ListStock(c *gin.Context) {
	filter := &repository.StockFilter{
		Page:  getPageParam(c),
		Limit: getLimitParam(c),
	}

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		id, _ := uuid.Parse(warehouseID)
		filter.WarehouseID = &id
	}
	if materialID := c.Query("material_id"); materialID != "" {
		id, _ := uuid.Parse(materialID)
		filter.MaterialID = &id
	}
	if hasStock := c.Query("has_stock"); hasStock == "true" {
		t := true
		filter.HasStock = &t
	}

	stocks, total, err := h.getStockUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.SuccessWithMeta(c, stocks, response.NewMeta(filter.Page, filter.Limit, int64(total)))
}

// GetStockByMaterial handles GET /stock/by-material/:id
func (h *StockHandler) GetStockByMaterial(c *gin.Context) {
	materialID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid material ID"))
		return
	}

	stocks, err := h.getStockUC.GetByMaterial(c.Request.Context(), materialID)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	summary, _ := h.getStockUC.GetSummary(c.Request.Context(), materialID)

	response.Success(c, gin.H{
		"stocks":  stocks,
		"summary": summary,
	})
}

// GetExpiringStock handles GET /stock/expiring
func (h *StockHandler) GetExpiringStock(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))

	stocks, err := h.getStockUC.GetExpiringStock(c.Request.Context(), days)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, stocks)
}

// GetLowStock handles GET /stock/low-stock
func (h *StockHandler) GetLowStock(c *gin.Context) {
	threshold, _ := strconv.ParseFloat(c.DefaultQuery("threshold", "100"), 64)

	stocks, err := h.getStockUC.GetLowStock(c.Request.Context(), threshold)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, stocks)
}

// LotHandler handles lot endpoints
type LotHandler struct {
	getLotUC          *lot.GetLotUseCase
	listLotsUC        *lot.ListLotsUseCase
	getExpiringLotsUC *lot.GetExpiringLotsUseCase
	getLotMovementsUC *lot.GetLotMovementsUseCase
}

// NewLotHandler creates a new lot handler
func NewLotHandler(
	getLotUC *lot.GetLotUseCase,
	listLotsUC *lot.ListLotsUseCase,
	getExpiringLotsUC *lot.GetExpiringLotsUseCase,
	getLotMovementsUC *lot.GetLotMovementsUseCase,
) *LotHandler {
	return &LotHandler{
		getLotUC:          getLotUC,
		listLotsUC:        listLotsUC,
		getExpiringLotsUC: getExpiringLotsUC,
		getLotMovementsUC: getLotMovementsUC,
	}
}

// ListLots handles GET /lots
func (h *LotHandler) ListLots(c *gin.Context) {
	filter := &repository.LotFilter{
		Status:   c.Query("status"),
		QCStatus: c.Query("qc_status"),
		Search:   c.Query("search"),
		Page:     getPageParam(c),
		Limit:    getLimitParam(c),
	}

	if materialID := c.Query("material_id"); materialID != "" {
		id, _ := uuid.Parse(materialID)
		filter.MaterialID = &id
	}

	lots, total, err := h.listLotsUC.Execute(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.SuccessWithMeta(c, lots, response.NewMeta(filter.Page, filter.Limit, int64(total)))
}

// GetLot handles GET /lots/:id
func (h *LotHandler) GetLot(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid lot ID"))
		return
	}

	l, err := h.getLotUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.NotFound("Lot"))
		return
	}

	response.Success(c, l)
}

// GetLotMovements handles GET /lots/:id/movements
func (h *LotHandler) GetLotMovements(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.BadRequest("Invalid lot ID"))
		return
	}

	movements, err := h.getLotMovementsUC.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.Internal(err))
		return
	}

	response.Success(c, movements)
}

// Helper functions
func getPageParam(c *gin.Context) int {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	return page
}

func getLimitParam(c *gin.Context) int {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return limit
}
