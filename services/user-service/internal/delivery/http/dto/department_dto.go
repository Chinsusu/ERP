package dto

// CreateDepartmentRequest represents department creation request
type CreateDepartmentRequest struct {
	Code     string  `json:"code" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parent_id"`
}

// UpdateDepartmentRequest represents department update request
type UpdateDepartmentRequest struct {
	Name      string  `json:"name"`
	ManagerID *string `json:"manager_id"`
}

// DepartmentResponse represents department response
type DepartmentResponse struct {
	ID       string                   `json:"id"`
	Code     string                   `json:"code"`
	Name     string                   `json:"name"`
	ParentID string                   `json:"parent_id,omitempty"`
	Level    int                      `json:"level"`
	Path     string                   `json:"path"`
	Status   string                   `json:"status"`
	Children []DepartmentResponse     `json:"children,omitempty"`
}
