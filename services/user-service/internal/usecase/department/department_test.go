package department_test

import (
	"context"
	"errors"
	"testing"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/usecase/department"
	"github.com/erp-cosmetics/user-service/internal/testmocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDeptUseCaseMocks() (
	*testmocks.MockDepartmentRepository,
	*testmocks.MockEventPublisher,
) {
	return &testmocks.MockDepartmentRepository{},
		&testmocks.MockEventPublisher{}
}

func createTestDept() *entity.Department {
	id := uuid.New()
	return &entity.Department{
		ID:     id,
		Code:   "DEPT001",
		Name:   "Test Department",
		Level:  0,
		Path:   "/DEPT001/",
		Status: "active",
	}
}

func TestCreateDepartmentUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	deptRepo, eventPub := setupDeptUseCaseMocks()
	
	req := &department.CreateDepartmentRequest{
		Code: "NEWDEPT",
		Name: "New Department",
	}

	deptRepo.On("GetByCode", ctx, req.Code).Return(nil, nil)
	deptRepo.On("Create", ctx, mock.AnythingOfType("*entity.Department")).Return(nil)
	eventPub.On("Publish", "department.created", mock.Anything).Return(nil)

	uc := department.NewCreateDepartmentUseCase(deptRepo, eventPub)
	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Code, res.Code)
}

func TestCreateDepartmentUseCase_Execute_WithParent(t *testing.T) {
	ctx := context.Background()
	deptRepo, eventPub := setupDeptUseCaseMocks()
	
	parentDept := createTestDept()
	parentIDStr := parentDept.ID.String()
	
	req := &department.CreateDepartmentRequest{
		Code:     "SUBDEPT",
		Name:     "Sub Department",
		ParentID: &parentIDStr,
	}

	deptRepo.On("GetByCode", ctx, req.Code).Return(nil, nil)
	deptRepo.On("GetByID", ctx, parentDept.ID).Return(parentDept, nil)
	deptRepo.On("Create", ctx, mock.MatchedBy(func(d *entity.Department) bool {
		return d.Code == req.Code && d.ParentID != nil && *d.ParentID == parentDept.ID && d.Level == 1
	})).Return(nil)
	eventPub.On("Publish", "department.created", mock.Anything).Return(nil)

	uc := department.NewCreateDepartmentUseCase(deptRepo, eventPub)
	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Level)
	assert.Contains(t, res.Path, parentDept.Code)
	deptRepo.AssertExpectations(t)
}

func TestUpdateDepartmentUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	targetDept := createTestDept()
	req := &department.UpdateDepartmentRequest{
		ID:   targetDept.ID.String(),
		Name: "Updated Dept Name",
	}

	deptRepo.On("GetByID", ctx, targetDept.ID).Return(targetDept, nil)
	deptRepo.On("Update", ctx, mock.AnythingOfType("*entity.Department")).Return(nil)

	uc := department.NewUpdateDepartmentUseCase(deptRepo)
	res, err := uc.Execute(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, "Updated Dept Name", res.Name)
}

func TestDeleteDepartmentUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	targetDept := createTestDept()
	deptRepo.On("GetByID", ctx, targetDept.ID).Return(targetDept, nil)
	deptRepo.On("GetUsers", ctx, targetDept.ID).Return([]entity.User{}, nil)
	deptRepo.On("GetChildren", ctx, targetDept.ID).Return([]entity.Department{}, nil)
	deptRepo.On("Delete", ctx, targetDept.ID).Return(nil)

	uc := department.NewDeleteDepartmentUseCase(deptRepo)
	err := uc.Execute(ctx, targetDept.ID.String())

	assert.NoError(t, err)
}

func TestDeleteDepartmentUseCase_Execute_HasUsers(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	targetDept := createTestDept()
	deptRepo.On("GetByID", ctx, targetDept.ID).Return(targetDept, nil)
	deptRepo.On("GetUsers", ctx, targetDept.ID).Return([]entity.User{{ID: uuid.New()}}, nil)

	uc := department.NewDeleteDepartmentUseCase(deptRepo)
	err := uc.Execute(ctx, targetDept.ID.String())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "with associated users")
}

func TestMoveDepartmentUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	targetDept := createTestDept()
	parentDept := createTestDept()
	parentDept.ID = uuid.New()
	parentDept.Code = "PARENT"
	parentDept.Path = "/PARENT/"

	deptRepo.On("GetByID", ctx, targetDept.ID).Return(targetDept, nil)
	deptRepo.On("GetByID", ctx, parentDept.ID).Return(parentDept, nil)
	deptRepo.On("Update", ctx, mock.AnythingOfType("*entity.Department")).Return(nil)

	uc := department.NewMoveDepartmentUseCase(deptRepo)
	res, err := uc.Execute(ctx, targetDept.ID.String(), parentDept.ID.String())

	assert.NoError(t, err)
	assert.Equal(t, &parentDept.ID, res.ParentID)
	assert.Contains(t, res.Path, "/PARENT/")
	assert.Equal(t, 1, res.Level)
}

func TestGetDepartmentTreeUseCase_Execute_Success(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	tree := []entity.Department{*createTestDept()}
	deptRepo.On("GetTree", ctx).Return(tree, nil)

	uc := department.NewGetDepartmentTreeUseCase(deptRepo)
	res, err := uc.Execute(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
}

func TestUpdateDepartmentUseCase_Execute_NotFound(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	id := uuid.New()
	deptRepo.On("GetByID", ctx, id).Return(nil, errors.New("not found"))

	uc := department.NewUpdateDepartmentUseCase(deptRepo)
	res, err := uc.Execute(ctx, &department.UpdateDepartmentRequest{ID: id.String()})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDeleteDepartmentUseCase_Execute_HasChildren(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	targetDept := createTestDept()
	deptRepo.On("GetByID", ctx, targetDept.ID).Return(targetDept, nil)
	deptRepo.On("GetUsers", ctx, targetDept.ID).Return([]entity.User{}, nil)
	deptRepo.On("GetChildren", ctx, targetDept.ID).Return([]entity.Department{{ID: uuid.New()}}, nil)

	uc := department.NewDeleteDepartmentUseCase(deptRepo)
	err := uc.Execute(ctx, targetDept.ID.String())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "with sub-departments")
}

func TestMoveDepartmentUseCase_Execute_SelfParent(t *testing.T) {
	ctx := context.Background()
	deptRepo, _ := setupDeptUseCaseMocks()
	
	targetDept := createTestDept()
	deptRepo.On("GetByID", ctx, targetDept.ID).Return(targetDept, nil)

	uc := department.NewMoveDepartmentUseCase(deptRepo)
	res, err := uc.Execute(ctx, targetDept.ID.String(), targetDept.ID.String())

	assert.Error(t, err)
	assert.Nil(t, res)
}
