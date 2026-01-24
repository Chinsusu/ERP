package postgres

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

// NewStockRepository creates a new stock repository
func NewStockRepository(db *gorm.DB) repository.StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(ctx context.Context, stock *entity.Stock) error {
	return r.db.WithContext(ctx).Create(stock).Error
}

func (r *stockRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Stock, error) {
	var stock entity.Stock
	err := r.db.WithContext(ctx).
		Preload("Lot").
		Preload("Location").
		First(&stock, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) Update(ctx context.Context, stock *entity.Stock) error {
	stock.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(stock).Error
}

func (r *stockRepository) GetByLocation(ctx context.Context, locationID uuid.UUID) ([]*entity.Stock, error) {
	var stocks []*entity.Stock
	err := r.db.WithContext(ctx).
		Preload("Lot").
		Where("location_id = ?", locationID).
		Where("quantity > 0").
		Find(&stocks).Error
	return stocks, err
}

func (r *stockRepository) GetByMaterial(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error) {
	var stocks []*entity.Stock
	err := r.db.WithContext(ctx).
		Preload("Lot").
		Preload("Location").
		Preload("Location.Zone").
		Where("material_id = ?", materialID).
		Where("quantity > 0").
		Find(&stocks).Error
	return stocks, err
}

func (r *stockRepository) GetByMaterialAndLot(ctx context.Context, materialID, lotID uuid.UUID) (*entity.Stock, error) {
	var stock entity.Stock
	err := r.db.WithContext(ctx).
		Where("material_id = ? AND lot_id = ?", materialID, lotID).
		First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) GetByLocationMaterialLot(ctx context.Context, locationID, materialID uuid.UUID, lotID *uuid.UUID) (*entity.Stock, error) {
	var stock entity.Stock
	query := r.db.WithContext(ctx).
		Where("location_id = ? AND material_id = ?", locationID, materialID)

	if lotID != nil {
		query = query.Where("lot_id = ?", *lotID)
	} else {
		query = query.Where("lot_id IS NULL")
	}

	err := query.First(&stock).Error
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) List(ctx context.Context, filter *repository.StockFilter) ([]*entity.Stock, int64, error) {
	var stocks []*entity.Stock
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Stock{})

	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.ZoneID != nil {
		query = query.Where("zone_id = ?", *filter.ZoneID)
	}
	if filter.LocationID != nil {
		query = query.Where("location_id = ?", *filter.LocationID)
	}
	if filter.MaterialID != nil {
		query = query.Where("material_id = ?", *filter.MaterialID)
	}
	if filter.LotID != nil {
		query = query.Where("lot_id = ?", *filter.LotID)
	}
	if filter.HasStock != nil && *filter.HasStock {
		query = query.Where("quantity > 0")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(50)
	}
	if filter.Page > 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if err := query.
		Preload("Lot").
		Preload("Location").
		Preload("Location.Zone").
		Find(&stocks).Error; err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}

// GetAvailableStockFEFO returns stock sorted by lot expiry date (FEFO)
// This is the CRITICAL method for cosmetics - First Expired First Out
func (r *stockRepository) GetAvailableStockFEFO(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error) {
	var stocks []*entity.Stock
	err := r.db.WithContext(ctx).
		Joins("JOIN lots ON lots.id = stock.lot_id").
		Where("stock.material_id = ?", materialID).
		Where("stock.quantity - stock.reserved_qty > 0"). // Has available qty
		Where("lots.status = ?", entity.LotStatusAvailable).
		Where("lots.qc_status = ?", entity.QCStatusPassed).
		Where("lots.expiry_date > ?", time.Now()). // Not expired
		Order("lots.expiry_date ASC"). // FEFO: earliest expiry first!
		Preload("Lot").
		Preload("Location").
		Find(&stocks).Error
	return stocks, err
}

// IssueStockFEFO issues stock using FEFO logic - THE CORE ALGORITHM
func (r *stockRepository) IssueStockFEFO(ctx context.Context, materialID uuid.UUID, quantity float64, createdBy uuid.UUID) ([]entity.LotIssued, error) {
	// Get available stocks sorted by expiry (earliest first)
	stocks, err := r.GetAvailableStockFEFO(ctx, materialID)
	if err != nil {
		return nil, err
	}

	remaining := quantity
	var issued []entity.LotIssued

	// Start transaction
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, stock := range stocks {
		if remaining <= 0 {
			break
		}

		availableQty := stock.Quantity - stock.ReservedQty
		issueQty := math.Min(availableQty, remaining)

		// Deduct from stock
		stock.Quantity -= issueQty
		stock.UpdatedAt = time.Now()
		if err := tx.Save(stock).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Record the issued lot
		if stock.Lot != nil && stock.LotID != nil {
			issued = append(issued, entity.LotIssued{
				LotID:      *stock.LotID,
				LotNumber:  stock.Lot.LotNumber,
				Quantity:   issueQty,
				ExpiryDate: stock.Lot.ExpiryDate,
				LocationID: stock.LocationID,
			})
		}

		remaining -= issueQty
	}

	// Check if we have enough stock
	if remaining > 0 {
		tx.Rollback()
		return nil, entity.ErrInsufficientStock
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return issued, nil
}

// ReceiveStock adds stock to a location
func (r *stockRepository) ReceiveStock(ctx context.Context, stock *entity.Stock, movement *entity.StockMovement) error {
	tx := r.db.WithContext(ctx).Begin()

	// Check if stock record exists
	var existing entity.Stock
	err := tx.Where("location_id = ? AND material_id = ? AND lot_id = ?",
		stock.LocationID, stock.MaterialID, stock.LotID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new stock record
		if err := tx.Create(stock).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		tx.Rollback()
		return err
	} else {
		// Update existing stock
		existing.Quantity += stock.Quantity
		existing.UpdatedAt = time.Now()
		if err := tx.Save(&existing).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create movement record
	if err := tx.Create(movement).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// IssueStock deducts stock from a location
func (r *stockRepository) IssueStock(ctx context.Context, stock *entity.Stock, movement *entity.StockMovement) error {
	tx := r.db.WithContext(ctx).Begin()

	if stock.Quantity < 0 {
		tx.Rollback()
		return entity.ErrInsufficientStock
	}

	stock.UpdatedAt = time.Now()
	if err := tx.Save(stock).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(movement).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// TransferStock transfers stock between locations
func (r *stockRepository) TransferStock(ctx context.Context, fromStock, toStock *entity.Stock, movement *entity.StockMovement) error {
	tx := r.db.WithContext(ctx).Begin()

	// Deduct from source
	fromStock.UpdatedAt = time.Now()
	if err := tx.Save(fromStock).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Check if destination stock exists
	var existing entity.Stock
	err := tx.Where("location_id = ? AND material_id = ? AND lot_id = ?",
		toStock.LocationID, toStock.MaterialID, toStock.LotID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		if err := tx.Create(toStock).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		tx.Rollback()
		return err
	} else {
		existing.Quantity += toStock.Quantity
		existing.UpdatedAt = time.Now()
		if err := tx.Save(&existing).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Create(movement).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// AdjustStock adjusts stock quantity
func (r *stockRepository) AdjustStock(ctx context.Context, stock *entity.Stock, adjustmentQty float64, movement *entity.StockMovement) error {
	tx := r.db.WithContext(ctx).Begin()

	stock.Quantity += adjustmentQty
	if stock.Quantity < 0 {
		tx.Rollback()
		return entity.ErrInvalidQuantity
	}
	stock.UpdatedAt = time.Now()

	if err := tx.Save(stock).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(movement).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ReserveStock reserves stock for an order
func (r *stockRepository) ReserveStock(ctx context.Context, materialID uuid.UUID, quantity float64, reservation *entity.StockReservation) error {
	tx := r.db.WithContext(ctx).Begin()

	// Get available stock using FEFO
	stocks, err := r.GetAvailableStockFEFO(ctx, materialID)
	if err != nil {
		tx.Rollback()
		return err
	}

	remaining := quantity
	for _, stock := range stocks {
		if remaining <= 0 {
			break
		}

		availableQty := stock.Quantity - stock.ReservedQty
		reserveQty := math.Min(availableQty, remaining)

		stock.ReservedQty += reserveQty
		stock.UpdatedAt = time.Now()
		if err := tx.Save(stock).Error; err != nil {
			tx.Rollback()
			return err
		}

		remaining -= reserveQty
	}

	if remaining > 0 {
		tx.Rollback()
		return entity.ErrInsufficientStock
	}

	if err := tx.Create(reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// ReleaseReservation releases a stock reservation
func (r *stockRepository) ReleaseReservation(ctx context.Context, reservationID uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()

	var reservation entity.StockReservation
	if err := tx.First(&reservation, "id = ?", reservationID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Find and release stock
	var stock entity.Stock
	query := tx.Where("material_id = ?", reservation.MaterialID)
	if reservation.LotID != nil {
		query = query.Where("lot_id = ?", *reservation.LotID)
	}
	if reservation.LocationID != nil {
		query = query.Where("location_id = ?", *reservation.LocationID)
	}

	if err := query.First(&stock).Error; err == nil {
		stock.ReservedQty -= reservation.Quantity
		if stock.ReservedQty < 0 {
			stock.ReservedQty = 0
		}
		stock.UpdatedAt = time.Now()
		tx.Save(&stock)
	}

	// Update reservation
	reservation.Release()
	if err := tx.Save(&reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetMaterialSummary returns aggregated stock for a material
func (r *stockRepository) GetMaterialSummary(ctx context.Context, materialID uuid.UUID) (*entity.StockSummary, error) {
	var result struct {
		TotalQuantity  float64
		TotalReserved  float64
		TotalAvailable float64
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("SUM(quantity) as total_quantity, SUM(reserved_qty) as total_reserved, SUM(quantity - reserved_qty) as total_available").
		Where("material_id = ?", materialID).
		Where("quantity > 0").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	return &entity.StockSummary{
		MaterialID:     materialID,
		TotalQuantity:  result.TotalQuantity,
		TotalReserved:  result.TotalReserved,
		TotalAvailable: result.TotalAvailable,
	}, nil
}

// GetLowStockMaterials returns materials below threshold
func (r *stockRepository) GetLowStockMaterials(ctx context.Context, threshold float64) ([]*entity.StockSummary, error) {
	var results []*entity.StockSummary

	err := r.db.WithContext(ctx).
		Model(&entity.Stock{}).
		Select("material_id, SUM(quantity) as total_quantity, SUM(reserved_qty) as total_reserved, SUM(quantity - reserved_qty) as total_available").
		Group("material_id").
		Having("SUM(quantity - reserved_qty) < ?", threshold).
		Scan(&results).Error

	return results, err
}

// GetExpiringStock returns stock with lots expiring within days
func (r *stockRepository) GetExpiringStock(ctx context.Context, days int) ([]*entity.Stock, error) {
	threshold := time.Now().AddDate(0, 0, days)
	var stocks []*entity.Stock

	err := r.db.WithContext(ctx).
		Joins("JOIN lots ON lots.id = stock.lot_id").
		Where("lots.expiry_date <= ?", threshold).
		Where("lots.expiry_date > ?", time.Now()).
		Where("stock.quantity > 0").
		Preload("Lot").
		Preload("Location").
		Order("lots.expiry_date ASC").
		Find(&stocks).Error

	return stocks, err
}

// CreateMovement creates a stock movement record
func (r *stockRepository) CreateMovement(ctx context.Context, movement *entity.StockMovement) error {
	return r.db.WithContext(ctx).Create(movement).Error
}

// GetMovementsByLot returns movements for a lot
func (r *stockRepository) GetMovementsByLot(ctx context.Context, lotID uuid.UUID) ([]*entity.StockMovement, error) {
	var movements []*entity.StockMovement
	err := r.db.WithContext(ctx).
		Where("lot_id = ?", lotID).
		Order("created_at DESC").
		Find(&movements).Error
	return movements, err
}

// GetMovementsByMaterial returns movements for a material
func (r *stockRepository) GetMovementsByMaterial(ctx context.Context, materialID uuid.UUID, limit int) ([]*entity.StockMovement, error) {
	var movements []*entity.StockMovement
	err := r.db.WithContext(ctx).
		Where("material_id = ?", materialID).
		Order("created_at DESC").
		Limit(limit).
		Find(&movements).Error
	return movements, err
}

// GetNextMovementNumber generates the next movement number
func (r *stockRepository) GetNextMovementNumber(ctx context.Context, movementType entity.MovementType) (string, error) {
	var count int64
	year := time.Now().Year()
	prefix := "MOV"

	switch movementType {
	case entity.MovementTypeIn:
		prefix = "MOV-IN"
	case entity.MovementTypeOut:
		prefix = "MOV-OUT"
	case entity.MovementTypeTransfer:
		prefix = "MOV-TRF"
	case entity.MovementTypeAdjustment:
		prefix = "MOV-ADJ"
	}

	r.db.WithContext(ctx).
		Model(&entity.StockMovement{}).
		Where("movement_number LIKE ?", fmt.Sprintf("%s-%d-%%", prefix, year)).
		Count(&count)

	return fmt.Sprintf("%s-%d-%05d", prefix, year, count+1), nil
}
