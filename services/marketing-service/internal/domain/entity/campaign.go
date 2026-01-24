package entity

import (
	"time"

	"github.com/google/uuid"
)

// CampaignType represents the type of marketing campaign
type CampaignType string

const (
	CampaignTypeProductLaunch CampaignType = "PRODUCT_LAUNCH"
	CampaignTypeSeasonal      CampaignType = "SEASONAL"
	CampaignTypePromotion     CampaignType = "PROMOTION"
	CampaignTypeAwareness     CampaignType = "AWARENESS"
	CampaignTypeInfluencer    CampaignType = "INFLUENCER"
)

// CampaignStatus represents campaign status
type CampaignStatus string

const (
	CampaignStatusDraft     CampaignStatus = "DRAFT"
	CampaignStatusPlanned   CampaignStatus = "PLANNED"
	CampaignStatusActive    CampaignStatus = "ACTIVE"
	CampaignStatusPaused    CampaignStatus = "PAUSED"
	CampaignStatusCompleted CampaignStatus = "COMPLETED"
	CampaignStatusCancelled CampaignStatus = "CANCELLED"
)

// Campaign represents a marketing campaign
type Campaign struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CampaignCode string         `json:"campaign_code" gorm:"uniqueIndex;size:50"`
	Name         string         `json:"name" gorm:"size:255;not null"`
	Description  string         `json:"description"`
	CampaignType CampaignType   `json:"campaign_type" gorm:"size:50;not null"`

	// Timeline
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	// Target
	TargetAudience string `json:"target_audience"`
	Channels       []byte `json:"channels" gorm:"type:jsonb"`

	// Budget
	Budget   float64 `json:"budget" gorm:"type:decimal(18,2)"`
	Spent    float64 `json:"spent" gorm:"type:decimal(18,2)"`
	Currency string  `json:"currency" gorm:"size:10;default:VND"`

	// Goals
	Goals []byte `json:"goals" gorm:"type:jsonb"`

	// Products
	Products []byte `json:"products" gorm:"type:jsonb"`

	// Performance
	Impressions      int     `json:"impressions"`
	Reach            int     `json:"reach"`
	Engagement       int     `json:"engagement"`
	Conversions      int     `json:"conversions"`
	RevenueGenerated float64 `json:"revenue_generated" gorm:"type:decimal(18,2)"`

	Status CampaignStatus `json:"status" gorm:"size:50;default:DRAFT"`

	Notes       string     `json:"notes"`
	CreatedBy   *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	UpdatedBy   *uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	LaunchedAt  *time.Time `json:"launched_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (Campaign) TableName() string {
	return "campaigns"
}

// CanBeLaunched checks if campaign can be launched
func (c *Campaign) CanBeLaunched() bool {
	return c.Status == CampaignStatusDraft || c.Status == CampaignStatusPlanned
}

// Launch sets campaign to active
func (c *Campaign) Launch() {
	now := time.Now()
	c.Status = CampaignStatusActive
	c.LaunchedAt = &now
	c.UpdatedAt = now
}

// Pause pauses the campaign
func (c *Campaign) Pause() {
	c.Status = CampaignStatusPaused
	c.UpdatedAt = time.Now()
}

// Complete marks campaign as completed
func (c *Campaign) Complete() {
	now := time.Now()
	c.Status = CampaignStatusCompleted
	c.CompletedAt = &now
	c.UpdatedAt = now
}

// Cancel cancels the campaign
func (c *Campaign) Cancel() {
	c.Status = CampaignStatusCancelled
	c.UpdatedAt = time.Now()
}

// GetROI calculates Return on Investment
func (c *Campaign) GetROI() float64 {
	if c.Spent == 0 {
		return 0
	}
	return ((c.RevenueGenerated - c.Spent) / c.Spent) * 100
}

// GetBudgetUtilization returns budget utilization percentage
func (c *Campaign) GetBudgetUtilization() float64 {
	if c.Budget == 0 {
		return 0
	}
	return (c.Spent / c.Budget) * 100
}

// IsActive checks if campaign is currently active
func (c *Campaign) IsActive() bool {
	now := time.Now()
	return c.Status == CampaignStatusActive && 
		now.After(c.StartDate) && 
		now.Before(c.EndDate)
}
