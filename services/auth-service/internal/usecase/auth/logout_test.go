package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth"
	repoMocks "github.com/erp-cosmetics/auth-service/internal/domain/repository/mocks"
	"github.com/erp-cosmetics/shared/pkg/jwt"
)

type LogoutUseCaseTestSuite struct {
	suite.Suite
	ctx            context.Context
	userRepo       *repoMocks.MockUserRepository
	tokenRepo      *repoMocks.MockTokenRepository
	cacheRepo      *repoMocks.MockCacheRepository
	eventPublisher *repoMocks.MockEventPublisher
	jwtManager     *jwt.Manager
	useCase        *auth.LogoutUseCase
}

func (s *LogoutUseCaseTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.userRepo = new(repoMocks.MockUserRepository)
	s.tokenRepo = new(repoMocks.MockTokenRepository)
	s.cacheRepo = new(repoMocks.MockCacheRepository)
	s.eventPublisher = new(repoMocks.MockEventPublisher)
	s.jwtManager = jwt.NewManager("test-secret", time.Minute*15, time.Hour*24)
	s.useCase = auth.NewLogoutUseCase(
		s.userRepo,
		s.tokenRepo,
		s.cacheRepo,
		s.jwtManager,
		s.eventPublisher,
	)
}

func TestLogoutUseCase(t *testing.T) {
	suite.Run(t, new(LogoutUseCaseTestSuite))
}

func (s *LogoutUseCaseTestSuite) TestLogout_Success() {
	userID := uuid.New()
	email := "test@example.com"
	accessToken, _ := s.jwtManager.GenerateAccessToken(userID.String(), email, []string{"role1"})
	refreshToken, _ := s.jwtManager.GenerateRefreshToken(userID.String(), email)

	claims, _ := s.jwtManager.VerifyToken(accessToken)

	// Mock expectations
	s.cacheRepo.On("BlacklistToken", s.ctx, claims.ID, mock.AnythingOfType("time.Time")).Return(nil)
	s.tokenRepo.On("DeleteSession", s.ctx, claims.ID).Return(nil)
	s.cacheRepo.On("DeleteUserPermissions", s.ctx, userID).Return(nil)

	s.userRepo.On("GetByUserID", s.ctx, userID).Return(&entity.User{
		Email: email,
	}, nil)

	s.tokenRepo.On("GetRefreshToken", s.ctx, mock.AnythingOfType("string")).Return(&entity.RefreshToken{
		UserID: userID,
	}, nil)
	s.tokenRepo.On("RevokeRefreshToken", s.ctx, mock.AnythingOfType("string"), userID).Return(nil)

	s.eventPublisher.On("Publish", "auth.user.logged_out", mock.MatchedBy(func(payload map[string]interface{}) bool {
		return payload["user_id"] == userID.String()
	})).Return(nil)

	req := &auth.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	err := s.useCase.Execute(s.ctx, req)

	assert.NoError(s.T(), err)
	s.cacheRepo.AssertExpectations(s.T())
	s.tokenRepo.AssertExpectations(s.T())
	s.eventPublisher.AssertExpectations(s.T())
}

func (s *LogoutUseCaseTestSuite) TestLogout_ExpiredAccessToken() {
	userID := uuid.New()
	email := "test@example.com"
	
	// Create a manager with 0 duration for instant expiry
	shortLivedManager := jwt.NewManager("test-secret", -time.Minute, time.Hour)
	expiredAccessToken, _ := shortLivedManager.GenerateAccessToken(userID.String(), email, nil)
	refreshToken, _ := s.jwtManager.GenerateRefreshToken(userID.String(), email)

	// When access token is expired, VerifyToken returns error, but we still proceed to revoke refresh token
	s.userRepo.On("GetByUserID", s.ctx, userID).Return(&entity.User{
		Email: email,
	}, nil)

	s.tokenRepo.On("GetRefreshToken", s.ctx, mock.AnythingOfType("string")).Return(&entity.RefreshToken{
		UserID: userID,
	}, nil)
	s.tokenRepo.On("RevokeRefreshToken", s.ctx, mock.AnythingOfType("string"), userID).Return(nil)

	s.eventPublisher.On("Publish", "auth.user.logged_out", mock.Anything).Return(nil)

	req := &auth.LogoutRequest{
		AccessToken:  expiredAccessToken,
		RefreshToken: refreshToken,
	}

	err := s.useCase.Execute(s.ctx, req)

	assert.NoError(s.T(), err)
	// Cache repo methods (Blacklist, DeleteUserPermissions) should NOT be called if claims are nil due to error
	s.cacheRepo.AssertNotCalled(s.T(), "BlacklistToken", mock.Anything, mock.Anything, mock.Anything)
	s.tokenRepo.AssertExpectations(s.T())
}

func (s *LogoutUseCaseTestSuite) TestLogout_InvalidRefreshToken() {
	userID := uuid.New()
	email := "test@example.com"
	accessToken, _ := s.jwtManager.GenerateAccessToken(userID.String(), email, nil)
	claims, _ := s.jwtManager.VerifyToken(accessToken)

	s.cacheRepo.On("BlacklistToken", s.ctx, claims.ID, mock.Anything).Return(nil)
	s.tokenRepo.On("DeleteSession", s.ctx, claims.ID).Return(nil)
	s.cacheRepo.On("DeleteUserPermissions", s.ctx, userID).Return(nil)

	// Refresh token not found in DB
	// But userID might be known from access token
	s.userRepo.On("GetByUserID", s.ctx, userID).Return(&entity.User{
		Email: email,
	}, nil)

	s.tokenRepo.On("GetRefreshToken", s.ctx, mock.Anything).Return(nil, assert.AnError)

	s.eventPublisher.On("Publish", "auth.user.logged_out", mock.Anything).Return(nil)

	req := &auth.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: "invalid-token",
	}

	err := s.useCase.Execute(s.ctx, req)

	assert.NoError(s.T(), err) // Logout still returns success even if refresh token revocation fails
	s.cacheRepo.AssertExpectations(s.T())
}

func (s *LogoutUseCaseTestSuite) TestLogout_NoRefreshToken() {
	userID := uuid.New()
	email := "test@example.com"
	accessToken, _ := s.jwtManager.GenerateAccessToken(userID.String(), email, nil)
	claims, _ := s.jwtManager.VerifyToken(accessToken)

	s.cacheRepo.On("BlacklistToken", s.ctx, claims.ID, mock.Anything).Return(nil)
	s.tokenRepo.On("DeleteSession", s.ctx, claims.ID).Return(nil)
	s.cacheRepo.On("DeleteUserPermissions", s.ctx, userID).Return(nil)

	s.userRepo.On("GetByUserID", s.ctx, userID).Return(&entity.User{
		Email: email,
	}, nil)

	s.eventPublisher.On("Publish", "auth.user.logged_out", mock.Anything).Return(nil)

	req := &auth.LogoutRequest{
		AccessToken: accessToken,
	}

	err := s.useCase.Execute(s.ctx, req)

	assert.NoError(s.T(), err)
	s.tokenRepo.AssertNotCalled(s.T(), "GetRefreshToken", mock.Anything, mock.Anything)
}
