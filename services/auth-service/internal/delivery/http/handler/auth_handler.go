package handler

import (
	"github.com/erp-cosmetics/auth-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth"
	"github.com/erp-cosmetics/shared/pkg/middleware"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	loginUC   *auth.LoginUseCase
	logoutUC  *auth.LogoutUseCase
	refreshUC *auth.RefreshTokenUseCase
}

func NewAuthHandler(
	loginUC *auth.LoginUseCase,
	logoutUC *auth.LogoutUseCase,
	refreshUC *auth.RefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		loginUC:   loginUC,
		logoutUC:  logoutUC,
		refreshUC: refreshUC,
	}
}

// Login handles user login
// @Summary User login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.loginUC.Execute(c.Request.Context(), &auth.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
		User: dto.UserInfo{
			ID:          result.User.ID.String(),
			UserID:      result.User.UserID.String(),
			Email:       result.User.Email,
			Roles:       result.User.Roles,
			Permissions: result.User.Permissions,
		},
	})
}

// Logout handles user logout
// @Summary User logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	c.ShouldBindJSON(&req)

	// Get access token from header
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	err := h.logoutUC.Execute(c.Request.Context(), &auth.LogoutRequest{
		RefreshToken: req.RefreshToken,
		AccessToken:  accessToken,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Logged out successfully"})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.refreshUC.Execute(c.Request.Context(), &auth.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    "Bearer",
	})
}

// GetMe returns current user info
// @Summary Get current user
// @Tags Auth
// @Produce json
// @Success 200 {object} dto.UserInfo
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)
	email := middleware.GetEmail(c)
	roleIDs := middleware.GetRoleIDs(c)

	response.Success(c, dto.UserInfo{
		UserID:      userID,
		Email:       email,
		Roles:       roleIDs,
		Permissions: []string{}, // TODO: Fetch from cache/DB
	})
}
