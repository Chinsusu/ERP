package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// File category constants
const (
	CategoryDocument     = "DOCUMENT"
	CategoryImage        = "IMAGE"
	CategoryCertificate  = "CERTIFICATE"
	CategoryContract     = "CONTRACT"
	CategoryReport       = "REPORT"
	CategoryAvatar       = "AVATAR"
	CategoryProductImage = "PRODUCT_IMAGE"
	CategoryQCPhoto      = "QC_PHOTO"
	CategorySignature    = "SIGNATURE"
	CategoryAttachment   = "ATTACHMENT"
)

// File represents a stored file
type File struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OriginalName string         `gorm:"type:varchar(500);not null" json:"original_name"`
	StoredName   string         `gorm:"type:varchar(500);not null" json:"stored_name"`
	ContentType  string         `gorm:"type:varchar(100);not null" json:"content_type"`
	FileSize     int64          `gorm:"type:bigint;not null" json:"file_size"`
	BucketName   string         `gorm:"type:varchar(100);not null" json:"bucket_name"`
	ObjectPath   string         `gorm:"type:text;not null" json:"object_path"`
	Category     string         `gorm:"type:varchar(50);not null" json:"category"`
	EntityType   string         `gorm:"type:varchar(50)" json:"entity_type,omitempty"`
	EntityID     *uuid.UUID     `gorm:"type:uuid" json:"entity_id,omitempty"`
	IsPublic     bool           `gorm:"default:false" json:"is_public"`
	AccessToken  string         `gorm:"type:varchar(255)" json:"access_token,omitempty"`
	Metadata     datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`
	Checksum     string         `gorm:"type:varchar(64)" json:"checksum,omitempty"`
	CreatedBy    *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	ExpiresAt    *time.Time     `gorm:"type:timestamp" json:"expires_at,omitempty"`
}

// TableName specifies table name
func (File) TableName() string {
	return "files"
}

// BeforeCreate generates stored name and access token
func (f *File) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	if f.StoredName == "" {
		ext := filepath.Ext(f.OriginalName)
		f.StoredName = f.ID.String() + ext
	}
	if f.AccessToken == "" {
		f.AccessToken = generateAccessToken()
	}
	if f.ObjectPath == "" {
		f.ObjectPath = f.StoredName
	}
	return nil
}

// GetExtension returns file extension
func (f *File) GetExtension() string {
	ext := filepath.Ext(f.OriginalName)
	return strings.TrimPrefix(strings.ToLower(ext), ".")
}

// IsImage checks if file is an image
func (f *File) IsImage() bool {
	imageTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml"}
	for _, t := range imageTypes {
		if f.ContentType == t {
			return true
		}
	}
	return false
}

// IsPDF checks if file is a PDF
func (f *File) IsPDF() bool {
	return f.ContentType == "application/pdf"
}

// GetFullPath returns full object path in storage
func (f *File) GetFullPath() string {
	return f.BucketName + "/" + f.ObjectPath
}

// IsExpired checks if file has expired
func (f *File) IsExpired() bool {
	if f.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*f.ExpiresAt)
}

// GetFileSizeFormatted returns human-readable file size
func (f *File) GetFileSizeFormatted() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case f.FileSize >= GB:
		return formatFloat(float64(f.FileSize)/GB) + " GB"
	case f.FileSize >= MB:
		return formatFloat(float64(f.FileSize)/MB) + " MB"
	case f.FileSize >= KB:
		return formatFloat(float64(f.FileSize)/KB) + " KB"
	default:
		return formatInt(f.FileSize) + " B"
	}
}

// Helper functions
func generateAccessToken() string {
	token := uuid.New().String() + uuid.New().String()
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:16])
}

func formatFloat(f float64) string {
	return strings.TrimRight(strings.TrimRight(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					formatFloatWithPrec(f, 2),
					".00", "", 1),
				".0", "", 1),
			"0", "", -1),
		"."),
		"0")
}

func formatFloatWithPrec(f float64, prec int) string {
	return strings.TrimRight(strings.TrimRight(
		formatFloatBasic(f, prec), "0"), ".")
}

func formatFloatBasic(f float64, prec int) string {
	// Simple formatting
	return ""
}

func formatInt(i int64) string {
	return ""
}
