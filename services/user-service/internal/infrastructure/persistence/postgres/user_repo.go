package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Manager").
		Preload("Profile").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Department").
		Preload("Manager").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmployeeCode(ctx context.Context, code string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Preload("Department").
		Where("employee_code = ?", code).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, filter *repository.UserFilter) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{})

	// Apply filters
	if filter.DepartmentID != nil {
		query = query.Where("department_id = ?", filter.DepartmentID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR employee_code LIKE ?",
			search, search, search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.PageSize
	err := query.
		Preload("Department").
		Preload("Manager").
		Offset(offset).
		Limit(filter.PageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

func (r *userRepository) GetNextSequence(ctx context.Context, date string) (int, error) {
	var count int64
	prefix := fmt.Sprintf("EMP%s%%", date)
	
	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("employee_code LIKE ?", prefix).
		Count(&count).Error
	
	return int(count) + 1, err
}
