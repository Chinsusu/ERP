package entity_test

import (
	"testing"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIssueStockFEFO_SingleLot_Success(t *testing.T) {
	// Given: 1 lot với 100 units, expiry 2027-12-31
	materialID := uuid.New()
	expiry := time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)
	
	lot := testutils.NewLotBuilder().
		WithMaterialID(materialID).
		WithExpiry(expiry).
		Build()
	
	stock := testutils.NewStockBuilder().
		WithMaterial(materialID).
		WithLot(lot).
		WithQuantity(100).
		Build()
	
	stocks := []*entity.Stock{stock}

	// When: Issue 50 units
	issued, remaining, err := entity.AllocateStockFEFO(stocks, 50)

	// Then: Issue từ lot đó, remaining = 50
	assert.NoError(t, err)
	assert.Equal(t, 0.0, remaining)
	assert.Len(t, issued, 1)
	assert.Equal(t, 50.0, issued[0].Quantity)
	assert.Equal(t, lot.ID, issued[0].LotID)
}

func TestIssueStockFEFO_MultipleLots_IssueFromEarliestExpiry(t *testing.T) {
	// Given:
	//   - Lot A: 50 units, expiry 2027-06-30 (sớm hơn)
	//   - Lot B: 100 units, expiry 2027-12-31
	materialID := uuid.New()
	
	lotA := testutils.NewLotBuilder().
		WithMaterialID(materialID).
		WithLotNumber("LOT-A").
		WithExpiry(time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)).
		Build()
	
	lotB := testutils.NewLotBuilder().
		WithMaterialID(materialID).
		WithLotNumber("LOT-B").
		WithExpiry(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)).
		Build()

	stockA := testutils.NewStockBuilder().WithLot(lotA).WithQuantity(50).Build()
	stockB := testutils.NewStockBuilder().WithLot(lotB).WithQuantity(100).Build()
	
	// Pre-sorted list for FEFO
	stocks := []*entity.Stock{stockA, stockB}

	// When: Issue 70 units
	issued, remaining, err := entity.AllocateStockFEFO(stocks, 70)

	// Then: 
	//   - Issue 50 từ Lot A (hết)
	//   - Issue 20 từ Lot B
	assert.NoError(t, err)
	assert.Equal(t, 0.0, remaining)
	assert.Len(t, issued, 2)
	assert.Equal(t, 50.0, issued[0].Quantity)
	assert.Equal(t, lotA.ID, issued[0].LotID)
	assert.Equal(t, 20.0, issued[1].Quantity)
	assert.Equal(t, lotB.ID, issued[1].LotID)
}

func TestIssueStockFEFO_SkipExpiredLots(t *testing.T) {
	// Given:
	//   - Lot A: 50 units, expiry 2024-01-01 (ĐÃ HẾT HẠN)
	//   - Lot B: 100 units, expiry 2027-12-31
	materialID := uuid.New()
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	
	lotA := testutils.NewLotBuilder().
		WithMaterialID(materialID).
		WithExpiry(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
		Build()
	
	lotB := testutils.NewLotBuilder().
		WithMaterialID(materialID).
		WithExpiry(time.Date(2027, 12, 31, 0, 0, 0, 0, time.UTC)).
		Build()

	stockA := testutils.NewStockBuilder().WithLot(lotA).WithQuantity(50).Build()
	stockB := testutils.NewStockBuilder().WithLot(lotB).WithQuantity(100).Build()
	
	allStocks := []*entity.Stock{stockA, stockB}

	// When: Filter stage
	availableStocks := entity.FilterAndSortForFEFO(allStocks, now)

	// Then: Only Lot B should be present
	assert.Len(t, availableStocks, 1)
	assert.Equal(t, lotB.ID, *availableStocks[0].LotID)

	// When: Issue 30 units
	issued, remaining, err := entity.AllocateStockFEFO(availableStocks, 30)

	// Then: Issue từ Lot B (skip Lot A)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, remaining)
	assert.Len(t, issued, 1)
	assert.Equal(t, 30.0, issued[0].Quantity)
	assert.Equal(t, lotB.ID, issued[0].LotID)
}

func TestIssueStockFEFO_InsufficientStock_Error(t *testing.T) {
	// Given: Total available = 100 units
	materialID := uuid.New()
	lot := testutils.NewLotBuilder().WithMaterialID(materialID).Build()
	stock := testutils.NewStockBuilder().WithLot(lot).WithQuantity(100).Build()
	stocks := []*entity.Stock{stock}

	// When: Issue 150 units
	issued, remaining, err := entity.AllocateStockFEFO(stocks, 150)

	// Then: Check remaining
	assert.NoError(t, err)
	assert.Equal(t, 50.0, remaining)
	assert.Len(t, issued, 1) // Lot A exhausted
}

func TestIssueStockFEFO_ReservedStockNotIssued(t *testing.T) {
	// Given:
	//   - Lot A: quantity=100, reserved=30, available=70
	materialID := uuid.New()
	lot := testutils.NewLotBuilder().WithMaterialID(materialID).Build()
	stock := testutils.NewStockBuilder().
		WithLot(lot).
		WithQuantity(100).
		WithReserved(30).
		Build()
	stocks := []*entity.Stock{stock}

	// When: Issue 80 units
	issued, remaining, err := entity.AllocateStockFEFO(stocks, 80)

	// Then: Chỉ có 70 available để issue
	assert.NoError(t, err)
	assert.Equal(t, 10.0, remaining)
	assert.Len(t, issued, 1)
	assert.Equal(t, 70.0, issued[0].Quantity)
}

func TestIssueStockFEFO_MultipleLocations(t *testing.T) {
	// Given: Same lot ở multiple locations
	materialID := uuid.New()
	lot := testutils.NewLotBuilder().WithMaterialID(materialID).Build()
	
	stockLoc1 := testutils.NewStockBuilder().WithLot(lot).WithQuantity(40).WithLocation(uuid.New()).Build()
	stockLoc2 := testutils.NewStockBuilder().WithLot(lot).WithQuantity(60).WithLocation(uuid.New()).Build()
	
	stocks := []*entity.Stock{stockLoc1, stockLoc2}

	// When: Issue 80 units
	issued, remaining, err := entity.AllocateStockFEFO(stocks, 80)

	// Then: Issue từ tất cả locations của lot đó
	assert.NoError(t, err)
	assert.Equal(t, 0.0, remaining)
	assert.Len(t, issued, 2)
	assert.Equal(t, 40.0, issued[0].Quantity)
	assert.Equal(t, 40.0, issued[1].Quantity)
	assert.Equal(t, lot.ID, issued[0].LotID)
	assert.Equal(t, lot.ID, issued[1].LotID)
}
