// TODO: Implement repository interfaces
// Reference: services/auth-service/internal/domain/repository/

package repository

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines user data access methods
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByEmployeeCode(ctx context.Context, code string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter *UserFilter) ([]entity.User, int64, error)
	GetNextSequence(ctx context.Context, date string) (int, error)
}

// UserFilter for listing users
type UserFilter struct {
	DepartmentID *uuid.UUID
	Status       string
	Search       string
	Page         int
	PageSize     int
}
