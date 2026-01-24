package repository

import (
	"context"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entity.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// GetByUserID retrieves a user by user_id (references user-service)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	
	// Update updates a user
	Update(ctx context.Context, user *entity.User) error
	
	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error
	
	// IncrementFailedAttempts increments failed login attempts
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error
	
	// ResetFailedAttempts resets failed login attempts
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	
	// UpdateLastLogin updates the last login timestamp
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}
