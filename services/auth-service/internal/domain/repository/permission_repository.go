package repository

import (
	"context"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

// PermissionRepository defines the interface for permission data access
type PermissionRepository interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *entity.Permission) error
	
	// GetByID retrieves a permission by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	
	// GetByCode retrieves a permission by code
	GetByCode(ctx context.Context, code string) (*entity.Permission, error)
	
	// List retrieves all permissions
	List(ctx context.Context) ([]entity.Permission, error)
	
	// ListByService retrieves permissions for a specific service
	ListByService(ctx context.Context, service string) ([]entity.Permission, error)
	
	// Update updates a permission
	Update(ctx context.Context, permission *entity.Permission) error
	
	// Delete deletes a permission
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetRolePermissions retrieves all permissions for a role
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
	
	// GetUserPermissions retrieves all permissions for a user (via roles)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error)
	
	// AssignPermissionToRole assigns a permission to a role
	AssignPermissionToRole(ctx context.Context, roleID, permissionID uuid.UUID) error
	
	// RemovePermissionFromRole removes a permission from a role
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID uuid.UUID) error
}
