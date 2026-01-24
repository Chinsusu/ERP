package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Permission represents a fine-grained permission in the system
type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Code        string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"code"`
	Name        string    `gorm:"type:varchar(200);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Service     string    `gorm:"type:varchar(50);not null;index:idx_permissions_service" json:"service"`
	Resource    string    `gorm:"type:varchar(50);not null" json:"resource"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (Permission) TableName() string {
	return "permissions"
}

// GenerateCode generates permission code from service, resource, and action
func GeneratePermissionCode(service, resource, action string) string {
	return fmt.Sprintf("%s:%s:%s", service, resource, action)
}

// ParseCode parses a permission code into service, resource, and action
func (p *Permission) ParseCode() (service, resource, action string) {
	parts := strings.Split(p.Code, ":")
	if len(parts) == 3 {
		return parts[0], parts[1], parts[2]
	}
	return "", "", ""
}

// Matches checks if this permission matches the given permission code
// Supports wildcard matching (e.g., wms:*:read matches wms:stock:read)
func (p *Permission) Matches(permissionCode string) bool {
	// Exact match
	if p.Code == permissionCode {
		return true
	}
	
	// Wildcard match
	thisParts := strings.Split(p.Code, ":")
	checkParts := strings.Split(permissionCode, ":")
	
	if len(thisParts) != 3 || len(checkParts) != 3 {
		return false
	}
	
	// Full wildcard
	if thisParts[0] == "*" && thisParts[1] == "*" && thisParts[2] == "*" {
		return true
	}
	
	// Check each part
	for i := 0; i < 3; i++ {
		if thisParts[i] != "*" && thisParts[i] != checkParts[i] {
			return false
		}
	}
	
	return true
}

// IsWildcard checks if the permission contains any wildcards
func (p *Permission) IsWildcard() bool {
	return strings.Contains(p.Code, "*")
}
