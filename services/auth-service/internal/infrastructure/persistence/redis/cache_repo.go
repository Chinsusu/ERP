package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type cacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(client *redis.Client) repository.CacheRepository {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) SetUserPermissions(ctx context.Context, userID uuid.UUID, permissions []entity.Permission, ttl time.Duration) error {
	key := fmt.Sprintf("user:permissions:%s", userID.String())
	
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *cacheRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, error) {
	key := fmt.Sprintf("user:permissions:%s", userID.String())
	
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}
	
	var permissions []entity.Permission
	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, err
	}
	
	return permissions, nil
}

func (r *cacheRepository) DeleteUserPermissions(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("user:permissions:%s", userID.String())
	return r.client.Del(ctx, key).Err()
}

func (r *cacheRepository) BlacklistToken(ctx context.Context, jti string, expiresAt time.Time) error {
	key := fmt.Sprintf("token:blacklist:%s", jti)
	ttl := time.Until(expiresAt)
	
	if ttl <= 0 {
		return nil // Token already expired, no need to blacklist
	}
	
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *cacheRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", jti)
	
	_, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil // Not blacklisted
		}
		return false, err
	}
	
	return true, nil // Blacklisted
}
