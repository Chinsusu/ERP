package entity_test

import (
	"testing"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestLot_IsExpiringSoon(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		expiryDate time.Time
		days       int
		expected   bool
	}{
		{
			name:       "Expiring within 30 days",
			expiryDate: now.AddDate(0, 0, 15),
			days:       30,
			expected:   true,
		},
		{
			name:       "Not expiring within 30 days",
			expiryDate: now.AddDate(0, 0, 45),
			days:       30,
			expected:   false,
		},
		{
			name:       "Already expired",
			expiryDate: now.AddDate(0, 0, -1),
			days:       30,
			expected:   false, // implementation returns false for expired
		},
		{
			name:       "Expiring exactly at threshold",
			expiryDate: now.AddDate(0, 0, 30).Add(time.Hour), // Add an hour to be sure it's AFTER threshold
			days:       30,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := &entity.Lot{ExpiryDate: tt.expiryDate}
			assert.Equal(t, tt.expected, lot.IsExpiringSoon(tt.days))
		})
	}
}

func TestLot_IsExpired(t *testing.T) {
	now := time.Now()
	
	lotPast := &entity.Lot{ExpiryDate: now.AddDate(0, 0, -1)}
	assert.True(t, lotPast.IsExpired())

	lotFuture := &entity.Lot{ExpiryDate: now.AddDate(0, 0, 1)}
	assert.False(t, lotFuture.IsExpired())
}
