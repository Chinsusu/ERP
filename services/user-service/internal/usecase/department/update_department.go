package department

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type UpdateDepartmentUseCase struct {
	deptRepo repository.DepartmentRepository
}

type UpdateDepartmentRequest struct {
	ID        string
	Name      string
	ManagerID *string
	Status    string
}

func NewUpdateDepartmentUseCase(deptRepo repository.DepartmentRepository) *UpdateDepartmentUseCase {
	return &UpdateDepartmentUseCase{deptRepo: deptRepo}
}

func (uc *UpdateDepartmentUseCase) Execute(ctx context.Context, req *UpdateDepartmentRequest) (*entity.Department, error) {
	deptID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errors.BadRequest("invalid department ID")
	}

	dept, err := uc.deptRepo.GetByID(ctx, deptID)
	if err != nil {
		return nil, errors.NotFound("department not found")
	}

	if req.Name != "" {
		dept.Name = req.Name
	}
	if req.Status != "" {
		dept.Status = req.Status
	}
	if req.ManagerID != nil {
		if *req.ManagerID == "" {
			dept.ManagerID = nil
		} else {
			mgrID, err := uuid.Parse(*req.ManagerID)
			if err == nil {
				dept.ManagerID = &mgrID
			}
		}
	}

	if err := dept.Validate(); err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	if err := uc.deptRepo.Update(ctx, dept); err != nil {
		return nil, errors.Internal(err)
	}

	return dept, nil
}
