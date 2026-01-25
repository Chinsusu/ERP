package handler

import (
	"net/http"
	"strconv"

	"github.com/erp-cosmetics/file-service/internal/usecase/file"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileHandler handles file requests
type FileHandler struct {
	fileUC file.UseCase
}

// NewFileHandler creates new handler
func NewFileHandler(uc file.UseCase) *FileHandler {
	return &FileHandler{fileUC: uc}
}

// Upload handles POST /api/v1/files/upload
func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "No file provided", err.Error())
		return
	}

	category := c.PostForm("category")
	if category == "" {
		category = "ATTACHMENT"
	}

	entityType := c.PostForm("entity_type")
	entityIDStr := c.PostForm("entity_id")
	isPublic := c.PostForm("is_public") == "true"

	var entityID *uuid.UUID
	if entityIDStr != "" {
		id, err := uuid.Parse(entityIDStr)
		if err == nil {
			entityID = &id
		}
	}

	// Open file
	src, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to read file", err.Error())
		return
	}
	defer src.Close()

	// Get user ID from context
	var createdBy *uuid.UUID
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			createdBy = &uid
		}
	}

	input := &file.UploadInput{
		FileName:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
		Size:        fileHeader.Size,
		Reader:      src,
		Category:    category,
		EntityType:  entityType,
		EntityID:    entityID,
		IsPublic:    isPublic,
		CreatedBy:   createdBy,
	}

	result, err := h.fileUC.Upload(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload file", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "File uploaded", result)
}

// UploadMultiple handles POST /api/v1/files/upload/multiple
func (h *FileHandler) UploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid form", err.Error())
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, http.StatusBadRequest, "No files provided", "")
		return
	}

	category := c.PostForm("category")
	if category == "" {
		category = "ATTACHMENT"
	}

	entityType := c.PostForm("entity_type")
	entityIDStr := c.PostForm("entity_id")

	var entityID *uuid.UUID
	if entityIDStr != "" {
		id, err := uuid.Parse(entityIDStr)
		if err == nil {
			entityID = &id
		}
	}

	var createdBy *uuid.UUID
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			createdBy = &uid
		}
	}

	var results []*struct {
		Filename string      `json:"filename"`
		Success  bool        `json:"success"`
		File     interface{} `json:"file,omitempty"`
		Error    string      `json:"error,omitempty"`
	}

	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			results = append(results, &struct {
				Filename string      `json:"filename"`
				Success  bool        `json:"success"`
				File     interface{} `json:"file,omitempty"`
				Error    string      `json:"error,omitempty"`
			}{
				Filename: fileHeader.Filename,
				Success:  false,
				Error:    err.Error(),
			})
			continue
		}

		input := &file.UploadInput{
			FileName:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Size:        fileHeader.Size,
			Reader:      src,
			Category:    category,
			EntityType:  entityType,
			EntityID:    entityID,
			CreatedBy:   createdBy,
		}

		result, err := h.fileUC.Upload(c.Request.Context(), input)
		src.Close()

		if err != nil {
			results = append(results, &struct {
				Filename string      `json:"filename"`
				Success  bool        `json:"success"`
				File     interface{} `json:"file,omitempty"`
				Error    string      `json:"error,omitempty"`
			}{
				Filename: fileHeader.Filename,
				Success:  false,
				Error:    err.Error(),
			})
		} else {
			results = append(results, &struct {
				Filename string      `json:"filename"`
				Success  bool        `json:"success"`
				File     interface{} `json:"file,omitempty"`
				Error    string      `json:"error,omitempty"`
			}{
				Filename: fileHeader.Filename,
				Success:  true,
				File:     result,
			})
		}
	}

	response.Success(c, http.StatusOK, "Files processed", results)
}

// Get handles GET /api/v1/files/:id
func (h *FileHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID", err.Error())
		return
	}

	file, err := h.fileUC.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "File not found", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "File retrieved", file)
}

// Download handles GET /api/v1/files/:id/download
func (h *FileHandler) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID", err.Error())
		return
	}

	reader, file, err := h.fileUC.Download(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "File not found", err.Error())
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Length", strconv.FormatInt(file.FileSize, 10))
	c.DataFromReader(http.StatusOK, file.FileSize, file.ContentType, reader, nil)
}

// GetDownloadURL handles GET /api/v1/files/:id/url
func (h *FileHandler) GetDownloadURL(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID", err.Error())
		return
	}

	url, err := h.fileUC.GetDownloadURL(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate URL", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Download URL generated", gin.H{"url": url})
}

// Delete handles DELETE /api/v1/files/:id
func (h *FileHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID", err.Error())
		return
	}

	if err := h.fileUC.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete file", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "File deleted", nil)
}

// GetByEntity handles GET /api/v1/files/entity/:type/:id
func (h *FileHandler) GetByEntity(c *gin.Context) {
	entityType := c.Param("type")
	entityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid entity ID", err.Error())
		return
	}

	files, err := h.fileUC.GetByEntity(c.Request.Context(), entityType, entityID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get files", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Files retrieved", files)
}

// ListCategories handles GET /api/v1/files/categories
func (h *FileHandler) ListCategories(c *gin.Context) {
	categories, err := h.fileUC.ListCategories(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list categories", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Categories retrieved", categories)
}
