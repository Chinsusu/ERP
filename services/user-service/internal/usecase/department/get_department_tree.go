package department

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
)

type GetDepartmentTreeUseCase struct {
	deptRepo repository.DepartmentRepository
}

func NewGetDepartmentTreeUseCase(deptRepo repository.DepartmentRepository) *GetDepartmentTreeUseCase {
	return &GetDepartmentTreeUseCase{deptRepo: deptRepo}
}

func (uc *GetDepartmentTreeUseCase) Execute(ctx context.Context) ([]entity.Department, error) {
	return uc.deptRepo.GetTree(ctx)
}
