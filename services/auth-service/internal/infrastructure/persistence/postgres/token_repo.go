package postgres

import (
	"context"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) repository.TokenRepository {
	return &tokenRepository{db: db}
}

// Refresh Tokens
func (r *tokenRepository) CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string, revokedBy uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"revoked_by": revokedBy,
		}).Error
}

func (r *tokenRepository) RevokeAllUserTokens(ctx context.Context, userID, revokedBy uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"revoked_by": revokedBy,
		}).Error
}

func (r *tokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? OR revoked_at IS NOT NULL", time.Now()).
		Delete(&entity.RefreshToken{}).Error
}

// Password Reset Tokens
func (r *tokenRepository) CreatePasswordResetToken(ctx context.Context, token *entity.PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	var token entity.PasswordResetToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *tokenRepository) MarkResetTokenAsUsed(ctx context.Context, tokenHash string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.PasswordResetToken{}).
		Where("token_hash = ?", tokenHash).
		Update("used_at", now).Error
}

// Sessions
func (r *tokenRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *tokenRepository) GetSession(ctx context.Context, jti string) (*entity.Session, error) {
	var session entity.Session
	err := r.db.WithContext(ctx).Where("access_token_jti = ?", jti).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *tokenRepository) DeleteSession(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).Where("access_token_jti = ?", jti).Delete(&entity.Session{}).Error
}

func (r *tokenRepository) UpdateSessionActivity(ctx context.Context, jti string) error {
	return r.db.WithContext(ctx).Model(&entity.Session{}).
		Where("access_token_jti = ?", jti).
		Update("last_activity", time.Now()).Error
}

func (r *tokenRepository) DeleteExpiredSessions(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&entity.Session{}).Error
}
