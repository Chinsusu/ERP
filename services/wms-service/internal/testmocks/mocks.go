package testmocks

import (
	"context"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/infrastructure/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockStockRepository is a mock implementation of repository.StockRepository
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) Create(ctx context.Context, stock *entity.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Stock, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Stock), args.Error(1)
}

func (m *MockStockRepository) Update(ctx context.Context, stock *entity.Stock) error {
	args := m.Called(ctx, stock)
	return args.Error(0)
}

func (m *MockStockRepository) GetAvailableStockFEFO(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error) {
	args := m.Called(ctx, materialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Stock), args.Error(1)
}

func (m *MockStockRepository) IssueStockFEFO(ctx context.Context, materialID uuid.UUID, quantity float64, createdBy uuid.UUID) ([]entity.LotIssued, error) {
	args := m.Called(ctx, materialID, quantity, createdBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.LotIssued), args.Error(1)
}

func (m *MockStockRepository) CreateMovement(ctx context.Context, movement *entity.StockMovement) error {
	args := m.Called(ctx, movement)
	return args.Error(0)
}

func (m *MockStockRepository) GetNextMovementNumber(ctx context.Context, movementType entity.MovementType) (string, error) {
	args := m.Called(ctx, movementType)
	return args.String(0), args.Error(1)
}

// Implement other methods
func (m *MockStockRepository) GetByLocation(ctx context.Context, locationID uuid.UUID) ([]*entity.Stock, error) { return nil, nil }
func (m *MockStockRepository) GetByMaterial(ctx context.Context, materialID uuid.UUID) ([]*entity.Stock, error) { return nil, nil }
func (m *MockStockRepository) GetByMaterialAndLot(ctx context.Context, materialID, lotID uuid.UUID) (*entity.Stock, error) { return nil, nil }
func (m *MockStockRepository) GetByLocationMaterialLot(ctx context.Context, locationID, materialID uuid.UUID, lotID *uuid.UUID) (*entity.Stock, error) { return nil, nil }
func (m *MockStockRepository) List(ctx context.Context, filter *repository.StockFilter) ([]*entity.Stock, int64, error) { return nil, 0, nil }
func (m *MockStockRepository) ReceiveStock(ctx context.Context, stock *entity.Stock, movement *entity.StockMovement) error { 
	args := m.Called(ctx, stock, movement)
	return args.Error(0)
}
func (m *MockStockRepository) IssueStock(ctx context.Context, stock *entity.Stock, movement *entity.StockMovement) error { return nil }
func (m *MockStockRepository) TransferStock(ctx context.Context, fromStock, toStock *entity.Stock, movement *entity.StockMovement) error { return nil }
func (m *MockStockRepository) AdjustStock(ctx context.Context, stock *entity.Stock, adjustmentQty float64, movement *entity.StockMovement) error { return nil }
func (m *MockStockRepository) ReserveStock(ctx context.Context, materialID uuid.UUID, quantity float64, reservation *entity.StockReservation) error {
	args := m.Called(ctx, materialID, quantity, reservation)
	return args.Error(0)
}
func (m *MockStockRepository) ReleaseReservation(ctx context.Context, reservationID uuid.UUID) error {
	args := m.Called(ctx, reservationID)
	return args.Error(0)
}
func (m *MockStockRepository) GetMaterialSummary(ctx context.Context, materialID uuid.UUID) (*entity.StockSummary, error) {
	args := m.Called(ctx, materialID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.StockSummary), args.Error(1)
}
func (m *MockStockRepository) GetLowStockMaterials(ctx context.Context, threshold float64) ([]*entity.StockSummary, error) { return nil, nil }
func (m *MockStockRepository) GetExpiringStock(ctx context.Context, days int) ([]*entity.Stock, error) { return nil, nil }
func (m *MockStockRepository) GetMovementsByLot(ctx context.Context, lotID uuid.UUID) ([]*entity.StockMovement, error) { return nil, nil }
func (m *MockStockRepository) GetMovementsByMaterial(ctx context.Context, materialID uuid.UUID, limit int) ([]*entity.StockMovement, error) { return nil, nil }

// MockGoodsIssueRepository
type MockGoodsIssueRepository struct {
	mock.Mock
}

func (m *MockGoodsIssueRepository) Create(ctx context.Context, issue *entity.GoodsIssue) error {
	args := m.Called(ctx, issue)
	return args.Error(0)
}
func (m *MockGoodsIssueRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.GoodsIssue, error) { return nil, nil }
func (m *MockGoodsIssueRepository) GetByNumber(ctx context.Context, num string) (*entity.GoodsIssue, error) { return nil, nil }
func (m *MockGoodsIssueRepository) List(ctx context.Context, filter *repository.GoodsIssueFilter) ([]*entity.GoodsIssue, int64, error) { return nil, 0, nil }
func (m *MockGoodsIssueRepository) Update(ctx context.Context, issue *entity.GoodsIssue) error {
	args := m.Called(ctx, issue)
	return args.Error(0)
}
func (m *MockGoodsIssueRepository) GetNextIssueNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
func (m *MockGoodsIssueRepository) CreateLineItem(ctx context.Context, item *entity.GILineItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *MockGoodsIssueRepository) CreateLineItems(ctx context.Context, items []*entity.GILineItem) error {
	args := m.Called(ctx, items)
	return args.Error(0)
}
func (m *MockGoodsIssueRepository) GetLineItemsByIssueID(ctx context.Context, id uuid.UUID) ([]*entity.GILineItem, error) { return nil, nil }

// MockGRNRepository
type MockGRNRepository struct {
	mock.Mock
}

func (m *MockGRNRepository) Create(ctx context.Context, grn *entity.GRN) error {
	args := m.Called(ctx, grn)
	return args.Error(0)
}
func (m *MockGRNRepository) Update(ctx context.Context, grn *entity.GRN) error {
	args := m.Called(ctx, grn)
	return args.Error(0)
}
func (m *MockGRNRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.GRN, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.GRN), args.Error(1)
}
func (m *MockGRNRepository) GetByNumber(ctx context.Context, num string) (*entity.GRN, error) { return nil, nil }
func (m *MockGRNRepository) GetByPOID(ctx context.Context, poID uuid.UUID) ([]*entity.GRN, error) { return nil, nil }
func (m *MockGRNRepository) CreateLineItem(ctx context.Context, item *entity.GRNLineItem) error { 
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *MockGRNRepository) GetNextGRNNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
func (m *MockGRNRepository) List(ctx context.Context, filter *repository.GRNFilter) ([]*entity.GRN, int64, error) { return nil, 0, nil }
func (m *MockGRNRepository) GetLineItemsByGRNID(ctx context.Context, grnID uuid.UUID) ([]*entity.GRNLineItem, error) {
	args := m.Called(ctx, grnID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.GRNLineItem), args.Error(1)
}
func (m *MockGRNRepository) UpdateLineItem(ctx context.Context, item *entity.GRNLineItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

// MockZoneRepository
type MockZoneRepository struct {
	mock.Mock
}

func (m *MockZoneRepository) GetQuarantineZone(ctx context.Context, whID uuid.UUID) (*entity.Zone, error) {
	args := m.Called(ctx, whID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Zone), args.Error(1)
}
func (m *MockZoneRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Zone, error) { return nil, nil }
func (m *MockZoneRepository) GetByWarehouseID(ctx context.Context, whID uuid.UUID) ([]*entity.Zone, error) { return nil, nil }
func (m *MockZoneRepository) GetStorageZone(ctx context.Context, whID uuid.UUID) (*entity.Zone, error) { return nil, nil }
func (m *MockZoneRepository) Create(ctx context.Context, zone *entity.Zone) error { return nil }
func (m *MockZoneRepository) Update(ctx context.Context, zone *entity.Zone) error { return nil }
func (m *MockZoneRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }

// MockLocationRepository
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) GetByZoneID(ctx context.Context, zoneID uuid.UUID) ([]*entity.Location, error) {
	args := m.Called(ctx, zoneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Location), args.Error(1)
}
func (m *MockLocationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Location, error) { return nil, nil }
func (m *MockLocationRepository) GetByCode(ctx context.Context, zID uuid.UUID, code string) (*entity.Location, error) { return nil, nil }
func (m *MockLocationRepository) Create(ctx context.Context, location *entity.Location) error { return nil }
func (m *MockLocationRepository) Update(ctx context.Context, location *entity.Location) error { return nil }
func (m *MockLocationRepository) Delete(ctx context.Context, id uuid.UUID) error { return nil }

// MockLotRepository
type MockLotRepository struct {
	mock.Mock
}

func (m *MockLotRepository) Create(ctx context.Context, lot *entity.Lot) error {
	args := m.Called(ctx, lot)
	return args.Error(0)
}
func (m *MockLotRepository) Update(ctx context.Context, lot *entity.Lot) error {
	args := m.Called(ctx, lot)
	return args.Error(0)
}
func (m *MockLotRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Lot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Lot), args.Error(1)
}
func (m *MockLotRepository) GetByLotNumber(ctx context.Context, num string) (*entity.Lot, error) { return nil, nil }
func (m *MockLotRepository) List(ctx context.Context, filter *repository.LotFilter) ([]*entity.Lot, int64, error) { return nil, 0, nil }
func (m *MockLotRepository) GetAvailableLots(ctx context.Context, matID uuid.UUID) ([]*entity.Lot, error) { return nil, nil }
func (m *MockLotRepository) GetExpiringLots(ctx context.Context, days int) ([]*entity.Lot, error) { return nil, nil }
func (m *MockLotRepository) GetExpiredLots(ctx context.Context) ([]*entity.Lot, error) { return nil, nil }
func (m *MockLotRepository) GetNextLotNumber(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
func (m *MockLotRepository) MarkExpired(ctx context.Context, ids []uuid.UUID) error { return nil }

// MockEventPublisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) PublishGRNCreated(e *event.GRNCreatedEvent) error {
	args := m.Called(e)
	return args.Error(0)
}
func (m *MockEventPublisher) PublishGRNCompleted(e *event.GRNCompletedEvent) error {
	args := m.Called(e)
	return args.Error(0)
}
func (m *MockEventPublisher) PublishStockReceived(e *event.StockReceivedEvent) error {
	args := m.Called(e)
	return args.Error(0)
}
func (m *MockEventPublisher) PublishStockIssued(e *event.StockIssuedEvent) error {
	args := m.Called(e)
	return args.Error(0)
}
func (m *MockEventPublisher) PublishStockReserved(e *event.StockReservedEvent) error {
	args := m.Called(e)
	return args.Error(0)
}
func (m *MockEventPublisher) PublishLowStockAlert(e *event.LowStockAlertEvent) error { return nil }
func (m *MockEventPublisher) PublishLotExpiringSoon(e *event.LotExpiringEvent) error { return nil }
func (m *MockEventPublisher) PublishLotExpired(e *event.LotExpiringEvent) error { return nil }
