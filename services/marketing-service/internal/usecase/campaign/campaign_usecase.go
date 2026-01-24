package campaign

import (
	"context"
	"time"

	"github.com/erp-cosmetics/marketing-service/internal/domain/entity"
	"github.com/erp-cosmetics/marketing-service/internal/domain/repository"
	"github.com/erp-cosmetics/marketing-service/internal/infrastructure/event"
	"github.com/google/uuid"
)

// CreateCampaignInput represents input for creating campaign
type CreateCampaignInput struct {
	Name           string
	Description    string
	CampaignType   entity.CampaignType
	StartDate      time.Time
	EndDate        time.Time
	TargetAudience string
	Channels       []byte
	Budget         float64
	Currency       string
	Goals          []byte
	Products       []byte
	Notes          string
}

// CreateCampaignUseCase handles campaign creation
type CreateCampaignUseCase struct {
	repo      repository.CampaignRepository
	publisher *event.Publisher
}

// NewCreateCampaignUseCase creates a new use case
func NewCreateCampaignUseCase(repo repository.CampaignRepository, publisher *event.Publisher) *CreateCampaignUseCase {
	return &CreateCampaignUseCase{repo: repo, publisher: publisher}
}

// Execute creates a new campaign
func (uc *CreateCampaignUseCase) Execute(ctx context.Context, input *CreateCampaignInput) (*entity.Campaign, error) {
	code, err := uc.repo.GenerateCampaignCode(ctx, "GEN")
	if err != nil {
		return nil, err
	}

	campaign := &entity.Campaign{
		CampaignCode:   code,
		Name:           input.Name,
		Description:    input.Description,
		CampaignType:   input.CampaignType,
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		TargetAudience: input.TargetAudience,
		Channels:       input.Channels,
		Budget:         input.Budget,
		Currency:       input.Currency,
		Goals:          input.Goals,
		Products:       input.Products,
		Notes:          input.Notes,
		Status:         entity.CampaignStatusDraft,
	}

	if err := uc.repo.Create(ctx, campaign); err != nil {
		return nil, err
	}

	if uc.publisher != nil {
		uc.publisher.PublishCampaignCreated(ctx, map[string]interface{}{
			"campaign_id":   campaign.ID,
			"campaign_code": campaign.CampaignCode,
			"name":          campaign.Name,
			"type":          campaign.CampaignType,
		})
	}

	return campaign, nil
}

// GetCampaignUseCase handles getting a campaign
type GetCampaignUseCase struct {
	repo repository.CampaignRepository
}

// NewGetCampaignUseCase creates a new use case
func NewGetCampaignUseCase(repo repository.CampaignRepository) *GetCampaignUseCase {
	return &GetCampaignUseCase{repo: repo}
}

// Execute gets a campaign by ID
func (uc *GetCampaignUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Campaign, error) {
	return uc.repo.GetByID(ctx, id)
}

// ListCampaignsUseCase handles listing campaigns
type ListCampaignsUseCase struct {
	repo repository.CampaignRepository
}

// NewListCampaignsUseCase creates a new use case
func NewListCampaignsUseCase(repo repository.CampaignRepository) *ListCampaignsUseCase {
	return &ListCampaignsUseCase{repo: repo}
}

// Execute lists campaigns with filter
func (uc *ListCampaignsUseCase) Execute(ctx context.Context, filter *repository.CampaignFilter) ([]*entity.Campaign, int64, error) {
	return uc.repo.List(ctx, filter)
}

// LaunchCampaignUseCase handles launching a campaign
type LaunchCampaignUseCase struct {
	repo      repository.CampaignRepository
	publisher *event.Publisher
}

// NewLaunchCampaignUseCase creates a new use case
func NewLaunchCampaignUseCase(repo repository.CampaignRepository, publisher *event.Publisher) *LaunchCampaignUseCase {
	return &LaunchCampaignUseCase{repo: repo, publisher: publisher}
}

// Execute launches a campaign
func (uc *LaunchCampaignUseCase) Execute(ctx context.Context, id uuid.UUID) (*entity.Campaign, error) {
	campaign, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !campaign.CanBeLaunched() {
		return nil, ErrCampaignCannotBeLaunched
	}

	campaign.Launch()

	if err := uc.repo.Update(ctx, campaign); err != nil {
		return nil, err
	}

	if uc.publisher != nil {
		uc.publisher.PublishCampaignLaunched(ctx, map[string]interface{}{
			"campaign_id":   campaign.ID,
			"campaign_code": campaign.CampaignCode,
			"name":          campaign.Name,
			"launched_at":   campaign.LaunchedAt,
		})
	}

	return campaign, nil
}

// UpdateCampaignInput represents input for updating campaign
type UpdateCampaignInput struct {
	Name           string
	Description    string
	TargetAudience string
	Budget         float64
	Goals          []byte
	Notes          string
}

// UpdateCampaignUseCase handles updating a campaign
type UpdateCampaignUseCase struct {
	repo repository.CampaignRepository
}

// NewUpdateCampaignUseCase creates a new use case
func NewUpdateCampaignUseCase(repo repository.CampaignRepository) *UpdateCampaignUseCase {
	return &UpdateCampaignUseCase{repo: repo}
}

// Execute updates a campaign
func (uc *UpdateCampaignUseCase) Execute(ctx context.Context, id uuid.UUID, input *UpdateCampaignInput) (*entity.Campaign, error) {
	campaign, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	campaign.Name = input.Name
	campaign.Description = input.Description
	campaign.TargetAudience = input.TargetAudience
	campaign.Budget = input.Budget
	campaign.Goals = input.Goals
	campaign.Notes = input.Notes

	if err := uc.repo.Update(ctx, campaign); err != nil {
		return nil, err
	}

	return campaign, nil
}

// Custom errors
var (
	ErrCampaignCannotBeLaunched = &CampaignError{Message: "campaign cannot be launched"}
)

// CampaignError represents a campaign-related error
type CampaignError struct {
	Message string
}

func (e *CampaignError) Error() string {
	return e.Message
}
