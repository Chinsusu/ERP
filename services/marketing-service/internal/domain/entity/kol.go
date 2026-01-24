package entity

import (
	"time"

	"github.com/google/uuid"
)

// KOLTier represents a tier classification for KOLs
type KOLTier struct {
	ID                 uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code               string    `json:"code" gorm:"uniqueIndex;size:50"`
	Name               string    `json:"name" gorm:"size:100"`
	Description        string    `json:"description"`
	MinFollowers       int       `json:"min_followers"`
	MaxFollowers       *int      `json:"max_followers"`
	AutoApproveSamples bool      `json:"auto_approve_samples"`
	DiscountPercent    float64   `json:"discount_percent" gorm:"type:decimal(5,2)"`
	Priority           int       `json:"priority"`
	IsActive           bool      `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (KOLTier) TableName() string {
	return "kol_tiers"
}

// KOLCategory represents the type of KOL
type KOLCategory string

const (
	KOLCategoryBeautyBlogger KOLCategory = "BEAUTY_BLOGGER"
	KOLCategoryInfluencer    KOLCategory = "INFLUENCER"
	KOLCategoryCelebrity     KOLCategory = "CELEBRITY"
	KOLCategoryExpert        KOLCategory = "EXPERT"
)

// KOLStatus represents KOL status
type KOLStatus string

const (
	KOLStatusActive      KOLStatus = "ACTIVE"
	KOLStatusInactive    KOLStatus = "INACTIVE"
	KOLStatusBlacklisted KOLStatus = "BLACKLISTED"
)

// Platform represents social media platform
type Platform string

const (
	PlatformInstagram Platform = "INSTAGRAM"
	PlatformYouTube   Platform = "YOUTUBE"
	PlatformTikTok    Platform = "TIKTOK"
	PlatformFacebook  Platform = "FACEBOOK"
)

// KOL represents a Key Opinion Leader / Influencer
type KOL struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	KOLCode  string    `json:"kol_code" gorm:"uniqueIndex;size:50"`
	Name     string    `json:"name" gorm:"size:200;not null"`
	Email    string    `json:"email" gorm:"size:255"`
	Phone    string    `json:"phone" gorm:"size:50"`

	TierID   *uuid.UUID `json:"tier_id" gorm:"type:uuid"`
	Tier     *KOLTier   `json:"tier,omitempty" gorm:"foreignKey:TierID"`
	Category KOLCategory `json:"category" gorm:"size:100"`

	// Social Media
	InstagramHandle    string `json:"instagram_handle" gorm:"size:100"`
	InstagramFollowers int    `json:"instagram_followers"`
	YouTubeChannel     string `json:"youtube_channel" gorm:"size:200"`
	YouTubeSubscribers int    `json:"youtube_subscribers"`
	TikTokHandle       string `json:"tiktok_handle" gorm:"size:100"`
	TikTokFollowers    int    `json:"tiktok_followers"`
	FacebookPage       string `json:"facebook_page" gorm:"size:200"`
	FacebookFollowers  int    `json:"facebook_followers"`

	// Engagement
	AvgEngagementRate float64 `json:"avg_engagement_rate" gorm:"type:decimal(5,2)"`
	Niche             string  `json:"niche" gorm:"size:100"`

	// Business
	CollaborationRate float64 `json:"collaboration_rate" gorm:"type:decimal(18,2)"`
	Currency          string  `json:"currency" gorm:"size:10;default:VND"`
	PreferredProducts []byte  `json:"preferred_products" gorm:"type:jsonb"`

	// Address
	AddressLine1 string `json:"address_line1" gorm:"size:255"`
	AddressLine2 string `json:"address_line2" gorm:"size:255"`
	City         string `json:"city" gorm:"size:100"`
	State        string `json:"state" gorm:"size:100"`
	PostalCode   string `json:"postal_code" gorm:"size:20"`
	Country      string `json:"country" gorm:"size:100;default:Vietnam"`

	// Performance
	TotalPosts             int        `json:"total_posts"`
	TotalSamplesReceived   int        `json:"total_samples_received"`
	TotalCollaborations    int        `json:"total_collaborations"`
	LastCollaborationDate  *time.Time `json:"last_collaboration_date"`

	Notes  string    `json:"notes"`
	Status KOLStatus `json:"status" gorm:"size:50;default:ACTIVE"`

	CreatedBy *uuid.UUID `json:"created_by" gorm:"type:uuid"`
	UpdatedBy *uuid.UUID `json:"updated_by" gorm:"type:uuid"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (KOL) TableName() string {
	return "kols"
}

// GetTotalFollowers returns the sum of followers across all platforms
func (k *KOL) GetTotalFollowers() int {
	return k.InstagramFollowers + k.YouTubeSubscribers + k.TikTokFollowers + k.FacebookFollowers
}

// GetPrimaryPlatform returns the platform with the most followers
func (k *KOL) GetPrimaryPlatform() Platform {
	max := k.InstagramFollowers
	platform := PlatformInstagram

	if k.YouTubeSubscribers > max {
		max = k.YouTubeSubscribers
		platform = PlatformYouTube
	}
	if k.TikTokFollowers > max {
		max = k.TikTokFollowers
		platform = PlatformTikTok
	}
	if k.FacebookFollowers > max {
		platform = PlatformFacebook
	}

	return platform
}

// Block sets KOL status to blacklisted
func (k *KOL) Block() {
	k.Status = KOLStatusBlacklisted
	k.UpdatedAt = time.Now()
}

// Activate sets KOL status to active
func (k *KOL) Activate() {
	k.Status = KOLStatusActive
	k.UpdatedAt = time.Now()
}
