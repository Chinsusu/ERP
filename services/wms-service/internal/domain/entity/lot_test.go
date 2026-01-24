package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestLot_DaysUntilExpiry(t *testing.T) {
	tests := []struct {
		name     string
		expiry   time.Time
		expected int
	}{
		{
			name:     "10 days until expiry",
			expiry:   time.Now().AddDate(0, 0, 10),
			expected: 10,
		},
		{
			name:     "0 days (today)",
			expiry:   time.Now(),
			expected: 0,
		},
		{
			name:     "Already expired (-5 days)",
			expiry:   time.Now().AddDate(0, 0, -5),
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := &Lot{
				ExpiryDate: tt.expiry,
			}
			// Allow 1 day tolerance due to time calculation
			result := lot.DaysUntilExpiry()
			if result < tt.expected-1 || result > tt.expected+1 {
				t.Errorf("DaysUntilExpiry() = %d, expected ~%d", result, tt.expected)
			}
		})
	}
}

func TestLot_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		expiry   time.Time
		expected bool
	}{
		{
			name:     "Future expiry",
			expiry:   time.Now().AddDate(0, 0, 30),
			expected: false,
		},
		{
			name:     "Past expiry",
			expiry:   time.Now().AddDate(0, 0, -1),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := &Lot{
				ExpiryDate: tt.expiry,
			}
			if lot.IsExpired() != tt.expected {
				t.Errorf("IsExpired() = %v, expected %v", lot.IsExpired(), tt.expected)
			}
		})
	}
}

func TestLot_IsExpiringSoon(t *testing.T) {
	tests := []struct {
		name        string
		expiry      time.Time
		alertDays   int
		expected    bool
	}{
		{
			name:      "Expiring within 30 days",
			expiry:    time.Now().AddDate(0, 0, 20),
			alertDays: 30,
			expected:  true,
		},
		{
			name:      "Not expiring soon",
			expiry:    time.Now().AddDate(0, 0, 100),
			alertDays: 30,
			expected:  false,
		},
		{
			name:      "Expiring within 7 days",
			expiry:    time.Now().AddDate(0, 0, 5),
			alertDays: 7,
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lot := &Lot{
				ExpiryDate: tt.expiry,
			}
			if lot.IsExpiringSoon(tt.alertDays) != tt.expected {
				t.Errorf("IsExpiringSoon(%d) = %v, expected %v", 
					tt.alertDays, lot.IsExpiringSoon(tt.alertDays), tt.expected)
			}
		})
	}
}

func TestLot_QCWorkflow(t *testing.T) {
	lot := &Lot{
		ID:        uuid.New(),
		QCStatus:  QCStatusPending,
		Status:    LotStatusAvailable,
	}

	// Test Pass QC
	lot.PassQC()
	if lot.QCStatus != QCStatusPassed {
		t.Errorf("After PassQC(), QCStatus = %s, expected %s", lot.QCStatus, QCStatusPassed)
	}

	// Reset and test Fail QC
	lot.QCStatus = QCStatusPending
	lot.FailQC()
	if lot.QCStatus != QCStatusFailed {
		t.Errorf("After FailQC(), QCStatus = %s, expected %s", lot.QCStatus, QCStatusFailed)
	}
	if lot.Status != LotStatusBlocked {
		t.Errorf("After FailQC(), Status = %s, expected %s", lot.Status, LotStatusBlocked)
	}
}

func TestLot_CanBeIssued(t *testing.T) {
	tests := []struct {
		name     string
		lot      *Lot
		expected bool
	}{
		{
			name: "Available and QC passed",
			lot: &Lot{
				Status:     LotStatusAvailable,
				QCStatus:   QCStatusPassed,
				ExpiryDate: time.Now().AddDate(0, 0, 30),
			},
			expected: true,
		},
		{
			name: "Blocked",
			lot: &Lot{
				Status:     LotStatusBlocked,
				QCStatus:   QCStatusPassed,
				ExpiryDate: time.Now().AddDate(0, 0, 30),
			},
			expected: false,
		},
		{
			name: "QC Failed",
			lot: &Lot{
				Status:     LotStatusAvailable,
				QCStatus:   QCStatusFailed,
				ExpiryDate: time.Now().AddDate(0, 0, 30),
			},
			expected: false,
		},
		{
			name: "Expired",
			lot: &Lot{
				Status:     LotStatusAvailable,
				QCStatus:   QCStatusPassed,
				ExpiryDate: time.Now().AddDate(0, 0, -1),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.lot.CanBeIssued() != tt.expected {
				t.Errorf("CanBeIssued() = %v, expected %v", tt.lot.CanBeIssued(), tt.expected)
			}
		})
	}
}
