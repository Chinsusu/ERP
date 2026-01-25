package mocks

import (
	"context"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCacheRepository is a mock implementation of repository.CacheRepository
type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) SetUserPermissions(ctx context.Context, userID uuid.UUID, permissions []entity.Permission, ttl time.Duration) error {
	args := m.Called(ctx, userID, permissions, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Permission), args.Error(1)
}

func (m *MockCacheRepository) DeleteUserPermissions(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCacheRepository) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
	args := m.Called(ctx, jti, expiresAt)
	return args.Error(0)
}

func (m *MockCacheRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	args := m.Called(ctx, jti)
	return args.Bool(0), args.Error(1)
}

// Ensure MockCacheRepository implements repository.CacheRepository
var _ repository.CacheRepository = (*MockCacheRepository)(nil)
