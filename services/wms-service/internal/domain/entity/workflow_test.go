package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGRN_Workflow(t *testing.T) {
	grn := &GRN{
		ID:        uuid.New(),
		GRNNumber: "GRN-2026-0001",
		Status:    GRNStatusDraft,
		QCStatus:  QCStatusPending,
	}

	// Test CanComplete
	if !grn.CanComplete() {
		t.Error("CanComplete() = false for DRAFT GRN, expected true")
	}

	// Complete with QC Pass
	grn.Complete(QCStatusPassed, "All items meet quality standards")
	if grn.Status != GRNStatusCompleted {
		t.Errorf("After Complete(), Status = %s, expected %s", grn.Status, GRNStatusCompleted)
	}
	if grn.QCStatus != QCStatusPassed {
		t.Errorf("After Complete(), QCStatus = %s, expected %s", grn.QCStatus, QCStatusPassed)
	}
	if grn.CompletedAt == nil {
		t.Error("After Complete(), CompletedAt should not be nil")
	}

	// Test CanComplete after already completed
	if grn.CanComplete() {
		t.Error("CanComplete() = true for COMPLETED GRN, expected false")
	}
}

func TestGRN_CompleteWithQCFail(t *testing.T) {
	grn := &GRN{
		ID:       uuid.New(),
		Status:   GRNStatusDraft,
		QCStatus: QCStatusPending,
	}

	grn.Complete(QCStatusFailed, "Materials do not meet specifications")
	if grn.Status != GRNStatusCompleted {
		t.Errorf("Status = %s, expected %s", grn.Status, GRNStatusCompleted)
	}
	if grn.QCStatus != QCStatusFailed {
		t.Errorf("QCStatus = %s, expected %s", grn.QCStatus, QCStatusFailed)
	}
	if grn.QCNotes != "Materials do not meet specifications" {
		t.Errorf("QCNotes = %s, expected notes about failure", grn.QCNotes)
	}
}

func TestGRN_Cancel(t *testing.T) {
	grn := &GRN{
		ID:     uuid.New(),
		Status: GRNStatusDraft,
	}

	grn.Cancel()
	if grn.Status != GRNStatusCancelled {
		t.Errorf("After Cancel(), Status = %s, expected %s", grn.Status, GRNStatusCancelled)
	}
}

func TestGRNLineItem_PassQC(t *testing.T) {
	lineItem := &GRNLineItem{
		ID:          uuid.New(),
		ReceivedQty: 100,
		QCStatus:    QCStatusPending,
	}

	lineItem.PassQC(100)
	if lineItem.QCStatus != QCStatusPassed {
		t.Errorf("After PassQC(), QCStatus = %s, expected %s", lineItem.QCStatus, QCStatusPassed)
	}
	if lineItem.AcceptedQty == nil || *lineItem.AcceptedQty != 100 {
		t.Error("After PassQC(100), AcceptedQty should be 100")
	}
}

func TestGRNLineItem_FailQC(t *testing.T) {
	lineItem := &GRNLineItem{
		ID:          uuid.New(),
		ReceivedQty: 100,
		QCStatus:    QCStatusPending,
	}

	lineItem.FailQC("Contamination detected")
	if lineItem.QCStatus != QCStatusFailed {
		t.Errorf("After FailQC(), QCStatus = %s, expected %s", lineItem.QCStatus, QCStatusFailed)
	}
	if lineItem.RejectedQty != 100 {
		t.Errorf("After FailQC(), RejectedQty = %f, expected 100", lineItem.RejectedQty)
	}
}

func TestGoodsIssue_Workflow(t *testing.T) {
	issue := &GoodsIssue{
		ID:     uuid.New(),
		Status: GoodsIssueStatusDraft,
	}

	// Confirm
	issue.Confirm()
	if issue.Status != GoodsIssueStatusConfirmed {
		t.Errorf("After Confirm(), Status = %s, expected %s", issue.Status, GoodsIssueStatusConfirmed)
	}

	// Complete
	issue.Complete()
	if issue.Status != GoodsIssueStatusCompleted {
		t.Errorf("After Complete(), Status = %s, expected %s", issue.Status, GoodsIssueStatusCompleted)
	}
}

func TestInventoryCount_Workflow(t *testing.T) {
	count := &InventoryCount{
		ID:          uuid.New(),
		CountNumber: "IC-2026-0001",
		Status:      InventoryCountStatusDraft,
	}

	// Test CanStart
	if !count.CanStart() {
		t.Error("CanStart() = false for DRAFT count, expected true")
	}

	// Start
	count.Start()
	if count.Status != InventoryCountStatusInProgress {
		t.Errorf("After Start(), Status = %s, expected %s", count.Status, InventoryCountStatusInProgress)
	}
	if count.StartedAt == nil {
		t.Error("After Start(), StartedAt should not be nil")
	}

	// Test CanComplete
	if !count.CanComplete() {
		t.Error("CanComplete() = false for IN_PROGRESS count, expected true")
	}

	// Complete
	approvedBy := uuid.New()
	count.Complete(approvedBy)
	if count.Status != InventoryCountStatusCompleted {
		t.Errorf("After Complete(), Status = %s, expected %s", count.Status, InventoryCountStatusCompleted)
	}
	if count.CompletedAt == nil {
		t.Error("After Complete(), CompletedAt should not be nil")
	}
}

func TestInventoryCountLineItem_RecordCount(t *testing.T) {
	lineItem := &InventoryCountLineItem{
		ID:        uuid.New(),
		SystemQty: 100,
		IsCounted: false,
	}

	countedBy := uuid.New()

	// Record count with variance
	lineItem.RecordCount(95, countedBy)
	if !lineItem.IsCounted {
		t.Error("After RecordCount(), IsCounted should be true")
	}
	if lineItem.CountedQty == nil || *lineItem.CountedQty != 95 {
		t.Error("After RecordCount(95), CountedQty should be 95")
	}
	if lineItem.Variance != -5 {
		t.Errorf("Variance = %f, expected -5", lineItem.Variance)
	}
	if lineItem.VariancePercent != -5 {
		t.Errorf("VariancePercent = %f, expected -5%%", lineItem.VariancePercent)
	}
	if !lineItem.HasVariance() {
		t.Error("HasVariance() = false, expected true")
	}
	if lineItem.CountedAt == nil {
		t.Error("CountedAt should not be nil")
	}
}

func TestInventoryCountLineItem_NoVariance(t *testing.T) {
	lineItem := &InventoryCountLineItem{
		ID:        uuid.New(),
		SystemQty: 100,
		IsCounted: false,
	}

	countedBy := uuid.New()
	lineItem.RecordCount(100, countedBy)
	
	if lineItem.Variance != 0 {
		t.Errorf("Variance = %f, expected 0", lineItem.Variance)
	}
	if lineItem.HasVariance() {
		t.Error("HasVariance() = true, expected false")
	}
}

func TestStockMovement_NewMovements(t *testing.T) {
	materialID := uuid.New()
	lotID := uuid.New()
	locationID := uuid.New()
	unitID := uuid.New()
	createdBy := uuid.New()

	// Test IN movement
	movement := NewStockMovementIn(materialID, lotID, locationID, unitID, createdBy, 100, ReferenceTypeGRN, nil, "MOV-IN-2026-0001")
	if movement.MovementType != MovementTypeIn {
		t.Errorf("MovementType = %s, expected %s", movement.MovementType, MovementTypeIn)
	}
	if movement.Quantity != 100 {
		t.Errorf("Quantity = %f, expected 100", movement.Quantity)
	}
	if movement.ToLocationID == nil {
		t.Error("ToLocationID should not be nil for IN movement")
	}

	// Test OUT movement
	movement = NewStockMovementOut(materialID, &lotID, &locationID, unitID, createdBy, 50, ReferenceTypeGI, nil, "MOV-OUT-2026-0001")
	if movement.MovementType != MovementTypeOut {
		t.Errorf("MovementType = %s, expected %s", movement.MovementType, MovementTypeOut)
	}
	if movement.FromLocationID == nil {
		t.Error("FromLocationID should not be nil for OUT movement")
	}

	// Test TRANSFER movement
	toLocationID := uuid.New()
	movement = NewStockMovementTransfer(materialID, &lotID, locationID, toLocationID, unitID, createdBy, 25, "MOV-TRF-2026-0001")
	if movement.MovementType != MovementTypeTransfer {
		t.Errorf("MovementType = %s, expected %s", movement.MovementType, MovementTypeTransfer)
	}
	if movement.FromLocationID == nil || movement.ToLocationID == nil {
		t.Error("Both FromLocationID and ToLocationID should not be nil for TRANSFER movement")
	}
}

func TestStockReservation_Workflow(t *testing.T) {
	reservation := &StockReservation{
		ID:       uuid.New(),
		Quantity: 50,
		Status:   ReservationStatusActive,
	}

	// Test IsActive
	if !reservation.IsActive() {
		t.Error("IsActive() = false for ACTIVE reservation, expected true")
	}

	// Release
	reservation.Release()
	if reservation.Status != ReservationStatusReleased {
		t.Errorf("After Release(), Status = %s, expected %s", reservation.Status, ReservationStatusReleased)
	}

	// Test with new reservation for Fulfill
	reservation2 := &StockReservation{
		ID:       uuid.New(),
		Quantity: 50,
		Status:   ReservationStatusActive,
	}

	reservation2.Fulfill()
	if reservation2.Status != ReservationStatusFulfilled {
		t.Errorf("After Fulfill(), Status = %s, expected %s", reservation2.Status, ReservationStatusFulfilled)
	}
}

func TestStockReservation_IsExpired(t *testing.T) {
	pastTime := time.Now().Add(-24 * time.Hour)
	futureTime := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name      string
		expiresAt *time.Time
		expected  bool
	}{
		{
			name:      "No expiry",
			expiresAt: nil,
			expected:  false,
		},
		{
			name:      "Expired",
			expiresAt: &pastTime,
			expected:  true,
		},
		{
			name:      "Not expired",
			expiresAt: &futureTime,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reservation := &StockReservation{
				Status:    ReservationStatusActive,
				ExpiresAt: tt.expiresAt,
			}
			if reservation.IsExpired() != tt.expected {
				t.Errorf("IsExpired() = %v, expected %v", reservation.IsExpired(), tt.expected)
			}
		})
	}
}
