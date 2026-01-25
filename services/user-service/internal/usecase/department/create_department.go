package department

import (
	"context"

	"github.com/google/uuid"
	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
)

type CreateDepartmentUseCase struct {
	deptRepo repository.DepartmentRepository
	eventPub EventPublisher
}

type CreateDepartmentRequest struct {
	Code     string
	Name     string
	ParentID *string
}

func NewCreateDepartmentUseCase(
	deptRepo repository.DepartmentRepository,
	eventPub EventPublisher,
) *CreateDepartmentUseCase {
	return &CreateDepartmentUseCase{
		deptRepo: deptRepo,
		eventPub: eventPub,
	}
}

func (uc *CreateDepartmentUseCase) Execute(ctx context.Context, req *CreateDepartmentRequest) (*entity.Department, error) {
	// Validate
	if req.Code == "" || req.Name == "" {
		return nil, errors.BadRequest("code and name are required")
	}

	// Check if code already exists
	existing, _ := uc.deptRepo.GetByCode(ctx, req.Code)
	if existing != nil {
		return nil, errors.Conflict("department code already exists")
	}

	// Create department
	dept := &entity.Department{
		Code:   req.Code,
		Name:   req.Name,
		Status: "active",
	}

	// Get parent if specified
	var parentPath string
	if req.ParentID != nil {
		parentID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, errors.BadRequest("invalid parent_id")
		}

		parent, err := uc.deptRepo.GetByID(ctx, parentID)
		if err != nil {
			return nil, errors.NotFound("parent department not found")
		}

		dept.ParentID = &parentID
		parentPath = parent.Path
	}

	// Update path
	dept.UpdatePath(parentPath)

	// Validate
	if err := dept.Validate(); err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	// Save to database
	if err := uc.deptRepo.Create(ctx, dept); err != nil {
		return nil, errors.Internal(err)
	}

	// Publish event
	uc.eventPub.Publish("department.created", map[string]interface{}{
		"department_id": dept.ID.String(),
		"code":          dept.Code,
		"name":          dept.Name,
	})

	return dept, nil
}
