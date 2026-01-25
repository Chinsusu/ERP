package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth"
	repoMocks "github.com/erp-cosmetics/auth-service/internal/domain/repository/mocks"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test fixtures
func createTestUser(password string) *entity.User {
	user := &entity.User{
		ID:                  uuid.New(),
		UserID:              uuid.New(),
		Email:               "test@company.vn",
		IsActive:            true,
		EmailVerified:       true,
		FailedLoginAttempts: 0,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	user.SetPassword(password)
	return user
}

func createLockedUser(password string) *entity.User {
	user := createTestUser(password)
	lockTime := time.Now().Add(30 * time.Minute)
	user.LockedUntil = &lockTime
	user.FailedLoginAttempts = 5
	return user
}

func createInactiveUser(password string) *entity.User {
	user := createTestUser(password)
	user.IsActive = false
	return user
}

func createTestRoles() []entity.Role {
	return []entity.Role{
		{
			ID:   uuid.New(),
			Name: "Admin",
		},
	}
}

func createTestPermissions() []entity.Permission {
	return []entity.Permission{
		{
			ID:       uuid.New(),
			Code:     "user:user:read",
			Name:     "Read Users",
			Service:  "user",
			Resource: "user",
			Action:   "read",
		},
		{
			ID:       uuid.New(),
			Code:     "user:user:create",
			Name:     "Create Users",
			Service:  "user",
			Resource: "user",
			Action:   "create",
		},
	}
}

func setupLoginUseCaseMocks() (
	*repoMocks.MockUserRepository,
	*repoMocks.MockRoleRepository,
	*repoMocks.MockPermissionRepository,
	*repoMocks.MockTokenRepository,
	*repoMocks.MockCacheRepository,
	*repoMocks.MockEventPublisher,
) {
	return &repoMocks.MockUserRepository{},
		&repoMocks.MockRoleRepository{},
		&repoMocks.MockPermissionRepository{},
		&repoMocks.MockTokenRepository{},
		&repoMocks.MockCacheRepository{},
		&repoMocks.MockEventPublisher{}
}

// TestLoginUseCase_Execute_Success tests successful login
func TestLoginUseCase_Execute_Success(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)
	roles := createTestRoles()
	permissions := createTestPermissions()

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	// Setup expectations
	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return(roles, nil)
	permRepo.On("GetUserPermissions", ctx, user.UserID).Return(permissions, nil)
	tokenRepo.On("CreateRefreshToken", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
	tokenRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)
	cacheRepo.On("SetUserPermissions", ctx, user.UserID, permissions, mock.AnythingOfType("time.Duration")).Return(nil)
	eventPub.On("Publish", "auth.user.logged_in", mock.Anything).Return(nil)

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "test@company.vn",
		Password:  password,
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, int64(900), resp.ExpiresIn)
	assert.Equal(t, user.Email, resp.User.Email)
	assert.Contains(t, resp.User.Roles, "Admin")

	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
	permRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	cacheRepo.AssertExpectations(t)
	eventPub.AssertExpectations(t)
}

// TestLoginUseCase_Execute_WrongPassword tests login with wrong password
func TestLoginUseCase_Execute_WrongPassword(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "test@company.vn",
		Password:  "WrongPassword",
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Invalid email or password")

	userRepo.AssertExpectations(t)
}

// TestLoginUseCase_Execute_UserNotFound tests login with non-existent email
func TestLoginUseCase_Execute_UserNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "nonexistent@company.vn").Return(nil, errors.New("user not found"))

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "nonexistent@company.vn",
		Password:  "SomePassword",
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Invalid email or password")

	userRepo.AssertExpectations(t)
}

// TestLoginUseCase_Execute_AccountLocked tests login with locked account
func TestLoginUseCase_Execute_AccountLocked(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createLockedUser(password)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "test@company.vn",
		Password:  password,
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Account is locked")

	userRepo.AssertExpectations(t)
}

// TestLoginUseCase_Execute_AccountInactive tests login with inactive account
func TestLoginUseCase_Execute_AccountInactive(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createInactiveUser(password)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "test@company.vn",
		Password:  password,
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Account is disabled")

	userRepo.AssertExpectations(t)
}

// TestLoginUseCase_Execute_AccountLockAfter5Attempts tests account locking after 5 failed attempts
func TestLoginUseCase_Execute_AccountLockAfter5Attempts(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)
	user.FailedLoginAttempts = 4 // Already 4 failed attempts

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "test@company.vn",
		Password:  "WrongPassword", // 5th failed attempt
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Account locked due to too many failed attempts")

	// Verify that failed attempts were incremented
	assert.Equal(t, 5, user.FailedLoginAttempts)
	assert.NotNil(t, user.LockedUntil)

	userRepo.AssertExpectations(t)
}

// TestLoginUseCase_Execute_RolesFetchError tests error handling when fetching roles fails
func TestLoginUseCase_Execute_RolesFetchError(t *testing.T) {
	// Arrange
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return(nil, errors.New("database error"))

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)

	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:     "test@company.vn",
		Password:  password,
		IPAddress: "127.0.0.1",
		UserAgent: "Test Client",
	}

	// Act
	resp, err := uc.Execute(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resp)

	userRepo.AssertExpectations(t)
	roleRepo.AssertExpectations(t)
}

// Table-driven test for password validation scenarios
func TestLoginUseCase_Execute_PasswordValidation(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		inputPassword string
		wantErr       bool
		errContains   string
	}{
		{
			name:          "correct password",
			password:      "Admin@123",
			inputPassword: "Admin@123",
			wantErr:       false,
		},
		{
			name:          "wrong password - case sensitive",
			password:      "Admin@123",
			inputPassword: "admin@123",
			wantErr:       true,
			errContains:   "Invalid email or password",
		},
		{
			name:          "wrong password - completely wrong",
			password:      "Admin@123",
			inputPassword: "WrongPassword!",
			wantErr:       true,
			errContains:   "Invalid email or password",
		},
		{
			name:          "empty password",
			password:      "Admin@123",
			inputPassword: "",
			wantErr:       true,
			errContains:   "Invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			user := createTestUser(tt.password)

			userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

			userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
			userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil).Maybe()

			if !tt.wantErr {
				roles := createTestRoles()
				permissions := createTestPermissions()
				roleRepo.On("GetUserRoles", ctx, user.UserID).Return(roles, nil)
				permRepo.On("GetUserPermissions", ctx, user.UserID).Return(permissions, nil)
				tokenRepo.On("CreateRefreshToken", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(nil)
				tokenRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(nil)
				cacheRepo.On("SetUserPermissions", ctx, user.UserID, permissions, mock.AnythingOfType("time.Duration")).Return(nil)
				eventPub.On("Publish", "auth.user.logged_in", mock.Anything).Return(nil)
			}

			jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
			uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

			req := &auth.LoginRequest{
				Email:     "test@company.vn",
				Password:  tt.inputPassword,
				IPAddress: "127.0.0.1",
				UserAgent: "Test Client",
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
// TestLoginUseCase_Execute_UpdateError tests error handling when updating user fails
func TestLoginUseCase_Execute_UpdateError(t *testing.T) {
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)
	roles := createTestRoles()
	permissions := createTestPermissions()


	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return(roles, nil).Maybe()
	permRepo.On("GetUserPermissions", ctx, user.UserID).Return(permissions, nil).Maybe()

	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:    "test@company.vn",
		Password: password,
	}

	resp, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database error")
}

// TestLoginUseCase_Execute_TokenCreateError tests error handling when creating refresh token fails
func TestLoginUseCase_Execute_TokenCreateError(t *testing.T) {
	ctx := context.Background()
	password := "Admin@123"
	user := createTestUser(password)
	roles := createTestRoles()
	permissions := createTestPermissions()

	userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, eventPub := setupLoginUseCaseMocks()

	userRepo.On("GetByEmail", ctx, "test@company.vn").Return(user, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	roleRepo.On("GetUserRoles", ctx, user.UserID).Return(roles, nil)
	permRepo.On("GetUserPermissions", ctx, user.UserID).Return(permissions, nil)
	tokenRepo.On("CreateRefreshToken", ctx, mock.AnythingOfType("*entity.RefreshToken")).Return(errors.New("token creation error"))
	tokenRepo.On("CreateSession", ctx, mock.AnythingOfType("*entity.Session")).Return(nil).Maybe()
	cacheRepo.On("SetUserPermissions", ctx, user.UserID, permissions, mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
	eventPub.On("Publish", "auth.user.logged_in", mock.Anything).Return(nil).Maybe()


	jwtManager := jwt.NewManager("test-secret-key-32-characters-long", 15*time.Minute, 7*24*time.Hour)
	uc := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, cacheRepo, jwtManager, eventPub)

	req := &auth.LoginRequest{
		Email:    "test@company.vn",
		Password: password,
	}

	resp, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "token creation error")
}
