package user

import (
	"context"
	"time"

	"github.com/erp-cosmetics/user-service/internal/domain/entity"
	"github.com/erp-cosmetics/user-service/internal/domain/repository"
	"github.com/erp-cosmetics/shared/pkg/errors"
)

type CreateUserUseCase struct {
	userRepo    repository.UserRepository
	profileRepo repository.UserProfileRepository
	authClient  AuthServiceClient
	eventPub    EventPublisher
}

type CreateUserRequest struct {
	Email        string
	Password     string
	FirstName    string
	LastName     string
	Phone        string
	DepartmentID *string
	ManagerID    *string
	DateOfBirth  *time.Time
	Address      string
	JoinDate     *time.Time
}

func NewCreateUserUseCase(
	userRepo repository.UserRepository,
	profileRepo repository.UserProfileRepository,
	authClient AuthServiceClient,
	eventPub EventPublisher,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		authClient:  authClient,
		eventPub:    eventPub,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req *CreateUserRequest) (*entity.User, error) {
	// Validate input
	if req.Email == "" || req.FirstName == "" || req.LastName == "" {
		return nil, errors.BadRequest("email, first name, and last name are required")
	}

	// Check if email already exists
	existing, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.Conflict("email already exists")
	}

	// Create user entity
	user := &entity.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Status:    "active",
	}

	// Set department and manager if provided
	if req.DepartmentID != nil {
		// Parse and set department ID
	}
	if req.ManagerID != nil {
		// Parse and set manager ID
	}

	// Generate employee code
	dateStr := time.Now().Format("20060102")
	sequence, err := uc.userRepo.GetNextSequence(ctx, dateStr)
	if err != nil {
		return nil, errors.Internal(err)
	}
	user.GenerateEmployeeCode(sequence)

	// Validate user
	if err := user.Validate(); err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	// Create user in database
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, errors.Internal(err)
	}

	// Create user profile if additional info provided
	if req.DateOfBirth != nil || req.Address != "" || req.JoinDate != nil {
		profile := &entity.UserProfile{
			UserID:      user.ID,
			DateOfBirth: req.DateOfBirth,
			Address:     req.Address,
			JoinDate:    req.JoinDate,
		}
		if err := uc.profileRepo.Create(ctx, profile); err != nil {
			// Log error but don't fail the entire operation
		}
	}

	// Create credentials in Auth Service
	if req.Password != "" {
		if err := uc.authClient.CreateUserCredentials(ctx, user.ID.String(), user.Email, req.Password); err != nil {
			// Log error but don't fail - can be retried later
		}
	}

	// Publish event
	uc.eventPub.Publish("user.created", map[string]interface{}{
		"user_id":       user.ID.String(),
		"email":         user.Email,
		"employee_code": user.EmployeeCode,
		"created_at":    user.CreatedAt,
	})

	return user, nil
}
