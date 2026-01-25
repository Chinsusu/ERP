package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/user-service/internal/usecase/user"
	"github.com/erp-cosmetics/user-service/internal/testmocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserUseCaseMocks() (
	*testmocks.MockUserRepository,
	*testmocks.MockUserProfileRepository,
	*testmocks.MockAuthServiceClient,
	*testmocks.MockEventPublisher,
) {
	return &testmocks.MockUserRepository{},
		&testmocks.MockUserProfileRepository{},
		&testmocks.MockAuthServiceClient{},
		&testmocks.MockEventPublisher{}
}

func createTestUser() *entity.User {
	id := uuid.New()
	deptID := uuid.New()
	return &entity.User{
		ID:           id,
		Email:        "test@company.vn",
		FirstName:    "Test",
		LastName:     "User",
		Phone:        "0123456789",
		EmployeeCode: "EMP20260125001",
		DepartmentID: &deptID,
		Status:       "active",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func TestCreateUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, profileRepo, authClient, eventPub := setupUserUseCaseMocks()
	
	req := &user.CreateUserRequest{
		Email:     "new@company.vn",
		FirstName: "New",
		LastName:  "User",
		Password:  "password123",
	}

	userRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)
	userRepo.On("GetNextSequence", ctx, mock.Anything).Return(1, nil)
	userRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	authClient.On("CreateUserCredentials", ctx, mock.Anything, req.Email, req.Password).Return(nil)
	eventPub.On("Publish", "user.created", mock.Anything).Return(nil)

	uc := user.NewCreateUserUseCase(userRepo, profileRepo, authClient, eventPub)
	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Email, res.Email)
	userRepo.AssertExpectations(t)
}

func TestCreateUserUseCase_Execute_DuplicateEmail(t *testing.T) {
	ctx := context.Background()
	userRepo, profileRepo, authClient, eventPub := setupUserUseCaseMocks()
	
	existingUser := createTestUser()
	req := &user.CreateUserRequest{
		Email:     existingUser.Email,
		FirstName: "New",
		LastName:  "User",
	}

	userRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

	uc := user.NewCreateUserUseCase(userRepo, profileRepo, authClient, eventPub)
	res, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Contains(t, err.Error(), "already exists")
}

func TestUpdateUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, profileRepo, authClient, eventPub := setupUserUseCaseMocks()
	
	targetUser := createTestUser()
	req := &user.UpdateUserRequest{
		ID:        targetUser.ID.String(),
		FirstName: "UpdatedName",
		Status:    "inactive",
	}

	userRepo.On("GetByID", ctx, targetUser.ID).Return(targetUser, nil)
	userRepo.On("Update", ctx, mock.AnythingOfType("*entity.User")).Return(nil)
	authClient.On("UpdateUserStatus", ctx, targetUser.ID.String(), false).Return(nil)
	eventPub.On("Publish", "user.updated", mock.Anything).Return(nil)

	uc := user.NewUpdateUserUseCase(userRepo, profileRepo, authClient, eventPub)
	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "UpdatedName", res.FirstName)
	assert.Equal(t, "inactive", res.Status)
}

func TestDeleteUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, _, _, eventPub := setupUserUseCaseMocks()
	
	targetUser := createTestUser()
	userRepo.On("GetByID", ctx, targetUser.ID).Return(targetUser, nil)
	userRepo.On("Delete", ctx, targetUser.ID).Return(nil)
	eventPub.On("Publish", "user.deleted", mock.Anything).Return(nil)

	uc := user.NewDeleteUserUseCase(userRepo, eventPub)
	err := uc.Execute(ctx, targetUser.ID.String())

	assert.NoError(t, err)
}

func TestGetUserUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, _, _, _ := setupUserUseCaseMocks()
	
	targetUser := createTestUser()
	userRepo.On("GetByID", ctx, targetUser.ID).Return(targetUser, nil)

	uc := user.NewGetUserUseCase(userRepo)
	res, err := uc.Execute(ctx, targetUser.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, targetUser.ID, res.ID)
}

func TestListUsersUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, _, _, _ := setupUserUseCaseMocks()
	
	users := []entity.User{*createTestUser()}
	req := &user.ListUsersRequest{Page: 1, PageSize: 10}
	
	userRepo.On("List", ctx, mock.MatchedBy(func(f *repository.UserFilter) bool {
		return f.Page == 1 && f.PageSize == 10
	})).Return(users, int64(1), nil)

	uc := user.NewListUsersUseCase(userRepo)
	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.Total)
	assert.Len(t, res.Users, 1)
}

func TestAssignRoleUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, _, authClient, eventPub := setupUserUseCaseMocks()
	
	targetUser := createTestUser()
	roleID := uuid.New().String()

	userRepo.On("GetByID", ctx, targetUser.ID).Return(targetUser, nil)
	authClient.On("AssignRole", ctx, targetUser.ID.String(), roleID).Return(nil)
	eventPub.On("Publish", "user.role_assigned", mock.Anything).Return(nil)

	uc := user.NewAssignRoleUseCase(userRepo, authClient, eventPub)
	err := uc.Execute(ctx, targetUser.ID.String(), roleID)

	assert.NoError(t, err)
}

func TestRemoveRoleUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	userRepo, _, authClient, eventPub := setupUserUseCaseMocks()
	
	targetUser := createTestUser()
	roleID := uuid.New().String()

	userRepo.On("GetByID", ctx, targetUser.ID).Return(targetUser, nil)
	authClient.On("RemoveRole", ctx, targetUser.ID.String(), roleID).Return(nil)
	eventPub.On("Publish", "user.role_removed", mock.Anything).Return(nil)

	uc := user.NewRemoveRoleUseCase(userRepo, authClient, eventPub)
	err := uc.Execute(ctx, targetUser.ID.String(), roleID)

	assert.NoError(t, err)
}

func TestUpdateUserUseCase_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	userRepo, profileRepo, authClient, eventPub := setupUserUseCaseMocks()
	
	id := uuid.New()
	userRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	uc := user.NewUpdateUserUseCase(userRepo, profileRepo, authClient, eventPub)
	res, err := uc.Execute(ctx, &user.UpdateUserRequest{ID: id.String()})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestUpdateUserUseCase_Execute_InvalidID(t *testing.T) {
	ctx := context.Background()
	uc := user.NewUpdateUserUseCase(nil, nil, nil, nil)
	res, err := uc.Execute(ctx, &user.UpdateUserRequest{ID: "invalid-uuid"})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDeleteUserUseCase_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	userRepo, _, _, eventPub := setupUserUseCaseMocks()
	
	id := uuid.New()
	userRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	uc := user.NewDeleteUserUseCase(userRepo, eventPub)
	err := uc.Execute(ctx, id.String())

	assert.Error(t, err)
}

func TestGetUserUseCase_Execute_InvalidID(t *testing.T) {
	ctx := context.Background()
	uc := user.NewGetUserUseCase(nil)
	res, err := uc.Execute(ctx, "invalid-uuid")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestListUsersUseCase_Execute_RepoError(t *testing.T) {
	ctx := context.Background()
	userRepo, _, _, _ := setupUserUseCaseMocks()
	
	req := &user.ListUsersRequest{Page: 1, PageSize: 10}
	userRepo.On("List", ctx, mock.Anything).Return([]entity.User{}, int64(0), errors.New("db error"))

	uc := user.NewListUsersUseCase(userRepo)
	res, err := uc.Execute(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, res)
}
