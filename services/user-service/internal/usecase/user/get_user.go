package user

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type GetUserUseCase struct {
	userRepo repository.UserRepository
}

func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{userRepo: userRepo}
}

func (uc *GetUserUseCase) Execute(ctx context.Context, userID string) (*entity.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.BadRequest("invalid user ID")
	}

	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NotFound("user not found")
	}

	return user, nil
}
