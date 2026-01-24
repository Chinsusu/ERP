package kol

import (
	"context"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	"github.com/erp-cosmetics/marketing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateKOLInput represents input for creating KOL
type CreateKOLInput struct {
	Name               string
	Email              string
	Phone              string
	TierID             *uuid.UUID
	Category           entity.KOLCategory
	InstagramHandle    string
	InstagramFollowers int
	YouTubeChannel     string
	YouTubeSubscribers int
	TikTokHandle       string
	TikTokFollowers    int
	FacebookPage       string
	FacebookFollowers  int
	AvgEngagementRate  float64
	Niche              string
	CollaborationRate  float64
	Currency           string
	AddressLine1       string
	City               string
	Notes              string
}

// CreateKOLUseCase handles KOL creation
type CreateKOLUseCase struct {
	repo      repository.KOLRepository
	publisher *event.Publisher
}

// NewCreateKOLUseCase creates a new use case
func NewCreateKOLUseCase(repo repository.KOLRepository, publisher *event.Publisher) *CreateKOLUseCase {
	return &CreateKOLUseCase{repo: repo, publisher: publisher}
}

// Execute creates a new KOL
func (uc *CreateKOLUseCase) Execute(ctx context.Context, input *CreateKOLInput) (*entity.KOL, error) {
	code, err := uc.repo.GenerateKOLCode(ctx)
	if err != nil {
		return nil, err
	}

	kol := &entity.KOL{
		KOLCode:            code,
		Name:               input.Name,
		Email:              input.Email,
		Phone:              input.Phone,
		TierID:             input.TierID,
		Category:           input.Category,
		InstagramHandle:    input.InstagramHandle,
		InstagramFollowers: input.InstagramFollowers,
		YouTubeChannel:     input.YouTubeChannel,
		YouTubeSubscribers: input.YouTubeSubscribers,
		TikTokHandle:       input.TikTokHandle,
		TikTokFollowers:    input.TikTokFollowers,
		FacebookPage:       input.FacebookPage,
		FacebookFollowers:  input.FacebookFollowers,
		AvgEngagementRate:  input.AvgEngagementRate,
		Niche:              input.Niche,
		CollaborationRate:  input.CollaborationRate,
		Currency:           input.Currency,
		AddressLine1:       input.AddressLine1,
		City:               input.City,
		Notes:              input.Notes,
		Status:             entity.KOLStatusActive,
	}

	if err := uc.repo.Create(ctx, kol); err != nil {
		return nil, err
	}

	if uc.publisher != nil {
		uc.publisher.PublishKOLCreated(ctx, map[string]interface{}{
			"kol_id":   kol.ID,
			"kol_code": kol.KOLCode,
			"name":     kol.Name,
		})
	}

	return kol, nil
}

// GetKOLUseCase handles getting a KOL
type GetKOLUseCase struct {
	repo repository.KOLRepository
}

// NewGetKOLUseCase creates a new use case
func NewGetKOLUseCase(repo repository.KOLRepository) *GetKOLUseCase {
	return &GetKOLUseCase{repo: repo}
}

// Execute gets a KOL by ID
func (uc *GetKOLUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.KOL, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListKOLsUseCase handles listing KOLs
type ListKOLsUseCase struct {
	repo repository.KOLRepository
}

// NewListKOLsUseCase creates a new use case
func NewListKOLsUseCase(repo repository.KOLRepository) *ListKOLsUseCase {
	return &ListKOLsUseCase{repo: repo}
}

// Execute lists KOLs with filter
func (uc *ListKOLsUseCase) Execute(ctx context.Context, filter *repository.KOLFilter) ([]*entity.KOL, int64, error) {
	return uc.repo.List(ctx, filter)
}

// UpdateKOLInput represents input for updating KOL
type UpdateKOLInput struct {
	Name               string
	Email              string
	Phone              string
	TierID             *uuid.UUID
	Category           entity.KOLCategory
	InstagramHandle    string
	InstagramFollowers int
	YouTubeChannel     string
	YouTubeSubscribers int
	TikTokHandle       string
	TikTokFollowers    int
	AvgEngagementRate  float64
	Niche              string
	CollaborationRate  float64
	AddressLine1       string
	City               string
	Status             entity.KOLStatus
	Notes              string
}

// UpdateKOLUseCase handles updating a KOL
type UpdateKOLUseCase struct {
	repo repository.KOLRepository
}

// NewUpdateKOLUseCase creates a new use case
func NewUpdateKOLUseCase(repo repository.KOLRepository) *UpdateKOLUseCase {
	return &UpdateKOLUseCase{repo: repo}
}

// Execute updates a KOL
func (uc *UpdateKOLUseCase) Execute(ctx context.Context, id uuid.UUID, input *UpdateKOLInput) (*entity.KOL, error) {
	kol, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	kol.Name = input.Name
	kol.Email = input.Email
	kol.Phone = input.Phone
	kol.TierID = input.TierID
	kol.Category = input.Category
	kol.InstagramHandle = input.InstagramHandle
	kol.InstagramFollowers = input.InstagramFollowers
	kol.YouTubeChannel = input.YouTubeChannel
	kol.YouTubeSubscribers = input.YouTubeSubscribers
	kol.TikTokHandle = input.TikTokHandle
	kol.TikTokFollowers = input.TikTokFollowers
	kol.AvgEngagementRate = input.AvgEngagementRate
	kol.Niche = input.Niche
	kol.CollaborationRate = input.CollaborationRate
	kol.AddressLine1 = input.AddressLine1
	kol.City = input.City
	kol.Notes = input.Notes
	if input.Status != "" {
		kol.Status = input.Status
	}

	if err := uc.repo.Update(ctx, kol); err != nil {
		return nil, err
	}

	return kol, nil
}

// DeleteKOLUseCase handles deleting a KOL
type DeleteKOLUseCase struct {
	repo repository.KOLRepository
}

// NewDeleteKOLUseCase creates a new use case
func NewDeleteKOLUseCase(repo repository.KOLRepository) *DeleteKOLUseCase {
	return &DeleteKOLUseCase{repo: repo}
}

// Execute soft deletes a KOL (sets to inactive)
func (uc *DeleteKOLUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	kol, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	kol.Status = entity.KOLStatusInactive
	return uc.repo.Update(ctx, kol)
}
