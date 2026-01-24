package grpc

import (
	"context"
	"time"

	"github.com/erp-cosmetics/wms-service/internal/domain/entity"
	"github.com/erp-cosmetics/wms-service/internal/domain/repository"
	"github.com/erp-cosmetics/wms-service/internal/usecase/reservation"
	"github.com/erp-cosmetics/wms-service/internal/usecase/stock"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// WMSServer implements the WMS gRPC service
type WMSServer struct {
	UnimplementedWMSServiceServer
	stockRepo            repository.StockRepository
	lotRepo              repository.LotRepository
	issueStockFEFOUC     *stock.IssueStockFEFOUseCase
	reserveStockUC       *reservation.CreateReservationUseCase
	releaseReservationUC *reservation.ReleaseReservationUseCase
	checkAvailabilityUC  *reservation.CheckAvailabilityUseCase
	logger               *zap.Logger
}

// NewWMSServer creates a new WMS gRPC server
func NewWMSServer(
	stockRepo repository.StockRepository,
	lotRepo repository.LotRepository,
	issueStockFEFOUC *stock.IssueStockFEFOUseCase,
	reserveStockUC *reservation.CreateReservationUseCase,
	releaseReservationUC *reservation.ReleaseReservationUseCase,
	checkAvailabilityUC *reservation.CheckAvailabilityUseCase,
	logger *zap.Logger,
) *WMSServer {
	return &WMSServer{
		stockRepo:            stockRepo,
		lotRepo:              lotRepo,
		issueStockFEFOUC:     issueStockFEFOUC,
		reserveStockUC:       reserveStockUC,
		releaseReservationUC: releaseReservationUC,
		checkAvailabilityUC:  checkAvailabilityUC,
		logger:               logger,
	}
}

// CheckStockAvailability checks if material has sufficient stock
func (s *WMSServer) CheckStockAvailability(ctx context.Context, req *CheckStockRequest) (*StockAvailabilityResponse, error) {
	materialID, err := uuid.Parse(req.MaterialId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid material_id")
	}

	result, err := s.checkAvailabilityUC.Execute(ctx, materialID, req.RequestedQuantity)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &StockAvailabilityResponse{
		MaterialId:        req.MaterialId,
		TotalQuantity:     result.TotalQuantity,
		ReservedQuantity:  result.ReservedQty,
		AvailableQuantity: result.AvailableQty,
		IsAvailable:       result.IsAvailable,
		ShortageQuantity:  result.ShortageQty,
	}, nil
}

// ReserveStock reserves stock for an order
func (s *WMSServer) ReserveStock(ctx context.Context, req *ReserveStockRequest) (*ReserveStockResponse, error) {
	materialID, err := uuid.Parse(req.MaterialId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid material_id")
	}

	unitID, err := uuid.Parse(req.UnitId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid unit_id")
	}

	referenceID, err := uuid.Parse(req.ReferenceId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reference_id")
	}

	input := &reservation.CreateReservationInput{
		MaterialID:      materialID,
		Quantity:        req.Quantity,
		UnitID:          unitID,
		ReservationType: entity.ReservationType(req.ReservationType),
		ReferenceID:     referenceID,
		ReferenceNumber: req.ReferenceNumber,
		CreatedBy:       uuid.New(), // Should come from context
	}

	result, err := s.reserveStockUC.Execute(ctx, input)
	if err != nil {
		if err == entity.ErrInsufficientStock {
			return &ReserveStockResponse{
				Success: false,
				Message: "Insufficient stock for reservation",
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ReserveStockResponse{
		ReservationId:    result.ID.String(),
		ReservedQuantity: result.Quantity,
		Success:          true,
		Message:          "Stock reserved successfully",
	}, nil
}

// ReleaseReservation releases a stock reservation
func (s *WMSServer) ReleaseReservation(ctx context.Context, req *ReleaseReservationRequest) (*ReleaseReservationResponse, error) {
	reservationID, err := uuid.Parse(req.ReservationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid reservation_id")
	}

	if err := s.releaseReservationUC.Execute(ctx, reservationID); err != nil {
		return &ReleaseReservationResponse{Success: false}, nil
	}

	return &ReleaseReservationResponse{Success: true}, nil
}

// IssueStock issues stock using FEFO logic
func (s *WMSServer) IssueStock(ctx context.Context, req *IssueStockRequest) (*IssueStockResponse, error) {
	issuedBy, _ := uuid.Parse(req.IssuedBy)
	referenceID, _ := uuid.Parse(req.ReferenceId)

	var lineItems []*IssueLineItem

	for _, item := range req.Items {
		materialID, err := uuid.Parse(item.MaterialId)
		if err != nil {
			continue
		}

		unitID, err := uuid.Parse(item.UnitId)
		if err != nil {
			continue
		}

		input := &stock.IssueStockInput{
			MaterialID:      materialID,
			Quantity:        item.Quantity,
			UnitID:          unitID,
			ReferenceType:   entity.ReferenceType(req.ReferenceType),
			ReferenceID:     &referenceID,
			ReferenceNumber: req.ReferenceNumber,
			CreatedBy:       issuedBy,
		}

		result, err := s.issueStockFEFOUC.Execute(ctx, input)
		if err != nil {
			if err == entity.ErrInsufficientStock {
				return &IssueStockResponse{
					Success: false,
					Message: "Insufficient stock for material: " + item.MaterialId,
				}, nil
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		// Convert lots issued to proto format
		lotsUsed := make([]*LotIssued, len(result.LotsIssued))
		for i, lot := range result.LotsIssued {
			lotsUsed[i] = &LotIssued{
				LotId:      lot.LotID.String(),
				LotNumber:  lot.LotNumber,
				Quantity:   lot.Quantity,
				ExpiryDate: timestamppb.New(lot.ExpiryDate),
				LocationId: lot.LocationID.String(),
			}
		}

		lineItems = append(lineItems, &IssueLineItem{
			MaterialId:     item.MaterialId,
			IssuedQuantity: item.Quantity,
			LotsUsed:       lotsUsed,
		})
	}

	return &IssueStockResponse{
		IssueNumber: time.Now().Format("GI-20060102-1504"),
		LineItems:   lineItems,
		Success:     true,
		Message:     "Stock issued successfully",
	}, nil
}

// GetLotInfo gets lot information
func (s *WMSServer) GetLotInfo(ctx context.Context, req *GetLotRequest) (*LotInfoResponse, error) {
	lotID, err := uuid.Parse(req.LotId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid lot_id")
	}

	lot, err := s.lotRepo.GetByID(ctx, lotID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "lot not found")
	}

	// Get available quantity from stock
	stocks, _ := s.stockRepo.GetByMaterialAndLot(ctx, lot.MaterialID, lot.ID)
	availableQty := 0.0
	if stocks != nil {
		availableQty = stocks.Quantity - stocks.ReservedQty
	}

	return &LotInfoResponse{
		Lot: &LotInfo{
			LotId:             lot.ID.String(),
			LotNumber:         lot.LotNumber,
			MaterialId:        lot.MaterialID.String(),
			ExpiryDate:        timestamppb.New(lot.ExpiryDate),
			QcStatus:          string(lot.QCStatus),
			Status:            string(lot.Status),
			AvailableQuantity: availableQty,
		},
	}, nil
}

// GetLotsByMaterial gets lots for a material
func (s *WMSServer) GetLotsByMaterial(ctx context.Context, req *GetLotsByMaterialRequest) (*LotsResponse, error) {
	materialID, err := uuid.Parse(req.MaterialId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid material_id")
	}

	var lots []*entity.Lot
	if req.AvailableOnly {
		lots, err = s.lotRepo.GetAvailableLots(ctx, materialID)
	} else {
		filter := &repository.LotFilter{MaterialID: &materialID, Limit: 100}
		lots, _, err = s.lotRepo.List(ctx, filter)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	lotInfos := make([]*LotInfo, len(lots))
	for i, lot := range lots {
		lotInfos[i] = &LotInfo{
			LotId:      lot.ID.String(),
			LotNumber:  lot.LotNumber,
			MaterialId: lot.MaterialID.String(),
			ExpiryDate: timestamppb.New(lot.ExpiryDate),
			QcStatus:   string(lot.QCStatus),
			Status:     string(lot.Status),
		}
	}

	return &LotsResponse{Lots: lotInfos}, nil
}

// ReceiveStock handles receiving stock from procurement
func (s *WMSServer) ReceiveStock(ctx context.Context, req *ReceiveStockRequest) (*ReceiveStockResponse, error) {
	// This would create GRN and receive stock
	// For now, just return success
	return &ReceiveStockResponse{
		GrnNumber: "GRN-" + time.Now().Format("20060102-1504"),
		Success:   true,
	}, nil
}

// UnimplementedWMSServiceServer is embedded to ensure forward compatibility
type UnimplementedWMSServiceServer struct{}

func (UnimplementedWMSServiceServer) CheckStockAvailability(context.Context, *CheckStockRequest) (*StockAvailabilityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckStockAvailability not implemented")
}
func (UnimplementedWMSServiceServer) ReserveStock(context.Context, *ReserveStockRequest) (*ReserveStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReserveStock not implemented")
}
func (UnimplementedWMSServiceServer) ReleaseReservation(context.Context, *ReleaseReservationRequest) (*ReleaseReservationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseReservation not implemented")
}
func (UnimplementedWMSServiceServer) IssueStock(context.Context, *IssueStockRequest) (*IssueStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueStock not implemented")
}
func (UnimplementedWMSServiceServer) GetLotInfo(context.Context, *GetLotRequest) (*LotInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLotInfo not implemented")
}
func (UnimplementedWMSServiceServer) GetLotsByMaterial(context.Context, *GetLotsByMaterialRequest) (*LotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLotsByMaterial not implemented")
}
func (UnimplementedWMSServiceServer) ReceiveStock(context.Context, *ReceiveStockRequest) (*ReceiveStockResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveStock not implemented")
}
