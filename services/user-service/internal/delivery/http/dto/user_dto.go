package dto

import "time"

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Email        string     `json:"email" binding:"required,email"`
	Password     string     `json:"password" binding:"required,min=8"`
	FirstName    string     `json:"first_name" binding:"required"`
	LastName     string     `json:"last_name" binding:"required"`
	Phone        string     `json:"phone"`
	DepartmentID *string    `json:"department_id"`
	ManagerID    *string    `json:"manager_id"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	Address      string     `json:"address"`
	JoinDate     *time.Time `json:"join_date"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Phone        string  `json:"phone"`
	DepartmentID *string `json:"department_id"`
	ManagerID    *string `json:"manager_id"`
}

// UserResponse represents user response
type UserResponse struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	EmployeeCode string `json:"employee_code"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone,omitempty"`
	AvatarURL    string `json:"avatar_url,omitempty"`
	DepartmentID string `json:"department_id,omitempty"`
	ManagerID    string `json:"manager_id,omitempty"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}
