package evaluation

import (
	"context"
	"time"

	"github.com/erp-cosmetics/supplier-service/internal/domain/entity"
	"github.com/erp-cosmetics/supplier-service/internal/domain/repository"
	"github.com/erp-cosmetics/supplier-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateEvaluationRequest represents the request to create an evaluation
type CreateEvaluationRequest struct {
	SupplierID            uuid.UUID `json:"-"`
	EvaluationDate        string    `json:"evaluation_date" validate:"required"`
	EvaluationPeriod      string    `json:"evaluation_period" validate:"required"`
	QualityScore          float64   `json:"quality_score" validate:"required,min=1,max=5"`
	DeliveryScore         float64   `json:"delivery_score" validate:"required,min=1,max=5"`
	PriceScore            float64   `json:"price_score" validate:"required,min=1,max=5"`
	ServiceScore          float64   `json:"service_score" validate:"required,min=1,max=5"`
	DocumentationScore    float64   `json:"documentation_score" validate:"required,min=1,max=5"`
	OnTimeDeliveryRate    float64   `json:"on_time_delivery_rate"`
	QualityAcceptanceRate float64   `json:"quality_acceptance_rate"`
	LeadTimeAdherence     float64   `json:"lead_time_adherence"`
	Strengths             string    `json:"strengths"`
	Weaknesses            string    `json:"weaknesses"`
	ActionItems           string    `json:"action_items"`
	EvaluatedBy           uuid.UUID `json:"-"`
}

// CreateEvaluationUseCase handles creating a new evaluation
type CreateEvaluationUseCase struct {
	evalRepo     repository.EvaluationRepository
	supplierRepo repository.SupplierRepository
	eventPub     *event.Publisher
}

// NewCreateEvaluationUseCase creates a new CreateEvaluationUseCase
func NewCreateEvaluationUseCase(
	evalRepo repository.EvaluationRepository,
	supplierRepo repository.SupplierRepository,
	eventPub *event.Publisher,
) *CreateEvaluationUseCase {
	return &CreateEvaluationUseCase{
		evalRepo:     evalRepo,
		supplierRepo: supplierRepo,
		eventPub:     eventPub,
	}
}

// Execute creates a new evaluation
func (uc *CreateEvaluationUseCase) Execute(ctx context.Context, req *CreateEvaluationRequest) (*entity.Evaluation, error) {
	// Validate supplier exists
	supplier, err := uc.supplierRepo.GetByID(ctx, req.SupplierID)
	if err != nil {
		return nil, err
	}

	// Parse date
	evalDate, err := time.Parse("2006-01-02", req.EvaluationDate)
	if err != nil {
		evalDate = time.Now()
	}

	eval := &entity.Evaluation{
		ID:                    uuid.New(),
		SupplierID:            req.SupplierID,
		EvaluationDate:        evalDate,
		EvaluationPeriod:      req.EvaluationPeriod,
		QualityScore:          req.QualityScore,
		DeliveryScore:         req.DeliveryScore,
		PriceScore:            req.PriceScore,
		ServiceScore:          req.ServiceScore,
		DocumentationScore:    req.DocumentationScore,
		OnTimeDeliveryRate:    req.OnTimeDeliveryRate,
		QualityAcceptanceRate: req.QualityAcceptanceRate,
		LeadTimeAdherence:     req.LeadTimeAdherence,
		Strengths:             req.Strengths,
		Weaknesses:            req.Weaknesses,
		ActionItems:           req.ActionItems,
		EvaluatedBy:           req.EvaluatedBy,
		Status:                entity.EvaluationStatusSubmitted,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	// Calculate overall score
	eval.CalculateOverallScore()

	if err := uc.evalRepo.Create(ctx, eval); err != nil {
		return nil, err
	}

	// Update supplier ratings
	quality, delivery, service, err := uc.evalRepo.GetAverageScoresBySupplierID(ctx, req.SupplierID)
	if err == nil {
		overall := (quality + delivery + service) / 3
		uc.supplierRepo.UpdateRating(ctx, req.SupplierID, quality, delivery, service, overall)
	}

	// Publish event
	uc.eventPub.PublishEvaluationCompleted(ctx, &event.EvaluationEvent{
		SupplierID:       supplier.ID.String(),
		EvaluationID:     eval.ID.String(),
		EvaluationPeriod: eval.EvaluationPeriod,
		OverallScore:     eval.OverallScore,
		EvaluatedBy:      eval.EvaluatedBy.String(),
		EvaluatedAt:      eval.EvaluationDate,
	})

	return eval, nil
}

// GetEvaluationsUseCase handles getting evaluations for a supplier
type GetEvaluationsUseCase struct {
	evalRepo repository.EvaluationRepository
}

// NewGetEvaluationsUseCase creates a new GetEvaluationsUseCase
func NewGetEvaluationsUseCase(evalRepo repository.EvaluationRepository) *GetEvaluationsUseCase {
	return &GetEvaluationsUseCase{evalRepo: evalRepo}
}

// Execute gets evaluations for a supplier
func (uc *GetEvaluationsUseCase) Execute(ctx context.Context, supplierID uuid.UUID) ([]*entity.Evaluation, error) {
	return uc.evalRepo.GetBySupplierID(ctx, supplierID)
}
