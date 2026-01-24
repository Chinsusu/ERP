package repository

import (
	"context"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

// TokenRepository defines the interface for refresh token data access
type TokenRepository interface {
	// CreateRefreshToken creates a new refresh token
	CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	
	// GetRefreshToken retrieves a refresh token by hash
	GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	
	// RevokeRefreshToken revokes a refresh token
	RevokeRefreshToken(ctx context.Context, tokenHash string, revokedBy uuid.UUID) error
	
	// RevokeAllUserTokens revokes all tokens for a user
	RevokeAllUserTokens(ctx context.Context, userID, revokedBy uuid.UUID) error
	
	// DeleteExpiredTokens deletes expired tokens (cleanup)
	DeleteExpiredTokens(ctx context.Context) error
	
	// CreatePasswordResetToken creates a new password reset token
	CreatePasswordResetToken(ctx context.Context, token *entity.PasswordResetToken) error
	
	// GetPasswordResetToken retrieves a password reset token by hash
	GetPasswordResetToken(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	
	// MarkResetTokenAsUsed marks a reset token as used
	MarkResetTokenAsUsed(ctx context.Context, tokenHash string) error
	
	// CreateSession creates a new session
	CreateSession(ctx context.Context, session *entity.Session) error
	
	// GetSession retrieves a session by JWT ID
	GetSession(ctx context.Context, jti string) (*entity.Session, error)
	
	// DeleteSession deletes a session
	DeleteSession(ctx context.Context, jti string) error
	
	// UpdateSessionActivity updates the last activity timestamp
	UpdateSessionActivity(ctx context.Context, jti string) error
	
	// DeleteExpiredSessions deletes expired sessions (cleanup)
	DeleteExpiredSessions(ctx context.Context) error
}

// CacheRepository defines the interface for caching operations
type CacheRepository interface {
	// SetUserPermissions caches user permissions
	SetUserPermissions(ctx context.Context, userID uuid.UUID, permissions []entity.Permission, ttl time.Duration) error
	
	// GetUserPermissions retrieves cached user permissions
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error)
	
	// DeleteUserPermissions deletes cached user permissions
	DeleteUserPermissions(ctx context.Context, userID uuid.UUID) error
	
	// BlacklistToken adds a token to the blacklist
	BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error
	
	// IsTokenBlacklisted checks if a token is blacklisted
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
}
