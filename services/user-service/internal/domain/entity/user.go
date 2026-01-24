package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	EmployeeCode string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"employee_code"`
	FirstName    string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName     string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Phone        string         `gorm:"type:varchar(20)" json:"phone,omitempty"`
	AvatarURL    string         `gorm:"type:varchar(500)" json:"avatar_url,omitempty"`
	DepartmentID *uuid.UUID     `gorm:"type:uuid" json:"department_id,omitempty"`
	ManagerID    *uuid.UUID     `gorm:"type:uuid" json:"manager_id,omitempty"`
	Status       string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relationships
	Department *Department  `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Manager    *User        `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
	Profile    *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// IsActive checks if user is active
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// Validate validates user data
func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.FirstName == "" {
		return fmt.Errorf("first name is required")
	}
	if u.LastName == "" {
		return fmt.Errorf("last name is required")
	}
	if u.Status != "active" && u.Status != "inactive" && u.Status != "suspended" {
		return fmt.Errorf("invalid status: %s", u.Status)
	}
	return nil
}

// GenerateEmployeeCode generates employee code if not provided
func (u *User) GenerateEmployeeCode(sequence int) {
	if u.EmployeeCode == "" {
		u.EmployeeCode = fmt.Sprintf("EMP%s%03d", time.Now().Format("20060102"), sequence)
	}
}
