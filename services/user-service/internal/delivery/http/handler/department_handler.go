package handler

import (
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/erp-cosmetics/user-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/user-service/internal/usecase/department"
	"github.com/gin-gonic/gin"
)


type DepartmentHandler struct {
	createUC  *department.CreateDepartmentUseCase
	getTreeUC *department.GetDepartmentTreeUseCase
}

func NewDepartmentHandler(
	createUC *department.CreateDepartmentUseCase,
	getTreeUC *department.GetDepartmentTreeUseCase,
) *DepartmentHandler {
	return &DepartmentHandler{
		createUC:  createUC,
		getTreeUC: getTreeUC,
	}
}

// CreateDepartment creates a new department
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	ucReq := &department.CreateDepartmentRequest{
		Code:     req.Code,
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	dept, err := h.createUC.Execute(c.Request.Context(), ucReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dept)
}

// GetDepartmentTree gets department hierarchy as tree
func (h *DepartmentHandler) GetDepartmentTree(c *gin.Context) {
	tree, err := h.getTreeUC.Execute(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, tree)
}
