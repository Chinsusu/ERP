package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/google/uuid"
)

type LoginUseCase struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permRepo       repository.PermissionRepository
	tokenRepo      repository.TokenRepository
	cacheRepo      repository.CacheRepository
	jwtManager     *jwt.Manager
	eventPublisher EventPublisher
}

type EventPublisher interface {
	Publish(subject string, data interface{}) error
}

type LoginRequest struct {
	Email     string
	Password  string
	IPAddress string
	UserAgent string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *UserInfo
}

type UserInfo struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Email       string
	Roles       []string
	Permissions []string
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	tokenRepo repository.TokenRepository,
	cacheRepo repository.CacheRepository,
	jwtManager *jwt.Manager,
	eventPublisher EventPublisher,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permRepo:       permRepo,
		tokenRepo:      tokenRepo,
		cacheRepo:      cacheRepo,
		jwtManager:     jwtManager,
		eventPublisher: eventPublisher,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.Unauthorized("Invalid email or password")
	}

	// Check if account is locked
	if user.IsLocked() {
		return nil, errors.Forbidden("Account is locked. Please try again later")
	}

	// Check if account is active
	if !user.IsActive {
		return nil, errors.Forbidden("Account is disabled")
	}

	// Verify password
	if !user.VerifyPassword(req.Password) {
		// Increment failed attempts
		user.IncrementFailedAttempts()
		uc.userRepo.Update(ctx, user)
		
		if user.IsLocked() {
			return nil, errors.Forbidden("Account locked due to too many failed attempts")
		}
		
		return nil, errors.Unauthorized("Invalid email or password")
	}

	// Reset failed attempts on successful login
	user.ResetFailedAttempts()
	user.UpdateLastLogin()
	uc.userRepo.Update(ctx, user)

	// Get user roles and permissions
	roles, err := uc.roleRepo.GetUserRoles(ctx, user.UserID)
	if err != nil {
		return nil, errors.Internal(err)
	}

	permissions, err := uc.permRepo.GetUserPermissions(ctx, user.UserID)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Extract role IDs and permission codes
	roleIDs := make([]string, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID.String()
	}

	permCodes := make([]string, len(permissions))
	for i, perm := range permissions {
		permCodes[i] = perm.Code
	}

	// Generate access token
	accessToken, err := uc.jwtManager.GenerateAccessToken(user.UserID.String(), user.Email, roleIDs)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Generate refresh token
	refreshTokenStr, err := uc.jwtManager.GenerateRefreshToken(user.UserID.String(), user.Email)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Hash and store refresh token
	refreshTokenHash := hashToken(refreshTokenStr)
	refreshToken := &entity.RefreshToken{
		UserID:    user.UserID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
	}
	if err := uc.tokenRepo.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, errors.Internal(err)
	}

	// Create session
	claims, _ := uc.jwtManager.VerifyToken(accessToken)
	session := &entity.Session{
		UserID:         user.UserID,
		AccessTokenJTI: claims.ID,
		RefreshTokenID: &refreshToken.ID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		ExpiresAt:      time.Now().Add(15 * time.Minute),
	}
	if err := uc.tokenRepo.CreateSession(ctx, session); err != nil {
		return nil, errors.Internal(err)
	}

	// Cache permissions
	uc.cacheRepo.SetUserPermissions(ctx, user.UserID, permissions, 15*time.Minute)

	// Publish event
	uc.eventPublisher.Publish("auth.user.logged_in", map[string]interface{}{
		"user_id":    user.UserID.String(),
		"email":      user.Email,
		"ip_address": req.IPAddress,
		"timestamp":  time.Now(),
	})

	// Extract role names
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    900, // 15 minutes
		User: &UserInfo{
			ID:          user.ID,
			UserID:      user.UserID,
			Email:       user.Email,
			Roles:       roleNames,
			Permissions: permCodes,
		},
	}, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
