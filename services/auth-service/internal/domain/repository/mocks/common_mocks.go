package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockEventPublisher is a mock implementation of auth.EventPublisher
type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) Publish(subject string, data interface{}) error {
	args := m.Called(subject, data)
	return args.Error(0)
}

// MockAuthClient is a mock implementation of the requested AuthClient interface
type MockAuthClient struct {
	mock.Mock
}

type UserInfo struct {
	ID          string
	UserID      string
	Email       string
	Roles       []string
	Permissions []string
}

func (m *MockAuthClient) ValidateToken(token string) (*UserInfo, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserInfo), args.Error(1)
}

func (m *MockAuthClient) CheckPermission(userID, permission string) bool {
	args := m.Called(userID, permission)
	return args.Bool(0)
}
