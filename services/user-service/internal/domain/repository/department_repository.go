// TODO: Implement department repository interface

package repository

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/google/uuid"
)

// DepartmentRepository defines department data access methods
type DepartmentRepository interface {
	Create(ctx context.Context, dept *entity.Department) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Department, error)
	GetByCode(ctx context.Context, code string) (*entity.Department, error)
	Update(ctx context.Context, dept *entity.Department) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTree(ctx context.Context) ([]entity.Department, error)
	GetChildren(ctx context.Context, parentID uuid.UUID) ([]entity.Department, error)
	GetUsers(ctx context.Context, deptID uuid.UUID) ([]entity.User, error)
}
