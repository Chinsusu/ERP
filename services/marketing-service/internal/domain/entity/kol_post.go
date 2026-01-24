package entity

import (
	"time"

	"github.com/google/uuid"
)

// PostType represents the type of social media post
type PostType string

const (
	PostTypePost  PostType = "POST"
	PostTypeStory PostType = "STORY"
	PostTypeReel  PostType = "REEL"
	PostTypeVideo PostType = "VIDEO"
	PostTypeLive  PostType = "LIVE"
)

// Sentiment represents post sentiment
type Sentiment string

const (
	SentimentPositive Sentiment = "POSITIVE"
	SentimentNeutral  Sentiment = "NEUTRAL"
	SentimentNegative Sentiment = "NEGATIVE"
)

// KOLPost represents a social media post by KOL
type KOLPost struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	KOLID           uuid.UUID  `json:"kol_id" gorm:"type:uuid;not null"`
	KOL             *KOL       `json:"kol,omitempty" gorm:"foreignKey:KOLID"`
	CollaborationID *uuid.UUID `json:"collaboration_id" gorm:"type:uuid"`
	CampaignID      *uuid.UUID `json:"campaign_id" gorm:"type:uuid"`

	Platform Platform `json:"platform" gorm:"size:50;not null"`
	PostType PostType `json:"post_type" gorm:"size:50"`

	PostURL  string    `json:"post_url"`
	PostDate time.Time `json:"post_date"`

	Caption           string `json:"caption"`
	ProductsMentioned []byte `json:"products_mentioned" gorm:"type:jsonb"`

	// Metrics
	Likes          int     `json:"likes"`
	Comments       int     `json:"comments"`
	Shares         int     `json:"shares"`
	Views          int     `json:"views"`
	Reach          int     `json:"reach"`
	EngagementRate float64 `json:"engagement_rate" gorm:"type:decimal(5,2)"`

	// Analysis
	Sentiment Sentiment `json:"sentiment" gorm:"size:50"`
	Summary   string    `json:"summary"`

	ScreenshotURLs []byte `json:"screenshot_urls" gorm:"type:jsonb"`

	Verified   bool       `json:"verified"`
	VerifiedBy *uuid.UUID `json:"verified_by" gorm:"type:uuid"`
	VerifiedAt *time.Time `json:"verified_at"`

	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (KOLPost) TableName() string {
	return "kol_posts"
}

// CalculateEngagementRate calculates engagement rate
func (p *KOLPost) CalculateEngagementRate(followers int) {
	if followers == 0 {
		p.EngagementRate = 0
		return
	}
	totalEngagement := p.Likes + p.Comments + p.Shares
	p.EngagementRate = (float64(totalEngagement) / float64(followers)) * 100
}

// Verify marks the post as verified
func (p *KOLPost) Verify(verifierID uuid.UUID) {
	now := time.Now()
	p.Verified = true
	p.VerifiedBy = &verifierID
	p.VerifiedAt = &now
	p.UpdatedAt = now
}

// CollaborationType represents type of collaboration
type CollaborationType string

const (
	CollabTypeSponsoredPost  CollaborationType = "SPONSORED_POST"
	CollabTypeProductReview  CollaborationType = "PRODUCT_REVIEW"
	CollabTypeGiveaway       CollaborationType = "GIVEAWAY"
	CollabTypeAmbassador     CollaborationType = "AMBASSADOR"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPartial PaymentStatus = "PARTIAL"
	PaymentStatusPaid    PaymentStatus = "PAID"
)

// CollaborationStatus represents collaboration status
type CollaborationStatus string

const (
	CollabStatusDraft     CollaborationStatus = "DRAFT"
	CollabStatusActive    CollaborationStatus = "ACTIVE"
	CollabStatusCompleted CollaborationStatus = "COMPLETED"
	CollabStatusCancelled CollaborationStatus = "CANCELLED"
)

// KOLCollaboration represents a collaboration between brand and KOL
type KOLCollaboration struct {
	ID                uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CollaborationCode string     `json:"collaboration_code" gorm:"uniqueIndex;size:50"`
	CampaignID        *uuid.UUID `json:"campaign_id" gorm:"type:uuid"`
	Campaign          *Campaign  `json:"campaign,omitempty" gorm:"foreignKey:CampaignID"`
	KOLID             uuid.UUID  `json:"kol_id" gorm:"type:uuid;not null"`
	KOL               *KOL       `json:"kol,omitempty" gorm:"foreignKey:KOLID"`

	CollaborationType CollaborationType `json:"collaboration_type" gorm:"size:50"`
	AgreedFee         float64           `json:"agreed_fee" gorm:"type:decimal(18,2)"`
	Currency          string            `json:"currency" gorm:"size:10;default:VND"`

	ExpectedPosts int    `json:"expected_posts"`
	ActualPosts   int    `json:"actual_posts"`
	Platforms     []byte `json:"platforms" gorm:"type:jsonb"`

	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	ContentDeadline *time.Time `json:"content_deadline"`

	// Performance
	TotalImpressions int `json:"total_impressions"`
	TotalEngagement  int `json:"total_engagement"`
	TotalReach       int `json:"total_reach"`

	// Payment
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"size:50;default:PENDING"`
	PaidAmount    float64       `json:"paid_amount" gorm:"type:decimal(18,2)"`

	Status CollaborationStatus `json:"status" gorm:"size:50;default:DRAFT"`
	Notes  string              `json:"notes"`

	CreatedBy *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (KOLCollaboration) TableName() string {
	return "kol_collaborations"
}

// Activate activates the collaboration
func (c *KOLCollaboration) Activate() {
	c.Status = CollabStatusActive
	c.UpdatedAt = time.Now()
}

// Complete marks collaboration as completed
func (c *KOLCollaboration) Complete() {
	c.Status = CollabStatusCompleted
	c.UpdatedAt = time.Now()
}

// Cancel cancels the collaboration
func (c *KOLCollaboration) Cancel() {
	c.Status = CollabStatusCancelled
	c.UpdatedAt = time.Now()
}

// MarkPaid marks collaboration as fully paid
func (c *KOLCollaboration) MarkPaid() {
	c.PaymentStatus = PaymentStatusPaid
	c.PaidAmount = c.AgreedFee
	c.UpdatedAt = time.Now()
}
