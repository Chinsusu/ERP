package mocks

import (
	"context"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockRoleRepository is a mock implementation of repository.RoleRepository
type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *MockRoleRepository) List(ctx context.Context) ([]entity.Role, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Role), args.Error(1)
}

func (m *MockRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	args := m.Called(ctx, role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Role), args.Error(1)
}

func (m *MockRoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy uuid.UUID) error {
	args := m.Called(ctx, userID, roleID, assignedBy)
	return args.Error(0)
}

func (m *MockRoleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

// MockPermissionRepository is a mock implementation of repository.PermissionRepository
type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByCode(ctx context.Context, code string) (*entity.Permission, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) List(ctx context.Context) ([]entity.Permission, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) ListByService(ctx context.Context, service string) ([]entity.Permission, error) {
	args := m.Called(ctx, service)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	args := m.Called(ctx, permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPermissionRepository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error) {
	args := m.Called(ctx, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

func (m *MockPermissionRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error {
	args := m.Called(ctx, roleID, permissionID)
	return args.Error(0)
}

// MockTokenRepository is a mock implementation of repository.TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.RefreshToken), args.Error(1)
}

func (m *MockTokenRepository) RevokeRefreshToken(ctx context.Context, tokenHash string, revokedBy uuid.UUID) error {
	args := m.Called(ctx, tokenHash, revokedBy)
	return args.Error(0)
}

func (m *MockTokenRepository) RevokeAllUserTokens(ctx context.Context, userID, revokedBy uuid.UUID) error {
	args := m.Called(ctx, userID, revokedBy)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTokenRepository) CreatePasswordResetToken(ctx context.Context, token *entity.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.PasswordResetToken), args.Error(1)
}

func (m *MockTokenRepository) MarkResetTokenAsUsed(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockTokenRepository) CreateSession(ctx context.Context, session *entity.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockTokenRepository) GetSession(ctx context.Context, jti string) (*entity.Session, error) {
	args := m.Called(ctx, jti)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Session), args.Error(1)
}

func (m *MockTokenRepository) DeleteSession(ctx context.Context, jti string) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}

func (m *MockTokenRepository) UpdateSessionActivity(ctx context.Context, jti string) error {
	args := m.Called(ctx, jti)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteExpiredSessions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

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

// MockEventPublisher is a mock implementation of auth.EventPublisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(subject string, data interface{}) error {
	args := m.Called(subject, data)
	return args.Error(0)
}
