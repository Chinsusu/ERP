package entity

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a long-lived refresh token for JWT rotation
type RefreshToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"type:timestamp;not null;index" json:"expires_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	RevokedAt *time.Time `gorm:"type:timestamp" json:"revoked_at,omitempty"`
	RevokedBy *uuid.UUID `gorm:"type:uuid" json:"revoked_by,omitempty"`
	IPAddress string     `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent string     `gorm:"type:text" json:"user_agent,omitempty"`
}

// TableName specifies the table name
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the token has expired
func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsRevoked checks if the token has been revoked
func (t *RefreshToken) IsRevoked() bool {
	return t.RevokedAt != nil
}

// IsValid checks if the token is valid (not expired and not revoked)
func (t *RefreshToken) IsValid() bool {
	return !t.IsExpired() && !t.IsRevoked()
}

// Revoke revokes the token
func (t *RefreshToken) Revoke(revokedBy uuid.UUID) {
	now := time.Now()
	t.RevokedAt = &now
	t.RevokedBy = &revokedBy
}

// PasswordResetToken represents a one-time password reset token
type PasswordResetToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string     `gorm:"type:varchar(255);not null;uniqueIndex" json:"-"`
	ExpiresAt time.Time  `gorm:"type:timestamp;not null;index" json:"expires_at"`
	UsedAt    *time.Time `gorm:"type:timestamp" json:"used_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

// TableName specifies the table name
func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

// IsExpired checks if the token has expired
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsUsed checks if the token has been used
func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}

// IsValid checks if the token is valid (not expired and not used)
func (t *PasswordResetToken) IsValid() bool {
	return !t.IsExpired() && !t.IsUsed()
}

// MarkAsUsed marks the token as used
func (t *PasswordResetToken) MarkAsUsed() {
	now := time.Now()
	t.UsedAt = &now
}

// Session represents an active user session
type Session struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	AccessTokenJTI  string     `gorm:"type:varchar(100);not null;uniqueIndex" json:"access_token_jti"`
	RefreshTokenID  *uuid.UUID `gorm:"type:uuid" json:"refresh_token_id,omitempty"`
	IPAddress       string     `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent       string     `gorm:"type:text" json:"user_agent,omitempty"`
	ExpiresAt       time.Time  `gorm:"type:timestamp;not null;index" json:"expires_at"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	LastActivity    time.Time  `gorm:"type:timestamp;not null" json:"last_activity"`
}

// TableName specifies the table name
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session has expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// UpdateActivity updates the last activity timestamp
func (s *Session) UpdateActivity() {
	s.LastActivity = time.Now()
}
