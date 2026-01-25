package user

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type AssignRoleUseCase struct {
	userRepo   repository.UserRepository
	authClient AuthServiceClient
	eventPub   EventPublisher
}

func NewAssignRoleUseCase(userRepo repository.UserRepository, authClient AuthServiceClient, eventPub EventPublisher) *AssignRoleUseCase {
	return &AssignRoleUseCase{
		userRepo:   userRepo,
		authClient: authClient,
		eventPub:   eventPub,
	}
}

func (uc *AssignRoleUseCase) Execute(ctx context.Context, userID, roleID string) error {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return errors.BadRequest("invalid user ID")
	}

	user, err := uc.userRepo.GetByID(ctx, uID)
	if err != nil {
		return errors.NotFound("user not found")
	}

	if err := uc.authClient.AssignRole(ctx, userID, roleID); err != nil {
		return errors.Internal(err)
	}

	// Publish event
	uc.eventPub.Publish("user.role_assigned", map[string]interface{}{
		"user_id": userID,
		"email":   user.Email,
		"role_id": roleID,
	})

	return nil
}
