package entity_test

import (
	"testing"
	"time"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_PasswordManagement(t *testing.T) {
	user := &entity.User{
		ID:    uuid.New(),
		Email: "test@company.vn",
	}

	password := "Secret@123"
	err := user.SetPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, password, user.PasswordHash)

	assert.True(t, user.VerifyPassword(password))
	assert.False(t, user.VerifyPassword("wrong-password"))
}

func TestUser_AccountLockout(t *testing.T) {
	user := &entity.User{
		ID:                  uuid.New(),
		FailedLoginAttempts: 0,
	}

	assert.False(t, user.IsLocked())

	// Increment 4 times
	for i := 0; i < 4; i++ {
		user.IncrementFailedAttempts()
		assert.False(t, user.IsLocked())
	}

	// 5th attempt should lock
	user.IncrementFailedAttempts()
	assert.True(t, user.IsLocked())
	assert.Equal(t, 5, user.FailedLoginAttempts)
	assert.NotNil(t, user.LockedUntil)

	// Reset should unlock
	user.ResetFailedAttempts()
	assert.False(t, user.IsLocked())
	assert.Equal(t, 0, user.FailedLoginAttempts)
	assert.Nil(t, user.LockedUntil)
}

func TestUser_IsLocked_Expired(t *testing.T) {
	user := &entity.User{}
	
	// Past lock time
	past := time.Now().Add(-1 * time.Minute)
	user.LockedUntil = &past
	assert.False(t, user.IsLocked())

	// Future lock time
	future := time.Now().Add(1 * time.Minute)
	user.LockedUntil = &future
	assert.True(t, user.IsLocked())
}

func TestUser_UpdateLastLogin(t *testing.T) {
	user := &entity.User{}
	assert.Nil(t, user.LastLoginAt)

	user.UpdateLastLogin()
	assert.NotNil(t, user.LastLoginAt)
	assert.WithinDuration(t, time.Now(), *user.LastLoginAt, 1*time.Second)
}
