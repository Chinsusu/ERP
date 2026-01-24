package permission

import (
	"context"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type GetPermissionsUseCase struct {
	permRepo  repository.PermissionRepository
	cacheRepo repository.CacheRepository
}

func NewGetPermissionsUseCase(
	permRepo repository.PermissionRepository,
	cacheRepo repository.CacheRepository,
) *GetPermissionsUseCase {
	return &GetPermissionsUseCase{
		permRepo:  permRepo,
		cacheRepo: cacheRepo,
	}
}

func (uc *GetPermissionsUseCase) Execute(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error) {
	// Try cache first
	cachedPerms, err := uc.cacheRepo.GetUserPermissions(ctx, userID)
	if err == nil && cachedPerms != nil {
		return cachedPerms, nil
	}

	// Cache miss - fetch from database
	permissions, err := uc.permRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Update cache
	uc.cacheRepo.SetUserPermissions(ctx, userID, permissions, 15*time.Minute)

	return permissions, nil
}

// CheckPermission checks if user has a specific permission
func (uc *GetPermissionsUseCase) CheckPermission(ctx context.Context, userID uuid.UUID, permissionCode string) (bool, error) {
	permissions, err := uc.Execute(ctx, userID)
	if err != nil {
		return false, err
	}

	// Check for exact match or wildcard match
	for _, perm := range permissions {
		if perm.Matches(permissionCode) {
			return true, nil
		}
	}

	return false, nil
}
