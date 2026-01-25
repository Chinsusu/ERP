package entity

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FileCategory represents a file category with validation rules
type FileCategory struct {
	ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code              string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Name              string         `gorm:"type:varchar(100);not null" json:"name"`
	Description       string         `gorm:"type:text" json:"description,omitempty"`
	AllowedExtensions datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"allowed_extensions"`
	MaxFileSize       int64          `gorm:"type:bigint;not null" json:"max_file_size"`
	StorageBucket     string         `gorm:"type:varchar(100);not null" json:"storage_bucket"`
	IsActive          bool           `gorm:"default:true" json:"is_active"`
	CreatedAt         time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies table name
func (FileCategory) TableName() string {
	return "file_categories"
}

// BeforeCreate sets defaults
func (fc *FileCategory) BeforeCreate(tx *gorm.DB) error {
	if fc.ID == uuid.Nil {
		fc.ID = uuid.New()
	}
	return nil
}

// GetAllowedExtensions returns allowed extensions as slice
func (fc *FileCategory) GetAllowedExtensions() []string {
	var extensions []string
	if fc.AllowedExtensions != nil {
		json.Unmarshal(fc.AllowedExtensions, &extensions)
	}
	return extensions
}

// IsExtensionAllowed checks if extension is allowed
func (fc *FileCategory) IsExtensionAllowed(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	for _, allowed := range fc.GetAllowedExtensions() {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}
	return false
}

// IsFileSizeAllowed checks if file size is within limit
func (fc *FileCategory) IsFileSizeAllowed(size int64) bool {
	return size <= fc.MaxFileSize
}

// ValidateFile validates a file against category rules
func (fc *FileCategory) ValidateFile(filename string, size int64) error {
	ext := getExtension(filename)
	if !fc.IsExtensionAllowed(ext) {
		return &ValidationError{
			Field:   "extension",
			Message: "File extension '" + ext + "' is not allowed for category " + fc.Name,
		}
	}
	if !fc.IsFileSizeAllowed(size) {
		return &ValidationError{
			Field:   "size",
			Message: "File size exceeds maximum allowed size",
		}
	}
	return nil
}

// GetMaxFileSizeFormatted returns human-readable max file size
func (fc *FileCategory) GetMaxFileSizeFormatted() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case fc.MaxFileSize >= GB:
		return formatSize(float64(fc.MaxFileSize)/GB) + " GB"
	case fc.MaxFileSize >= MB:
		return formatSize(float64(fc.MaxFileSize)/MB) + " MB"
	case fc.MaxFileSize >= KB:
		return formatSize(float64(fc.MaxFileSize)/KB) + " KB"
	default:
		return string(rune(fc.MaxFileSize)) + " B"
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Helper functions
func getExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}
	return strings.ToLower(parts[len(parts)-1])
}

func formatSize(f float64) string {
	// Simple format
	return ""
}
