package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestStock_GetAvailableQuantity(t *testing.T) {
	tests := []struct {
		name        string
		quantity    float64
		reservedQty float64
		expected    float64
	}{
		{
			name:        "No reservation",
			quantity:    100,
			reservedQty: 0,
			expected:    100,
		},
		{
			name:        "With reservation",
			quantity:    100,
			reservedQty: 30,
			expected:    70,
		},
		{
			name:        "Fully reserved",
			quantity:    100,
			reservedQty: 100,
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := &Stock{
				Quantity:    tt.quantity,
				ReservedQty: tt.reservedQty,
			}
			if stock.GetAvailableQuantity() != tt.expected {
				t.Errorf("GetAvailableQuantity() = %f, expected %f", 
					stock.GetAvailableQuantity(), tt.expected)
			}
		})
	}
}

func TestStock_CanIssue(t *testing.T) {
	tests := []struct {
		name        string
		quantity    float64
		reservedQty float64
		issueQty    float64
		expected    bool
	}{
		{
			name:        "Sufficient stock",
			quantity:    100,
			reservedQty: 0,
			issueQty:    50,
			expected:    true,
		},
		{
			name:        "Exactly available",
			quantity:    100,
			reservedQty: 50,
			issueQty:    50,
			expected:    true,
		},
		{
			name:        "Insufficient (over reserved)",
			quantity:    100,
			reservedQty: 60,
			issueQty:    50,
			expected:    false,
		},
		{
			name:        "Insufficient (over quantity)",
			quantity:    100,
			reservedQty: 0,
			issueQty:    150,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := &Stock{
				Quantity:    tt.quantity,
				ReservedQty: tt.reservedQty,
			}
			if stock.CanIssue(tt.issueQty) != tt.expected {
				t.Errorf("CanIssue(%f) = %v, expected %v", 
					tt.issueQty, stock.CanIssue(tt.issueQty), tt.expected)
			}
		})
	}
}

func TestStock_Reserve(t *testing.T) {
	stock := &Stock{
		Quantity:    100,
		ReservedQty: 0,
	}

	// Reserve 30
	err := stock.Reserve(30)
	if err != nil {
		t.Errorf("Reserve(30) error = %v, expected nil", err)
	}
	if stock.ReservedQty != 30 {
		t.Errorf("After Reserve(30), ReservedQty = %f, expected 30", stock.ReservedQty)
	}

	// Reserve another 50
	err = stock.Reserve(50)
	if err != nil {
		t.Errorf("Reserve(50) error = %v, expected nil", err)
	}
	if stock.ReservedQty != 80 {
		t.Errorf("After Reserve(50), ReservedQty = %f, expected 80", stock.ReservedQty)
	}

	// Try to reserve more than available
	err = stock.Reserve(30) // Only 20 available
	if err == nil {
		t.Error("Reserve(30) should have returned error for insufficient stock")
	}
}

func TestStock_ReleaseReservation(t *testing.T) {
	stock := &Stock{
		Quantity:    100,
		ReservedQty: 50,
	}

	// Release 20
	stock.ReleaseReservation(20)
	if stock.ReservedQty != 30 {
		t.Errorf("After ReleaseReservation(20), ReservedQty = %f, expected 30", stock.ReservedQty)
	}

	// Release more than reserved
	stock.ReleaseReservation(50)
	if stock.ReservedQty != 0 {
		t.Errorf("After ReleaseReservation(50), ReservedQty = %f, expected 0", stock.ReservedQty)
	}
}

func TestStock_Issue(t *testing.T) {
	stock := &Stock{
		ID:          uuid.New(),
		Quantity:    100,
		ReservedQty: 20,
	}

	// Issue 30
	err := stock.Issue(30)
	if err != nil {
		t.Errorf("Issue(30) error = %v, expected nil", err)
	}
	if stock.Quantity != 70 {
		t.Errorf("After Issue(30), Quantity = %f, expected 70", stock.Quantity)
	}

	// Try to issue more than total quantity
	err = stock.Issue(80) // Only 70 quantity left
	if err == nil {
		t.Error("Issue(80) should have returned error for insufficient stock")
	}
}

func TestStock_Receive(t *testing.T) {
	stock := &Stock{
		Quantity: 50,
	}

	stock.Receive(30)
	if stock.Quantity != 80 {
		t.Errorf("After Receive(30), Quantity = %f, expected 80", stock.Quantity)
	}

	stock.Receive(20)
	if stock.Quantity != 100 {
		t.Errorf("After Receive(20), Quantity = %f, expected 100", stock.Quantity)
	}
}
