package user

import (
	"context"
)

type AuthServiceClient interface {
	CreateUserCredentials(ctx context.Context, userID, email, password string) error
	UpdateUserStatus(ctx context.Context, userID string, isActive bool) error
	AssignRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
}

type EventPublisher interface {
	Publish(subject string, data interface{}) error
}
