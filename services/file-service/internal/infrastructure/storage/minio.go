package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// MinIOClient wraps MinIO client with helper methods
type MinIOClient struct {
	client *minio.Client
	logger *zap.Logger
}

// Config holds MinIO configuration
type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Region          string
}

// NewMinIOClient creates a new MinIO client
func NewMinIOClient(cfg *Config, logger *zap.Logger) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &MinIOClient{
		client: client,
		logger: logger,
	}, nil
}

// EnsureBucket creates bucket if it doesn't exist
func (m *MinIOClient) EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		m.logger.Info("Created bucket", zap.String("bucket", bucketName))
	}

	return nil
}

// Upload uploads a file to MinIO
func (m *MinIOClient) Upload(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (*UploadResult, error) {
	// Ensure bucket exists
	if err := m.EnsureBucket(ctx, bucketName); err != nil {
		return nil, err
	}

	info, err := m.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	m.logger.Info("File uploaded",
		zap.String("bucket", bucketName),
		zap.String("object", objectName),
		zap.Int64("size", info.Size),
	)

	return &UploadResult{
		Bucket:   bucketName,
		Key:      objectName,
		Size:     info.Size,
		ETag:     info.ETag,
		Location: fmt.Sprintf("%s/%s", bucketName, objectName),
	}, nil
}

// Download downloads a file from MinIO
func (m *MinIOClient) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return object, nil
}

// Delete deletes a file from MinIO
func (m *MinIOClient) Delete(ctx context.Context, bucketName, objectName string) error {
	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	m.logger.Info("File deleted",
		zap.String("bucket", bucketName),
		zap.String("object", objectName),
	)

	return nil
}

// GetPresignedURL generates a presigned URL for downloading
func (m *MinIOClient) GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// GetObjectInfo gets object metadata
func (m *MinIOClient) GetObjectInfo(ctx context.Context, bucketName, objectName string) (*minio.ObjectInfo, error) {
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// UploadResult represents upload result
type UploadResult struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Size     int64  `json:"size"`
	ETag     string `json:"etag"`
	Location string `json:"location"`
}

// InitDefaultBuckets creates default buckets
func (m *MinIOClient) InitDefaultBuckets(ctx context.Context) error {
	buckets := []string{
		"documents",
		"images",
		"certificates",
		"contracts",
		"reports",
		"avatars",
		"products",
		"qc-photos",
		"signatures",
		"attachments",
	}

	for _, bucket := range buckets {
		if err := m.EnsureBucket(ctx, bucket); err != nil {
			m.logger.Error("Failed to create bucket", zap.String("bucket", bucket), zap.Error(err))
		}
	}

	return nil
}
