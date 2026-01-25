package auth

import (
	"context"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/google/uuid"
)

type RefreshTokenUseCase struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permRepo       repository.PermissionRepository
	tokenRepo      repository.TokenRepository
	cacheRepo      repository.CacheRepository
	jwtManager     *jwt.Manager
}

type RefreshTokenRequest struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

func NewRefreshTokenUseCase(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permRepo repository.PermissionRepository,
	tokenRepo repository.TokenRepository,
	cacheRepo repository.CacheRepository,
	jwtManager *jwt.Manager,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		permRepo:   permRepo,
		tokenRepo:  tokenRepo,
		cacheRepo:  cacheRepo,
		jwtManager: jwtManager,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// Verify refresh token JWT
	claims, err := uc.jwtManager.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, errors.Unauthorized("Invalid or expired refresh token")
	}

	// Check in database
	refreshTokenHash := hashToken(req.RefreshToken)
	storedToken, err := uc.tokenRepo.GetRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		return nil, errors.Unauthorized("Invalid refresh token")
	}

	// Validate token
	if !storedToken.IsValid() {
		return nil, errors.Unauthorized("Refresh token is expired or revoked")
	}

	// Get user
	userID, _ := uuid.Parse(claims.UserID)
	user, err := uc.userRepo.GetByUserID(ctx, userID)
	if err != nil || !user.IsActive {
		return nil, errors.Unauthorized("User not found or inactive")
	}

	// Revoke old refresh token (token rotation)
	if err := uc.tokenRepo.RevokeRefreshToken(ctx, refreshTokenHash, userID); err != nil {
		return nil, errors.Internal(err)
	}

	// Get user roles
	roles, err := uc.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, errors.Internal(err)
	}

	roleIDs := make([]string, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID.String()
	}

	// Generate new access token
	newAccessToken, err := uc.jwtManager.GenerateAccessToken(userID.String(), user.Email, roleIDs)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Generate new refresh token
	newRefreshTokenStr, err := uc.jwtManager.GenerateRefreshToken(userID.String(), user.Email)
	if err != nil {
		return nil, errors.Internal(err)
	}

	// Store new refresh token
	newRefreshTokenHash := hashToken(newRefreshTokenStr)
	newRefreshToken := &entity.RefreshToken{
		UserID:    userID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
	}
	if err := uc.tokenRepo.CreateRefreshToken(ctx, newRefreshToken); err != nil {
		return nil, errors.Internal(err)
	}

	// Create new session
	newClaims, _ := uc.jwtManager.VerifyToken(newAccessToken)
	session := &entity.Session{
		UserID:         userID,
		AccessTokenJTI: newClaims.ID,
		RefreshTokenID: &newRefreshToken.ID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		ExpiresAt:      time.Now().Add(15 * time.Minute),
	}
	if err := uc.tokenRepo.CreateSession(ctx, session); err != nil {
		return nil, errors.Internal(err)
	}

	// Update cached permissions
	permissions, _ := uc.permRepo.GetUserPermissions(ctx, userID)
	uc.cacheRepo.SetUserPermissions(ctx, userID, permissions, 15*time.Minute)

	return &RefreshTokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshTokenStr,
		ExpiresIn:    900,
	}, nil
}
