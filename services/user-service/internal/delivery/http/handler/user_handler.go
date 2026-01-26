package handler

import (
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
