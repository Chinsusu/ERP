package postgres

import (
	"context"
	"time"

	"github.com/erp-cosmetics/file-service/internal/domain/entity"
	"github.com/erp-cosmetics/file-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates new repository
func NewFileRepository(db *gorm.DB) repository.FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *entity.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.File, error) {
	var file entity.File
	err := r.db.WithContext(ctx).First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetByAccessToken(ctx context.Context, token string) (*entity.File, error) {
	var file entity.File
	err := r.db.WithContext(ctx).First(&file, "access_token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Update(ctx context.Context, file *entity.File) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.File{}, "id = ?", id).Error
}

func (r *fileRepository) GetByEntityID(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.File, error) {
	var files []*entity.File
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&files).Error
	return files, err
}

func (r *fileRepository) GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.File, int64, error) {
	var files []*entity.File
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.File{}).Where("category = ?", category)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&files).Error
	return files, total, err
}

func (r *fileRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.File, int64, error) {
	var files []*entity.File
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.File{}).Where("created_by = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&files).Error
	return files, total, err
}

func (r *fileRepository) DeleteExpired(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).
		Delete(&entity.File{}, "expires_at IS NOT NULL AND expires_at < ?", time.Now())
	return result.RowsAffected, result.Error
}
