package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig holds rate limiter configuration
type RateLimiterConfig struct {
	Enabled       bool
	RequestsPerMin int
	Burst         int
	IPLimitPerMin int
}

// RateLimiter returns rate limiting middleware using Redis
type RateLimiter struct {
	redis  *redis.Client
	config RateLimiterConfig
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(redisClient *redis.Client, config RateLimiterConfig) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		config: config,
	}
}

// Middleware returns the rate limiting middleware
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.config.Enabled {
			c.Next()
			return
		}

		// Determine key and limit based on authentication
		var key string
		var limit int
		
		userID, exists := c.Get("user_id")
		if exists && userID != nil {
			key = fmt.Sprintf("ratelimit:user:%s", userID.(string))
			limit = rl.config.RequestsPerMin
		} else {
			key = fmt.Sprintf("ratelimit:ip:%s", c.ClientIP())
			limit = rl.config.IPLimitPerMin
		}

		// Sliding window rate limiting
		ctx := context.Background()
		now := time.Now().Unix()
		windowStart := now - 60 // 1 minute window

		// Remove old entries
		rl.redis.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))

		// Count current requests
		count, err := rl.redis.ZCard(ctx, key).Result()
		if err != nil {
			// If Redis fails, allow request
			c.Next()
			return
		}

		// Check limit
		if int(count) >= limit {
			resetTime := time.Unix(now+60, 0)
			
			c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Header("Retry-After", "60")

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": 60,
				"limit":       limit,
				"remaining":   0,
				"reset_at":    resetTime.Format(time.RFC3339),
			})
			return
		}

		// Add current request
		rl.redis.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: fmt.Sprintf("%d-%s", now, GetRequestID(c)),
		})
		rl.redis.Expire(ctx, key, 2*time.Minute)

		// Set headers
		remaining := limit - int(count) - 1
		if remaining < 0 {
			remaining = 0
		}
		
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(now+60, 10))

		c.Next()
	}
}
