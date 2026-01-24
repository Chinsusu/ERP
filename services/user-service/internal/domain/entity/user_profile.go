package entity

import (
	"time"

	"github.com/google/uuid"
)

// UserProfile represents extended user information
type UserProfile struct {
	UserID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"user_id"`
	DateOfBirth      *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Address          string     `gorm:"type:text" json:"address,omitempty"`
	EmergencyContact string     `gorm:"type:varchar(255)" json:"emergency_contact,omitempty"`
	JoinDate         *time.Time `gorm:"type:date" json:"join_date,omitempty"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationship
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name
func (UserProfile) TableName() string {
	return "user_profiles"
}
