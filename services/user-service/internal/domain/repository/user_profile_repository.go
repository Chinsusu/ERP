// TODO: Implement user profile repository interface

package repository

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/google/uuid"
)

// UserProfileRepository defines user profile data access methods
type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) error
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
}
