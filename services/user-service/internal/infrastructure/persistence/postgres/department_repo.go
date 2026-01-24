package postgres

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, dept *entity.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *departmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).
		Preload("Parent").
		Preload("Manager").
		First(&dept, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) GetByCode(ctx context.Context, code string) (*entity.Department, error) {
	var dept entity.Department
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *departmentRepository) Update(ctx context.Context, dept *entity.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Department{}, id).Error
}

func (r *departmentRepository) GetTree(ctx context.Context) ([]entity.Department, error) {
	var departments []entity.Department
	
	// Get all departments ordered by path for hierarchical display
	err := r.db.WithContext(ctx).
		Preload("Manager").
		Order("path ASC").
		Find(&departments).Error
	
	if err != nil {
		return nil, err
	}

	// Build tree structure
	return r.buildTree(departments), nil
}

func (r *departmentRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]entity.Department, error) {
	var children []entity.Department
	err := r.db.WithContext(ctx).
		Preload("Manager").
		Where("parent_id = ?", parentID).
		Order("name ASC").
		Find(&children).Error
	return children, err
}

func (r *departmentRepository) GetUsers(ctx context.Context, deptID uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	err := r.db.WithContext(ctx).
		Where("department_id = ?", deptID).
		Order("first_name, last_name").
		Find(&users).Error
	return users, err
}

// buildTree constructs hierarchical tree from flat list
func (r *departmentRepository) buildTree(departments []entity.Department) []entity.Department {
	deptMap := make(map[uuid.UUID]*entity.Department)
	var roots []entity.Department

	// First pass: create map
	for i := range departments {
		deptMap[departments[i].ID] = &departments[i]
		departments[i].Children = []entity.Department{}
	}

	// Second pass: build tree
	for i := range departments {
		if departments[i].ParentID == nil {
			roots = append(roots, departments[i])
		} else if parent, ok := deptMap[*departments[i].ParentID]; ok {
			parent.Children = append(parent.Children, departments[i])
		}
	}

	return roots
}
