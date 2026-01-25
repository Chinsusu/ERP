package qc_test

import (
	"context"
	"testing"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/testmocks"
	"github.com/erp-cosmetics/manufacturing-service/internal/usecase/qc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateInspectionUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	qcRepo := new(testmocks.MockQCRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := qc.NewCreateInspectionUseCase(qcRepo, eventPub)

	woID := uuid.New()
	input := qc.CreateInspectionInput{
		InspectionType: entity.CheckpointTypeIPQC,
		ReferenceType:  entity.ReferenceTypeWorkOrder,
		ReferenceID:    woID,
		InspectorID:    uuid.New(),
	}

	qcRepo.On("GenerateInspectionNumber", ctx).Return("QC-2026-0001", nil)
	qcRepo.On("GetCheckpointsByType", ctx, entity.CheckpointTypeIPQC).Return([]*entity.QCCheckpoint{{ID: uuid.New()}}, nil)
	qcRepo.On("CreateInspection", ctx, mock.AnythingOfType("*entity.QCInspection")).Return(nil)
	qcRepo.On("CreateInspectionItems", ctx, mock.Anything).Return(nil)

	// Act
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, entity.InspectionResultPending, res.Result)
}

func TestApproveInspectionUseCase_Execute_FailureTriggersQCFailedEvent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	qcRepo := new(testmocks.MockQCRepository)
	eventPub := new(testmocks.MockEventPublisher)
	
	uc := qc.NewApproveInspectionUseCase(qcRepo, eventPub)

	inspection := &entity.QCInspection{
		ID:            uuid.New(),
		ReferenceType: entity.ReferenceTypeWorkOrder,
		ReferenceID:   uuid.New(),
		Result:        entity.InspectionResultPending,
	}

	qcRepo.On("GetInspectionByID", ctx, inspection.ID).Return(inspection, nil)
	qcRepo.On("UpdateInspection", ctx, mock.Anything).Return(nil)
	
	eventPub.On("PublishQCFailed", mock.Anything).Return(nil)

	// Act
	input := qc.ApproveInspectionInput{
		InspectionID: inspection.ID,
		Result:       entity.InspectionResultFailed,
		ApproverID:   uuid.New(),
		Notes:        "Defect detected",
	}
	res, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, entity.InspectionResultFailed, res.Result)
	
	eventPub.AssertCalled(t, "PublishQCFailed", mock.Anything)
}
