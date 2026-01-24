package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents user credentials and authentication information
type User struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID               uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Email                string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email"`
	PasswordHash         string         `gorm:"type:varchar(255);not null" json:"-"`
	IsActive             bool           `gorm:"default:true" json:"is_active"`
	EmailVerified        bool           `gorm:"default:false" json:"email_verified"`
	FailedLoginAttempts  int            `gorm:"default:0" json:"failed_login_attempts"`
	LockedUntil          *time.Time     `gorm:"type:timestamp" json:"locked_until,omitempty"`
	LastLoginAt          *time.Time     `gorm:"type:timestamp" json:"last_login_at,omitempty"`
	CreatedAt            time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Relationships
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "user_credentials"
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// VerifyPassword checks if the provided password matches the stored hash
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsLocked checks if the account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IncrementFailedAttempts increments the failed login counter
func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++
	// Lock account for 30 minutes after 5 failed attempts
	if u.FailedLoginAttempts >= 5 {
		lockUntil := time.Now().Add(30 * time.Minute)
		u.LockedUntil = &lockUntil
	}
}

// ResetFailedAttempts resets the failed login counter
func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}
