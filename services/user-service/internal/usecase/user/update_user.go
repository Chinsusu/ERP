package user

import (
	"context"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/google/uuid"
)

type UpdateUserUseCase struct {
	userRepo    repository.UserRepository
	profileRepo repository.UserProfileRepository
	authClient  AuthServiceClient
	eventPub    EventPublisher
}

type UpdateUserRequest struct {
	ID           string
	FirstName    string
	LastName     string
	Phone        string
	AvatarURL    string
	DepartmentID *string
	ManagerID    *string
	Status       string
	// Profile info
	Address      *string
}

func NewUpdateUserUseCase(
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
	authClient AuthServiceClient,
	eventPub EventPublisher,
) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		authClient:  authClient,
		eventPub:    eventPub,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, req *UpdateUserRequest) (*entity.User, error) {
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errors.BadRequest("invalid user ID")
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NotFound("user not found")
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	if req.Status != "" {
		oldStatus := user.Status
		user.Status = req.Status
		if oldStatus != req.Status {
			// Notify auth service about status change
			isActive := req.Status == "active"
			_ = uc.authClient.UpdateUserStatus(ctx, user.ID.String(), isActive)
		}
	}

	if req.DepartmentID != nil {
		if *req.DepartmentID == "" {
			user.DepartmentID = nil
		} else {
			deptID, err := uuid.Parse(*req.DepartmentID)
			if err == nil {
				user.DepartmentID = &deptID
			}
		}
	}

	if req.ManagerID != nil {
		if *req.ManagerID == "" {
			user.ManagerID = nil
		} else {
			mgrID, err := uuid.Parse(*req.ManagerID)
			if err == nil {
				user.ManagerID = &mgrID
			}
		}
	}

	if err := user.Validate(); err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, errors.Internal(err)
	}

	// Update profile address if provided
	if req.Address != nil {
		profile, err := uc.profileRepo.GetByUserID(ctx, user.ID)
		if err == nil && profile != nil {
			profile.Address = *req.Address
			_ = uc.profileRepo.Update(ctx, profile)
		}
	}

	// Publish event
	uc.eventPub.Publish("user.updated", map[string]interface{}{
		"user_id": user.ID.String(),
		"status":  user.Status,
	})

	return user, nil
}
