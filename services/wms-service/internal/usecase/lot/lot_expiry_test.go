package lot_test

import (
	"testing"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestLot_ExpiryAlerts(t *testing.T) {
	// Given: A lot expiring in 15 days
	now := time.Now()
	expiry := now.AddDate(0, 0, 15)
	
	lot := &entity.Lot{
		ExpiryDate: expiry,
	}

	// Then: It should be expiring soon within 30 days
	assert.True(t, lot.IsExpiringSoon(30))
	assert.True(t, lot.IsExpiringSoon(20))
	assert.False(t, lot.IsExpiringSoon(10))
	
	// And not expired yet
	assert.False(t, lot.IsExpired())
}

func TestLot_AutoBlockOnExpiry(t *testing.T) {
	// Given: An expired lot
	lot := &entity.Lot{
		Status:   entity.LotStatusAvailable,
		QCStatus: entity.QCStatusPassed,
	}

	// When: Mark as expired
	lot.MarkExpired()

	// Then: Status should be EXPIRED
	assert.Equal(t, entity.LotStatusExpired, lot.Status)
	assert.False(t, lot.CanBeIssued())
}
