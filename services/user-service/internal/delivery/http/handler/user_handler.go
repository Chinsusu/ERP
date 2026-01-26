package handler

import (
	"strings"
	"time"

	"github.com/erp-cosmetics/shared/pkg/errors"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/erp-cosmetics/user-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/user-service/internal/usecase/user"
	"github.com/gin-gonic/gin"
)


type UserHandler struct {
	createUC *user.CreateUserUseCase
	getUC    *user.GetUserUseCase
	listUC   *user.ListUsersUseCase
}

func NewUserHandler(
	createUC *user.CreateUserUseCase,
	getUC *user.GetUserUseCase,
	listUC *user.ListUsersUseCase,
) *UserHandler {
	return &UserHandler{
		createUC: createUC,
		getUC:    getUC,
		listUC:   listUC,
	}
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	ucReq := &user.CreateUserRequest{
		Email:        req.Email,
		Password:     req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		DepartmentID: req.DepartmentID,
		ManagerID:    req.ManagerID,
		DateOfBirth:  req.DateOfBirth,
		Address:      req.Address,
		JoinDate:     req.JoinDate,
	}

	user, err := h.createUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// GetUser gets user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.getUC.Execute(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// GetMe returns the current logged in user
func (h *UserHandler) GetMe(c *gin.Context) {
	// Extract user ID from header (set by API Gateway)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		// Fallback for direct access if needed
		val, _ := c.Get("user_id")
		userID, _ = val.(string)
	}

	if userID == "" {
		response.Error(c, errors.Unauthorized("User ID not found in request"))
		return
	}

	user, err := h.getUC.Execute(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// Map to response DTO
	res := dto.UserResponse{
		ID:           user.ID.String(),
		Email:        user.Email,
		EmployeeCode: user.EmployeeCode,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		AvatarURL:    user.AvatarURL,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
	}

	if user.DepartmentID != nil {
		res.DepartmentID = user.DepartmentID.String()
	}
	if user.ManagerID != nil {
		res.ManagerID = user.ManagerID.String()
	}

	// Extract roles from header (set by API Gateway from JWT)
	rolesHeader := c.GetHeader("X-User-Roles")
	if rolesHeader != "" {
		roleNames := strings.Split(rolesHeader, ",")
		res.Roles = make([]dto.RoleResponse, len(roleNames))
		for i, name := range roleNames {
			// Convert to snake_case for frontend (e.g., "Super Admin" -> "super_admin")
			res.Roles[i] = dto.RoleResponse{Name: toSnakeCase(name)}
		}
	} else {
		res.Roles = []dto.RoleResponse{}
	}

	// Extract permissions from header (set by API Gateway from JWT)
	permsHeader := c.GetHeader("X-User-Permissions")
	if permsHeader != "" {
		res.Permissions = strings.Split(permsHeader, ",")
	} else {
		res.Permissions = []string{}
	}

	response.Success(c, res)
}

func toSnakeCase(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}

// ListUsers lists users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req user.ListUsersRequest
	
	// Parse query parameters
	req.Status = c.Query("status")
	req.Search = c.Query("search")
	req.Page = 1
	req.PageSize = 20

	if page := c.Query("page"); page != "" {
		// Parse page number
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		// Parse page size
	}

	result, err := h.listUC.Execute(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}
