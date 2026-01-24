package user

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
)

type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

type ListUsersRequest struct {
	DepartmentID *string
	Status       string
	Search       string
	Page         int
	PageSize     int
}

type ListUsersResponse struct {
	Users []entity.User `json:"users"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Pages int           `json:"pages"`
}

func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{userRepo: userRepo}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Build filter
	filter := &repository.UserFilter{
		Status:   req.Status,
		Search:   req.Search,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Get users
	users, total, err := uc.userRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate pages
	pages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		pages++
	}

	return &ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Pages: pages,
	}, nil
}
