package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth/mocks"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRefreshTokenUseCaseMocks() (
	*mocks.MockUserRepository,
	*mocks.MockRoleRepository,
	*mocks.MockPermissionRepository,
	*mocks.MockTokenRepository,
	*mocks.MockCacheRepository,
) {
	return &mocks.MockUserRepository{},
		&mocks.MockRoleRepository{},
		&mocks.MockPermissionRepository{},
		&mocks.MockTokenRepository{},
		&mocks.MockCacheRepository{}
}

func createValidRefreshToken(userID uuid.UUID) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "valid-hash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

func createExpiredRefreshToken(userID uuid.UUID) *entity.RefreshToken {
	return &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "expired-hash",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
		CreatedAt: time.Now().Add(-8 * 24 * time.Hour),
	}
}

func createRevokedRefreshToken(userID uuid.UUID) *entity.RefreshToken {
	revokedAt := time.Now().Add(-1 * time.Hour)
	return &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: "revoked-hash",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		RevokedAt: &revokedAt,
		CreatedAt: time.Now().Add(-1 * 24 * time.Hour),
	}
}

// TestRefreshTokenUseCase_Execute_Success tests successful token refresh
func TestRefreshTokenUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)
	roles := createTestRoles()
	permissions := createTestPermissions()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	// Generate a valid refresh token
	validRefreshTokenStr, _ := jwtManager.GenerateRefreshToken(user.UserID.String(), user.Email)
	storedToken := createValidRefreshToken(user.UserID)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	tokenRepo.On("GetRefreshToken", ctx, mock.AnythingOfType("string")).Return(storedToken, nil)
	tokenRepo.On("RevokeRefreshToken", ctx, mock.AnythingOfType("string"), user.UserID).Return(nil)
	userRepo.On("GetByUserID", ctx, user.UserID).Return(user, nil)
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return(roles, nil)
	tokenRepo.On("CreateRefreshToken", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
	tokenRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)
	permRepo.On("GetUserPermissions", ctx, user.UserID).Return(permissions, nil)
	cacheRepo.On("SetUserPermissions", ctx, user.UserID, permissions, mock.AnythingOfType("time.Duration")).Return(nil)

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: validRefreshTokenStr,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, int64(900), resp.ExpiresIn)
	// New refresh token should be different from old one
	assert.NotEqual(t, validRefreshTokenStr, resp.RefreshToken)

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

// TestRefreshTokenUseCase_Execute_InvalidToken tests refresh with invalid JWT
func TestRefreshTokenUseCase_Execute_InvalidToken(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: "invalid.jwt.token",
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Invalid or expired refresh token")
}

// TestRefreshTokenUseCase_Execute_ExpiredToken tests refresh with expired JWT
func TestRefreshTokenUseCase_Execute_ExpiredToken(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	// Create a JWT manager with very short expiration
	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 1*time.Millisecond, 1*time.Millisecond)

	// Generate a token and wait for it to expire
	user := createTestUser("password")
	expiredToken, _ := jwtManager.GenerateRefreshToken(user.UserID.String(), user.Email)
	time.Sleep(10 * time.Millisecond) // Wait for token to expire

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: expiredToken,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Invalid or expired refresh token")
}

// TestRefreshTokenUseCase_Execute_TokenNotInDB tests refresh with token not found in database
func TestRefreshTokenUseCase_Execute_TokenNotInDB(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	user := createTestUser("password")
	validToken, _ := jwtManager.GenerateRefreshToken(user.UserID.String(), user.Email)

	tokenRepo.On("GetRefreshToken", ctx, mock.AnythingOfType("string")).Return(nil, errors.New("token not found"))

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: validToken,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Invalid refresh token")

	tokenRepo.AssertExpectations(t)
}

// TestRefreshTokenUseCase_Execute_RevokedToken tests refresh with revoked token
func TestRefreshTokenUseCase_Execute_RevokedToken(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	user := createTestUser("password")
	validToken, _ := jwtManager.GenerateRefreshToken(user.UserID.String(), user.Email)
	revokedToken := createRevokedRefreshToken(user.UserID)

	tokenRepo.On("GetRefreshToken", ctx, mock.AnythingOfType("string")).Return(revokedToken, nil)

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: validToken,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Refresh token is expired or revoked")

	tokenRepo.AssertExpectations(t)
}

// TestRefreshTokenUseCase_Execute_UserInactive tests refresh when user is inactive
func TestRefreshTokenUseCase_Execute_UserInactive(t *testing.T) {
	// Arrange
	ctx := context.Background()
	inactiveUser := createInactiveUser("password")
	storedToken := createValidRefreshToken(inactiveUser.UserID)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
	validToken, _ := jwtManager.GenerateRefreshToken(inactiveUser.UserID.String(), inactiveUser.Email)

	tokenRepo.On("GetRefreshToken", ctx, mock.AnythingOfType("string")).Return(storedToken, nil)
	userRepo.On("GetByUserID", ctx, inactiveUser.UserID).Return(inactiveUser, nil)

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: validToken,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "User not found or inactive")

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// TestRefreshTokenUseCase_Execute_UserNotFound tests refresh when user is not found
func TestRefreshTokenUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userID := uuid.New()
	storedToken := createValidRefreshToken(userID)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
	validToken, _ := jwtManager.GenerateRefreshToken(userID.String(), "deleted@company.vn")

	tokenRepo.On("GetRefreshToken", ctx, mock.AnythingOfType("string")).Return(storedToken, nil)
	userRepo.On("GetByUserID", ctx, userID).Return(nil, errors.New("user not found"))

	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{
		RefreshToken: validToken,
		IPAddress:    "127.0.0.1",
		UserAgent:    "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "User not found or inactive")

	tokenRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

// Table-driven tests for refresh token scenarios
func TestRefreshTokenUseCase_Execute_Scenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func(*mocks.MockTokenRepository, *mocks.MockUserRepository, *entity.User)
		wantErr     bool
		errContains string
	}{
		{
			name: "database error fetching token",
			setupMocks: func(tokenRepo *mocks.MockTokenRepository, userRepo *mocks.MockUserRepository, user *entity.User) {
				tokenRepo.On("GetRefreshToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("database error"))
			},
			wantErr:     true,
			errContains: "Invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			user := createTestUser("password")

			userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()
			tt.setupMocks(tokenRepo, userRepo, user)

			jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
			validToken, _ := jwtManager.GenerateRefreshToken(user.UserID.String(), user.Email)

			uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

			req := &auth.RefreshTokenRequest{
				RefreshToken: validToken,
				IPAddress:    "127.0.0.1",
				UserAgent:    "Test Client",
			}

			resp, err := uc.Execute(ctx, req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
// TestRefreshTokenUseCase_Execute_RolesFetchError tests error handling when fetching roles fails
func TestRefreshTokenUseCase_Execute_RolesFetchError(t *testing.T) {
	ctx := context.Background()
	user := createTestUser("Admin@123")
	refreshToken := createValidRefreshToken(user.UserID)
	refreshTokenStr, _ := jwt.NewManager("test-secret", time.Minute, time.Hour).GenerateRefreshToken(user.UserID.String(), user.Email)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	tokenRepo.On("GetRefreshToken", ctx, mock.Anything).Return(refreshToken, nil)
	userRepo.On("GetByUserID", ctx, user.UserID).Return(user, nil)
	tokenRepo.On("RevokeRefreshToken", ctx, mock.Anything, user.UserID).Return(nil)
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return(nil, errors.New("roles fetch error"))

	jwtManager := jwt.NewManager("test-secret", time.Minute, time.Hour)
	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{RefreshToken: refreshTokenStr}
	resp, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "roles fetch error")
}

// TestRefreshTokenUseCase_Execute_TokenCreateError tests error handling when creating new refresh token fails
func TestRefreshTokenUseCase_Execute_TokenCreateError(t *testing.T) {
	ctx := context.Background()
	user := createTestUser("Admin@123")
	refreshToken := createValidRefreshToken(user.UserID)
	refreshTokenStr, _ := jwt.NewManager("test-secret", time.Minute, time.Hour).GenerateRefreshToken(user.UserID.String(), user.Email)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo := setupRefreshTokenUseCaseMocks()

	tokenRepo.On("GetRefreshToken", ctx, mock.Anything).Return(refreshToken, nil)
	userRepo.On("GetByUserID", ctx, user.UserID).Return(user, nil)
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return([]entity.Role{{Name: "Admin"}}, nil)
	tokenRepo.On("RevokeRefreshToken", ctx, mock.Anything, user.UserID).Return(nil)
	tokenRepo.On("CreateRefreshToken", ctx, mock.Anything).Return(errors.New("token creation error"))

	jwtManager := jwt.NewManager("test-secret", time.Minute, time.Hour)
	uc := auth.NewRefreshTokenUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager)

	req := &auth.RefreshTokenRequest{RefreshToken: refreshTokenStr}
	resp, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "token creation error")
}
