package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrExpiredToken is returned when token is expired
	ErrExpiredToken = errors.New("token expired")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	// AccessToken is used for API authentication
	AccessToken TokenType = "access"
	// RefreshToken is used to obtain new access tokens
	RefreshToken TokenType = "refresh"
)

// Claims represents JWT custom claims
type Claims struct {
	UserID      string   `json:"sub"`
	Email       string   `json:"email"`
	RoleIDs     []string `json:"role_ids,omitempty"`
	RoleNames   []string `json:"role_names,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Type        TokenType `json:"type"`
	jwt.RegisteredClaims
}

// Manager handles JWT token operations
type Manager struct {
	secret                []byte
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

// NewManager creates a new JWT manager
func NewManager(secret string, accessTokenExpire, refreshTokenExpire time.Duration) *Manager {
	return &Manager{
		secret:                []byte(secret),
		accessTokenDuration:   accessTokenExpire,
		refreshTokenDuration:  refreshTokenExpire,
	}
}

// GenerateAccessToken creates a new access token
func (m *Manager) GenerateAccessToken(userID, email string, roleIDs []string, roleNames []string, permissions []string) (string, error) {
	return m.generateToken(userID, email, roleIDs, roleNames, permissions, AccessToken, m.accessTokenDuration)
}

// GenerateRefreshToken creates a new refresh token
func (m *Manager) GenerateRefreshToken(userID, email string) (string, error) {
	return m.generateToken(userID, email, nil, nil, nil, RefreshToken, m.refreshTokenDuration)
}

// generateToken creates a JWT token
func (m *Manager) generateToken(userID, email string, roleIDs []string, roleNames []string, permissions []string, tokenType TokenType, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:      userID,
		Email:       email,
		RoleIDs:     roleIDs,
		RoleNames:   roleNames,
		Permissions: permissions,
		Type:        tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        generateJTI(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// VerifyToken verifies and parses a JWT token
func (m *Manager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractUserID extracts user ID from token
func (m *Manager) ExtractUserID(tokenString string) (string, error) {
	claims, err := m.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}

// generateJTI generates a unique JWT ID
func generateJTI() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// ParseDuration parses a duration string (e.g., "15m", "24h", "7d")
func ParseDuration(s string) (time.Duration, error) {
	// Handle days
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days := s[:len(s)-1]
		d, err := time.ParseDuration(days + "h")
		if err != nil {
			return 0, err
		}
		return d * 24, nil
	}
	return time.ParseDuration(s)
}
