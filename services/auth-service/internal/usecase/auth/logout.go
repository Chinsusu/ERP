package auth

import (
	"context"

	"github.com/erp-cosmetics/auth-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/google/uuid"
)

type LogoutUseCase struct {
	userRepo       repository.UserRepository
	tokenRepo      repository.TokenRepository
	cacheRepo      repository.CacheRepository
	jwtManager     *jwt.Manager
	eventPublisher EventPublisher
}

type LogoutRequest struct {
	RefreshToken string
	AccessToken  string
}

func NewLogoutUseCase(
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	cacheRepo repository.CacheRepository,
	jwtManager *jwt.Manager,
	eventPublisher EventPublisher,
) *LogoutUseCase {
	return &LogoutUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		cacheRepo:      cacheRepo,
		jwtManager:     jwtManager,
		eventPublisher: eventPublisher,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, req *LogoutRequest) error {
	// Verify access token to get user ID
	claims, err := uc.jwtManager.VerifyToken(req.AccessToken)
	if err != nil {
		// Token might be expired, but we still want to revoke refresh token
		// So we don't return error here
	}

	var userID uuid.UUID
	if claims != nil {
		userID, _ = uuid.Parse(claims.UserID)
		
		// Blacklist access token
		uc.cacheRepo.BlacklistToken(ctx, claims.ID, claims.ExpiresAt.Time)
		
		// Delete session
		uc.tokenRepo.DeleteSession(ctx, claims.ID)
		
		// Clear cached permissions
		uc.cacheRepo.DeleteUserPermissions(ctx, userID)
	}

	// Revoke refresh token
	if req.RefreshToken != "" {
		refreshTokenHash := hashToken(req.RefreshToken)
		refreshToken, err := uc.tokenRepo.GetRefreshToken(ctx, refreshTokenHash)
		if err == nil && refreshToken != nil {
			userID = refreshToken.UserID
			uc.tokenRepo.RevokeRefreshToken(ctx, refreshTokenHash, userID)
		}
	}

	// Publish event
	if userID != uuid.Nil {
		user, _ := uc.userRepo.GetByUserID(ctx, userID)
		if user != nil {
			uc.eventPublisher.Publish("auth.user.logged_out", map[string]interface{}{
				"user_id": userID.String(),
				"email":   user.Email,
				"reason":  "user_logout",
			})
		}
	}

	return nil
}
