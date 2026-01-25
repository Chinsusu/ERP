package department

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type MoveDepartmentUseCase struct {
	deptRepo repository.DepartmentRepository
}

func NewMoveDepartmentUseCase(deptRepo repository.DepartmentRepository) *MoveDepartmentUseCase {
	return &MoveDepartmentUseCase{deptRepo: deptRepo}
}

func (uc *MoveDepartmentUseCase) Execute(ctx context.Context, id, newParentID string) (*entity.Department, error) {
	deptID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.BadRequest("invalid department ID")
	}

	dept, err := uc.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return nil, errors.NotFound("department not found")
	}

	var parentPath string = "/"
	if newParentID == "" {
		dept.ParentID = nil
	} else {
		pID, err := uuid.Parse(newParentID)
		if err != nil {
			return nil, errors.BadRequest("invalid parent ID")
		}
		if pID == deptID {
			return nil, errors.BadRequest("cannot set department as its own parent")
		}
		parent, err := uc.deptRepo.GetByID(ctx, pID)
		if err != nil {
			return nil, errors.NotFound("parent department not found")
		}
		dept.ParentID = &pID
		parentPath = parent.Path
	}

	dept.UpdatePath(parentPath)

	if err := uc.deptRepo.Update(ctx, dept); err != nil {
		return nil, errors.Internal(err)
	}

	return dept, nil
}
