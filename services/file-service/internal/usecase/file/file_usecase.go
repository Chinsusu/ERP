package file

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/erp-cosmetics/file-service/internal/domain/entity"
	"github.com/erp-cosmetics/file-service/internal/domain/repository"
	"github.com/erp-cosmetics/file-service/internal/infrastructure/storage"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines file use case interface
type UseCase interface {
	Upload(ctx context.Context, input *UploadInput) (*entity.File, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.File, error)
	GetDownloadURL(ctx context.Context, id uuid.UUID) (string, error)
	Download(ctx context.Context, id uuid.UUID) (io.ReadCloser, *entity.File, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.File, error)
	ListCategories(ctx context.Context) ([]*entity.FileCategory, error)
}

// UploadInput for uploading a file
type UploadInput struct {
	FileName    string
	ContentType string
	Size        int64
	Reader      io.Reader
	Category    string
	EntityType  string
	EntityID    *uuid.UUID
	IsPublic    bool
	Metadata    map[string]interface{}
	CreatedBy   *uuid.UUID
	ExpiresAt   *time.Time
}

type useCase struct {
	fileRepo     repository.FileRepository
	categoryRepo repository.FileCategoryRepository
	storage      *storage.MinIOClient
	logger       *zap.Logger
}

// NewUseCase creates new file use case
func NewUseCase(
	fileRepo repository.FileRepository,
	categoryRepo repository.FileCategoryRepository,
	storageClient *storage.MinIOClient,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		fileRepo:     fileRepo,
		categoryRepo: categoryRepo,
		storage:      storageClient,
		logger:       logger,
	}
}

func (uc *useCase) Upload(ctx context.Context, input *UploadInput) (*entity.File, error) {
	// Get category for validation
	category, err := uc.categoryRepo.GetByCode(ctx, input.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Validate file
	if err := category.ValidateFile(input.FileName, input.Size); err != nil {
		return nil, err
	}

	// Generate stored name
	fileID := uuid.New()
	ext := filepath.Ext(input.FileName)
	storedName := fileID.String() + ext
	objectPath := time.Now().Format("2006/01/02") + "/" + storedName

	// Calculate checksum (read into buffer for dual use)
	hasher := sha256.New()
	teeReader := io.TeeReader(input.Reader, hasher)

	// Upload to storage
	result, err := uc.storage.Upload(ctx, category.StorageBucket, objectPath, teeReader, input.Size, input.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))

	// Create file record
	file := &entity.File{
		ID:           fileID,
		OriginalName: input.FileName,
		StoredName:   storedName,
		ContentType:  input.ContentType,
		FileSize:     result.Size,
		BucketName:   category.StorageBucket,
		ObjectPath:   objectPath,
		Category:     input.Category,
		EntityType:   input.EntityType,
		EntityID:     input.EntityID,
		IsPublic:     input.IsPublic,
		Checksum:     checksum,
		CreatedBy:    input.CreatedBy,
		ExpiresAt:    input.ExpiresAt,
	}

	if err := uc.fileRepo.Create(ctx, file); err != nil {
		// Cleanup uploaded file
		uc.storage.Delete(ctx, category.StorageBucket, objectPath)
		return nil, fmt.Errorf("failed to save file record: %w", err)
	}

	uc.logger.Info("File uploaded",
		zap.String("file_id", file.ID.String()),
		zap.String("filename", input.FileName),
		zap.Int64("size", result.Size),
	)

	return file, nil
}

func (uc *useCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.File, error) {
	return uc.fileRepo.GetByID(ctx, id)
}

func (uc *useCase) GetDownloadURL(ctx context.Context, id uuid.UUID) (string, error) {
	file, err := uc.fileRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}

	if file.IsExpired() {
		return "", fmt.Errorf("file has expired")
	}

	// Generate presigned URL (1 hour expiry)
	url, err := uc.storage.GetPresignedURL(ctx, file.BucketName, file.ObjectPath, time.Hour)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (uc *useCase) Download(ctx context.Context, id uuid.UUID) (io.ReadCloser, *entity.File, error) {
	file, err := uc.fileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	if file.IsExpired() {
		return nil, nil, fmt.Errorf("file has expired")
	}

	reader, err := uc.storage.Download(ctx, file.BucketName, file.ObjectPath)
	if err != nil {
		return nil, nil, err
	}

	return reader, file, nil
}

func (uc *useCase) Delete(ctx context.Context, id uuid.UUID) error {
	file, err := uc.fileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete from storage
	if err := uc.storage.Delete(ctx, file.BucketName, file.ObjectPath); err != nil {
		uc.logger.Error("Failed to delete from storage", zap.Error(err))
	}

	// Delete record
	return uc.fileRepo.Delete(ctx, id)
}

func (uc *useCase) GetByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]*entity.File, error) {
	return uc.fileRepo.GetByEntityID(ctx, entityType, entityID)
}

func (uc *useCase) ListCategories(ctx context.Context) ([]*entity.FileCategory, error) {
	return uc.categoryRepo.ListActive(ctx)
}
