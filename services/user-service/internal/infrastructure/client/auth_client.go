package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient provides methods to interact with Auth Service
type AuthClient struct {
	conn *grpc.ClientConn
	// authpb.AuthServiceClient will be added when proto is generated
}

// NewAuthClient creates a new Auth Service client
func NewAuthClient(authServiceAddr string) (*AuthClient, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		conn: conn,
	}, nil
}

// Close closes the gRPC connection
func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// CreateUserCredentials creates user credentials in Auth Service
func (c *AuthClient) CreateUserCredentials(ctx context.Context, userID, email, password string) error {
	// TODO: Implement when auth.proto is available
	// client := authpb.NewAuthServiceClient(c.conn)
	// _, err := client.CreateUser(ctx, &authpb.CreateUserRequest{
	// 	UserId:   userID,
	// 	Email:    email,
	// 	Password: password,
	// })
	// return err
	return nil
}

// AssignRole assigns a role to user in Auth Service
func (c *AuthClient) AssignRole(ctx context.Context, userID, roleID string) error {
	// TODO: Implement when auth.proto is available
	return nil
}

// RemoveRole removes a role from user in Auth Service
func (c *AuthClient) RemoveRole(ctx context.Context, userID, roleID string) error {
	// TODO: Implement when auth.proto is available
	return nil
}

// GetUserRoles gets user roles from Auth Service
func (c *AuthClient) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	// TODO: Implement when auth.proto is available
	return nil, nil
}

// UpdateUserStatus updates user status in Auth Service
func (c *AuthClient) UpdateUserStatus(ctx context.Context, userID string, isActive bool) error {
	// TODO: Implement when auth.proto is available
	return nil
}
