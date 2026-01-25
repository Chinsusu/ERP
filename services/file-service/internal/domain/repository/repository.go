package repository

import (
	"context"

	"github.com/erp-cosmetics/file-service/internal/domain/entity"
	"github.com/google/uuid"
)

// FileRepository interface
type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.File, error)
	GetByAccessToken(ctx context.Context, token string) (*entity.File, error)
	Update(ctx context.Context, file *entity.File) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEntityID(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.File, error)
	GetByCategory(ctx context.Context, category string, limit, offset int) ([]*entity.File, int64, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.File, int64, error)
	DeleteExpired(ctx context.Context) (int64, error)
}

// FileCategoryRepository interface
type FileCategoryRepository interface {
	GetByCode(ctx context.Context, code string) (*entity.FileCategory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.FileCategory, error)
	List(ctx context.Context) ([]*entity.FileCategory, error)
	ListActive(ctx context.Context) ([]*entity.FileCategory, error)
}
