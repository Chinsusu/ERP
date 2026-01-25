package permission_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth/mocks"
	"github.com/erp-cosmetics/auth-service/internal/usecase/permission"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createPermission(code string) entity.Permission {
	parts := splitCode(code)
	return entity.Permission{
		ID:       uuid.New(),
		Code:     code,
		Name:     "Test Permission",
		Service:  parts[0],
		Resource: parts[1],
		Action:   parts[2],
	}
}

func splitCode(code string) [3]string {
	result := [3]string{"", "", ""}
	start := 0
	idx := 0
	for i, ch := range code {
		if ch == ':' && idx < 2 {
			result[idx] = code[start:i]
			start = i + 1
			idx++
		}
	}
	result[idx] = code[start:]
	return result
}

// TestGetPermissionsUseCase_Execute_FromCache tests cache hit scenario
func TestGetPermissionsUseCase_Execute_FromCache(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	cachedPermissions := []entity.Permission{
		createPermission("user:user:read"),
		createPermission("user:user:create"),
	}

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(cachedPermissions, nil)

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act
	result, err := uc.Execute(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, cachedPermissions, result)
	assert.Len(t, result, 2)

	// permRepo should NOT be called
	permRepo.AssertNotCalled(t, "GetUserPermissions", mock.Anything, mock.Anything)
	cacheRepo.AssertExpectations(t)
}

// TestGetPermissionsUseCase_Execute_FromDatabase tests cache miss with DB fetch
func TestGetPermissionsUseCase_Execute_FromDatabase(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	dbPermissions := []entity.Permission{
		createPermission("wms:stock:read"),
		createPermission("wms:stock:update"),
	}

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(nil, errors.New("cache miss"))
	permRepo.On("GetUserPermissions", ctx, userID).Return(dbPermissions, nil)
	cacheRepo.On("SetUserPermissions", ctx, userID, dbPermissions, 15*time.Minute).Return(nil)

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act
	result, err := uc.Execute(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, dbPermissions, result)

	permRepo.AssertExpectations(t)
	cacheRepo.AssertExpectations(t)
}

// TestGetPermissionsUseCase_Execute_DatabaseError tests DB error handling
func TestGetPermissionsUseCase_Execute_DatabaseError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(nil, errors.New("cache miss"))
	permRepo.On("GetUserPermissions", ctx, userID).Return(nil, errors.New("database error"))

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act
	result, err := uc.Execute(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)

	permRepo.AssertExpectations(t)
	cacheRepo.AssertExpectations(t)
}

// TestGetPermissionsUseCase_CheckPermission_ExactMatch tests exact permission match
func TestGetPermissionsUseCase_CheckPermission_ExactMatch(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	permissions := []entity.Permission{
		createPermission("user:user:read"),
		createPermission("wms:stock:read"),
	}

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act
	allowed, err := uc.CheckPermission(ctx, userID, "wms:stock:read")

	// Assert
	assert.NoError(t, err)
	assert.True(t, allowed)

	cacheRepo.AssertExpectations(t)
}

// TestGetPermissionsUseCase_CheckPermission_Wildcard tests wildcard permission match
func TestGetPermissionsUseCase_CheckPermission_Wildcard(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	permissions := []entity.Permission{
		createPermission("wms:*:read"), // Wildcard: all wms resources with read action
	}

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act
	allowed, err := uc.CheckPermission(ctx, userID, "wms:stock:read")

	// Assert
	assert.NoError(t, err)
	assert.True(t, allowed)

	cacheRepo.AssertExpectations(t)
}

// TestGetPermissionsUseCase_CheckPermission_SuperWildcard tests super wildcard (*:*:*)
func TestGetPermissionsUseCase_CheckPermission_SuperWildcard(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	permissions := []entity.Permission{
		createPermission("*:*:*"), // Super admin - all permissions
	}

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act & Assert - should match any permission
	tests := []string{
		"user:user:read",
		"wms:stock:update",
		"procurement:po:create",
		"manufacturing:bom:delete",
	}

	for _, perm := range tests {
		allowed, err := uc.CheckPermission(ctx, userID, perm)
		assert.NoError(t, err)
		assert.True(t, allowed, "Super wildcard should match: %s", perm)
	}
}

// TestGetPermissionsUseCase_CheckPermission_Denied tests permission denied
func TestGetPermissionsUseCase_CheckPermission_Denied(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	permissions := []entity.Permission{
		createPermission("user:user:read"),
	}

	permRepo := &mocks.MockPermissionRepository{}
	cacheRepo := &mocks.MockCacheRepository{}

	cacheRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

	uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

	// Act
	allowed, err := uc.CheckPermission(ctx, userID, "wms:stock:delete")

	// Assert
	assert.NoError(t, err)
	assert.False(t, allowed)

	cacheRepo.AssertExpectations(t)
}

// Table-driven tests for permission matching scenarios
func TestGetPermissionsUseCase_CheckPermission_Scenarios(t *testing.T) {
	tests := []struct {
		name           string
		userPerms      []string
		checkPerm      string
		expectedResult bool
	}{
		{
			name:           "exact match",
			userPerms:      []string{"user:user:read"},
			checkPerm:      "user:user:read",
			expectedResult: true,
		},
		{
			name:           "no match - different service",
			userPerms:      []string{"user:user:read"},
			checkPerm:      "wms:user:read",
			expectedResult: false,
		},
		{
			name:           "no match - different action",
			userPerms:      []string{"user:user:read"},
			checkPerm:      "user:user:delete",
			expectedResult: false,
		},
		{
			name:           "wildcard resource match",
			userPerms:      []string{"wms:*:read"},
			checkPerm:      "wms:stock:read",
			expectedResult: true,
		},
		{
			name:           "wildcard resource - wrong action",
			userPerms:      []string{"wms:*:read"},
			checkPerm:      "wms:stock:delete",
			expectedResult: false,
		},
		{
			name:           "wildcard action match",
			userPerms:      []string{"wms:stock:*"},
			checkPerm:      "wms:stock:delete",
			expectedResult: true,
		},
		{
			name:           "wildcard service match",
			userPerms:      []string{"*:user:read"},
			checkPerm:      "auth:user:read",
			expectedResult: true,
		},
		{
			name:           "super wildcard matches everything",
			userPerms:      []string{"*:*:*"},
			checkPerm:      "any:thing:here",
			expectedResult: true,
		},
		{
			name:           "multiple permissions - one matches",
			userPerms:      []string{"user:user:read", "wms:stock:*", "procurement:po:create"},
			checkPerm:      "wms:stock:delete",
			expectedResult: true,
		},
		{
			name:           "multiple permissions - none match",
			userPerms:      []string{"user:user:read", "wms:stock:read"},
			checkPerm:      "manufacturing:bom:update",
			expectedResult: false,
		},
		{
			name:           "empty permissions",
			userPerms:      []string{},
			checkPerm:      "user:user:read",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			userID := uuid.New()

			permissions := make([]entity.Permission, len(tt.userPerms))
			for i, code := range tt.userPerms {
				permissions[i] = createPermission(code)
			}

			permRepo := &mocks.MockPermissionRepository{}
			cacheRepo := &mocks.MockCacheRepository{}

			cacheRepo.On("GetUserPermissions", ctx, userID).Return(permissions, nil)

			uc := permission.NewGetPermissionsUseCase(permRepo, cacheRepo)

			allowed, err := uc.CheckPermission(ctx, userID, tt.checkPerm)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, allowed)
		})
	}
}
