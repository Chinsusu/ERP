package testutils

import (
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/google/uuid"
)

// LotBuilder helps created test Lot entities
type LotBuilder struct {
	lot *entity.Lot
}

func NewLotBuilder() *LotBuilder {
	return &LotBuilder{
		lot: &entity.Lot{
			ID:           uuid.New(),
			LotNumber:    "LOT-" + uuid.New().String()[:8],
			QCStatus:     entity.QCStatusPassed,
			Status:       entity.LotStatusAvailable,
			ReceivedDate: time.Now(),
			ExpiryDate:   time.Now().AddDate(0, 0, 365),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
}

func (b *LotBuilder) WithID(id uuid.UUID) *LotBuilder {
	b.lot.ID = id
	return b
}

func (b *LotBuilder) WithLotNumber(n string) *LotBuilder {
	b.lot.LotNumber = n
	return b
}

func (b *LotBuilder) WithMaterialID(id uuid.UUID) *LotBuilder {
	b.lot.MaterialID = id
	return b
}

func (b *LotBuilder) WithExpiry(date time.Time) *LotBuilder {
	b.lot.ExpiryDate = date
	return b
}

func (b *LotBuilder) WithQCStatus(status entity.QCStatus) *LotBuilder {
	b.lot.QCStatus = status
	return b
}

func (b *LotBuilder) WithStatus(status entity.LotStatus) *LotBuilder {
	b.lot.Status = status
	return b
}

func (b *LotBuilder) Build() *entity.Lot {
	return b.lot
}

// StockBuilder helps create test Stock entities
type StockBuilder struct {
	stock *entity.Stock
}

func NewStockBuilder() *StockBuilder {
	return &StockBuilder{
		stock: &entity.Stock{
			ID:          uuid.New(),
			WarehouseID: uuid.New(),
			ZoneID:      uuid.New(),
			LocationID:  uuid.New(),
			Quantity:    100,
			ReservedQty: 0,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

func (b *StockBuilder) WithLocation(id uuid.UUID) *StockBuilder {
	b.stock.LocationID = id
	return b
}

func (b *StockBuilder) WithMaterial(id uuid.UUID) *StockBuilder {
	b.stock.MaterialID = id
	return b
}

func (b *StockBuilder) WithLot(lot *entity.Lot) *StockBuilder {
	if lot != nil {
		b.stock.LotID = &lot.ID
		b.stock.Lot = lot
		b.stock.MaterialID = lot.MaterialID
	}
	return b
}

func (b *StockBuilder) WithQuantity(qty float64) *StockBuilder {
	b.stock.Quantity = qty
	b.stock.AvailableQty = qty - b.stock.ReservedQty
	return b
}

func (b *StockBuilder) WithReserved(qty float64) *StockBuilder {
	b.stock.ReservedQty = qty
	b.stock.AvailableQty = b.stock.Quantity - qty
	return b
}

func (b *StockBuilder) Build() *entity.Stock {
	b.stock.AvailableQty = b.stock.Quantity - b.stock.ReservedQty
	return b.stock
}
