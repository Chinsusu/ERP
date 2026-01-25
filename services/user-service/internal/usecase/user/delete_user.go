package user

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type DeleteUserUseCase struct {
	userRepo repository.UserRepository
	eventPub EventPublisher
}

func NewDeleteUserUseCase(userRepo repository.UserRepository, eventPub EventPublisher) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
		eventPub: eventPub,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.BadRequest("invalid user ID")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.NotFound("user not found")
	}

	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return errors.Internal(err)
	}

	// Publish event
	uc.eventPub.Publish("user.deleted", map[string]interface{}{
		"user_id": id,
		"email":   user.Email,
	})

	return nil
}
