package postgres

import (
	"context"

	"github.com/erp-cosmetics/file-service/internal/domain/entity"
	"github.com/erp-cosmetics/file-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileCategoryRepository struct {
	db *gorm.DB
}

// NewFileCategoryRepository creates new repository
func NewFileCategoryRepository(db *gorm.DB) repository.FileCategoryRepository {
	return &fileCategoryRepository{db: db}
}

func (r *fileCategoryRepository) GetByCode(ctx context.Context, code string) (*entity.FileCategory, error) {
	var category entity.FileCategory
	err := r.db.WithContext(ctx).First(&category, "code = ? AND is_active = ?", code, true).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *fileCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.FileCategory, error) {
	var category entity.FileCategory
	err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *fileCategoryRepository) List(ctx context.Context) ([]*entity.FileCategory, error) {
	var categories []*entity.FileCategory
	err := r.db.WithContext(ctx).Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *fileCategoryRepository) ListActive(ctx context.Context) ([]*entity.FileCategory, error) {
	var categories []*entity.FileCategory
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&categories).Error
	return categories, err
}
