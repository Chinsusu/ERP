package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response helpers - extending shared package for this service

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *errorInfo  `json:"error,omitempty"`
	Meta    *meta       `json:"meta,omitempty"`
}

type errorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	TotalItems int64 `json:"total_items,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    data,
	})
}

func created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, apiResponse{
		Success: true,
		Data:    data,
	})
}

func successWithMeta(c *gin.Context, data interface{}, m *meta) {
	c.JSON(http.StatusOK, apiResponse{
		Success: true,
		Data:    data,
		Meta:    m,
	})
}

func badRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, apiResponse{
		Success: false,
		Error: &errorInfo{
			Code:    "BAD_REQUEST",
			Message: message,
		},
	})
}

func notFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, apiResponse{
		Success: false,
		Error: &errorInfo{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

func internalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, apiResponse{
		Success: false,
		Error: &errorInfo{
			Code:    "INTERNAL_ERROR",
			Message: message,
		},
	})
}

func newMeta(page, pageSize int, totalItems int64) *meta {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}
	return &meta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
