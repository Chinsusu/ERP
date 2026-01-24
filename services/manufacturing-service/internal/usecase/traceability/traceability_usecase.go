package traceability

import (
	"context"

	"github.com/erp-cosmetics/manufacturing-service/internal/domain/entity"
	"github.com/erp-cosmetics/manufacturing-service/internal/domain/repository"
	"github.com/google/uuid"
)

// TraceBackwardUseCase handles backward tracing (product lot → material lots)
type TraceBackwardUseCase struct {
	repo   repository.TraceabilityRepository
	woRepo repository.WorkOrderRepository
}

// NewTraceBackwardUseCase creates a new TraceBackwardUseCase
func NewTraceBackwardUseCase(repo repository.TraceabilityRepository, woRepo repository.WorkOrderRepository) *TraceBackwardUseCase {
	return &TraceBackwardUseCase{
		repo:   repo,
		woRepo: woRepo,
	}
}

// Execute traces backward from a product lot to material lots
func (uc *TraceBackwardUseCase) Execute(ctx context.Context, productLotID uuid.UUID) (*entity.BackwardTraceResult, error) {
	traces, err := uc.repo.GetByProductLot(ctx, productLotID)
	if err != nil || len(traces) == 0 {
		return nil, entity.ErrTraceNotFound
	}

	// Get work order info
	woID := traces[0].WorkOrderID
	wo, err := uc.woRepo.GetByID(ctx, woID)
	if err != nil {
		return nil, err
	}

	result := &entity.BackwardTraceResult{
		FinishedLot: entity.FinishedLotInfo{
			LotNumber:        traces[0].ProductLotNumber,
			ProductCode:      "", // Would need to join with master data
			ProductName:      "", // Would need to join with master data
			Quantity:         0,
			ManufacturedDate: wo.CreatedAt,
		},
		WorkOrder: entity.WorkOrderInfo{
			WONumber:   wo.WONumber,
			Supervisor: "", // Would need to join with user service
		},
		MaterialsUsed: []entity.MaterialTraceInfo{},
	}

	if wo.GoodQuantity != nil {
		result.FinishedLot.Quantity = *wo.GoodQuantity
	}
	if wo.ActualEndDate != nil {
		result.FinishedLot.ManufacturedDate = *wo.ActualEndDate
	}

	for _, trace := range traces {
		result.MaterialsUsed = append(result.MaterialsUsed, entity.MaterialTraceInfo{
			MaterialCode: "", // Would need to join with master data
			MaterialName: "", // Would need to join with master data
			LotNumber:    trace.MaterialLotNumber,
			Quantity:     trace.MaterialQuantity,
			UOM:          "", // Would need to join with master data
			Supplier:     "", // Would need to join with supplier service
			SupplierLot:  trace.SupplierLotNumber,
		})
	}

	return result, nil
}

// TraceForwardUseCase handles forward tracing (material lot → product lots)
type TraceForwardUseCase struct {
	repo repository.TraceabilityRepository
}

// NewTraceForwardUseCase creates a new TraceForwardUseCase
func NewTraceForwardUseCase(repo repository.TraceabilityRepository) *TraceForwardUseCase {
	return &TraceForwardUseCase{repo: repo}
}

// Execute traces forward from a material lot to product lots
func (uc *TraceForwardUseCase) Execute(ctx context.Context, materialLotID uuid.UUID) (*entity.ForwardTraceResult, error) {
	traces, err := uc.repo.GetByMaterialLot(ctx, materialLotID)
	if err != nil || len(traces) == 0 {
		return nil, entity.ErrTraceNotFound
	}

	result := &entity.ForwardTraceResult{
		MaterialLot: entity.MaterialLotInfo{
			LotNumber:    traces[0].MaterialLotNumber,
			MaterialCode: "", // Would need to join with master data
			MaterialName: "", // Would need to join with master data
			Supplier:     "", // Would need to join with supplier service
		},
		UsedInProducts: []entity.ProductTraceInfo{},
		TotalQtyUsed:   0,
	}

	for _, trace := range traces {
		result.UsedInProducts = append(result.UsedInProducts, entity.ProductTraceInfo{
			ProductLot:     trace.ProductLotNumber,
			ProductName:    "", // Would need to join with master data
			WONumber:       "", // Would need to join with WO
			QuantityUsed:   trace.MaterialQuantity,
			ProductionDate: trace.TraceDate,
		})
		result.TotalQtyUsed += trace.MaterialQuantity
	}

	return result, nil
}
