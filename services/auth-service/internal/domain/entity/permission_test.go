package entity_test

import (
	"testing"

	"github.com/erp-cosmetics/auth-service/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestPermission_Matches(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		checkCode      string
		expectedResult bool
	}{
		{
			name:           "exact match",
			code:           "user:user:read",
			checkCode:      "user:user:read",
			expectedResult: true,
		},
		{
			name:           "no match",
			code:           "user:user:read",
			checkCode:      "user:user:write",
			expectedResult: false,
		},
		{
			name:           "resource wildcard match",
			code:           "wms:*:read",
			checkCode:      "wms:stock:read",
			expectedResult: true,
		},
		{
			name:           "resource wildcard no match action",
			code:           "wms:*:read",
			checkCode:      "wms:stock:write",
			expectedResult: false,
		},
		{
			name:           "super wildcard match",
			code:           "*:*:*",
			checkCode:      "any:any:any",
			expectedResult: true,
		},
		{
			name:           "service wildcard match",
			code:           "*:user:read",
			checkCode:      "auth:user:read",
			expectedResult: true,
		},
		{
			name:           "action wildcard match",
			code:           "user:user:*",
			checkCode:      "user:user:delete",
			expectedResult: true,
		},
		{
			name:           "invalid format permission code",
			code:           "user:user:read",
			checkCode:      "invalid:code",
			expectedResult: false,
		},
		{
			name:           "invalid format this code",
			code:           "invalid:code",
			checkCode:      "user:user:read",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entity.Permission{Code: tt.code}
			assert.Equal(t, tt.expectedResult, p.Matches(tt.checkCode))
		})
	}
}

func TestPermission_IsWildcard(t *testing.T) {
	assert.True(t, (&entity.Permission{Code: "wms:*:read"}).IsWildcard())
	assert.True(t, (&entity.Permission{Code: "*:*:*"}).IsWildcard())
	assert.False(t, (&entity.Permission{Code: "user:user:read"}).IsWildcard())
}
