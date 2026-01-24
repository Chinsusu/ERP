package postgres

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) repository.UserProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *userProfileRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*entity.UserProfile, error) {
	var profile entity.UserProfile
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}
