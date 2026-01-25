package department

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type DeleteDepartmentUseCase struct {
	deptRepo repository.DepartmentRepository
}

func NewDeleteDepartmentUseCase(deptRepo repository.DepartmentRepository) *DeleteDepartmentUseCase {
	return &DeleteDepartmentUseCase{deptRepo: deptRepo}
}

func (uc *DeleteDepartmentUseCase) Execute(ctx context.Context, id string) error {
	deptID, err := uuid.Parse(id)
	if err != nil {
		return errors.BadRequest("invalid department ID")
	}

	_, err = uc.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return errors.NotFound("department not found")
	}

	// Check if department has users
	users, err := uc.deptRepo.GetUsers(ctx, deptID)
	if err == nil && len(users) > 0 {
		return errors.Conflict("cannot delete department with associated users")
	}

	// Check if department has children
	children, err := uc.deptRepo.GetChildren(ctx, deptID)
	if err == nil && len(children) > 0 {
		return errors.Conflict("cannot delete department with sub-departments")
	}

	if err := uc.deptRepo.Delete(ctx, deptID); err != nil {
		return errors.Internal(err)
	}

	return nil
}
