package repository

import (
	"context"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

// RoleRepository defines the interface for role data access
type RoleRepository interface {
	// Create creates a new role
	Create(ctx context.Context, role *entity.Role) error
	
	// GetByID retrieves a role by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	
	// GetByName retrieves a role by name
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	
	// List retrieves all roles
	List(ctx context.Context) ([]entity.Role, error)
	
	// Update updates a role
	Update(ctx context.Context, role *entity.Role) error
	
	// Delete soft deletes a role
	Delete(ctx context.Context, id uuid.UUID) error
	
	// GetUserRoles retrieves all roles for a user
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]entity.Role, error)
	
	// AssignRoleToUser assigns a role to a user
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy uuid.UUID) error
	
	// RemoveRoleFromUser removes a role from a user
	RemoveRoleFromUser(ctx context.Context, userID, roleID uuid.UUID) error
}
