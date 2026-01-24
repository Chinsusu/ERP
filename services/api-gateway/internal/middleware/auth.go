package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

// AuthConfig holds auth middleware configuration
type AuthConfig struct {
	JWTSecret string
	Redis     *redis.Client
}

// Auth returns authentication middleware
func Auth(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
				"code":  "MISSING_TOKEN",
			})
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  "INVALID_TOKEN_FORMAT",
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
				"code":  "INVALID_TOKEN",
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
				"code":  "INVALID_CLAIMS",
			})
			return
		}

		// Check token expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token has expired",
					"code":  "TOKEN_EXPIRED",
				})
				return
			}
		}

		// Check blacklist in Redis
		if config.Redis != nil {
			jti, _ := claims["jti"].(string)
			if jti != "" {
				blacklistKey := fmt.Sprintf("blacklist:token:%s", jti)
				exists, _ := config.Redis.Exists(context.Background(), blacklistKey).Result()
				if exists > 0 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Token has been revoked",
						"code":  "REVOKED_TOKEN",
					})
					return
				}
			}
		}

		// Extract user info
		userID, _ := claims["sub"].(string)
		if userID == "" {
			userID, _ = claims["user_id"].(string)
		}

		// Set user info in context
		c.Set("user_id", userID)
		c.Set("claims", claims)

		// Add user ID to request headers for downstream services
		c.Request.Header.Set("X-User-ID", userID)

		c.Next()
	}
}

// OptionalAuth validates JWT if present but doesn't require it
func OptionalAuth(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Use Auth middleware logic
		Auth(config)(c)
	}
}
