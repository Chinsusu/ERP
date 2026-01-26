package handler

import (
	"net/http"
	"strconv"

	"github.com/erp-cosmetics/reporting-service/internal/usecase/report"
	"github.com/erp-cosmetics/shared/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/erp-cosmetics/shared/pkg/errors")

// ReportHandler handles report requests
type ReportHandler struct {
	reportUC report.UseCase
}

// NewReportHandler creates new handler
func NewReportHandler(uc report.UseCase) *ReportHandler {
	return &ReportHandler{reportUC: uc}
}

// ListDefinitions handles GET /api/v1/reports
func (h *ReportHandler) ListDefinitions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	reportType := c.Query("type")

	if reportType != "" {
		reports, err := h.reportUC.ListByType(c.Request.Context(), reportType)
		if err != nil {
			response.Error(c, errors.New("ERROR", "Failed to list reports", http.StatusInternalServerError))
			return
		}
		response.Success(c, reports)
		return
	}

	output, err := h.reportUC.ListDefinitions(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to list reports", http.StatusInternalServerError))
		return
	}

	response.SuccessWithMeta(c, output.Reports, response.NewMeta(page, pageSize, output.Total))
}

// GetDefinition handles GET /api/v1/reports/:id
func (h *ReportHandler) GetDefinition(c *gin.Context) {
	idStr := c.Param("id")

	// Try to parse as UUID first
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Try as code
		report, err := h.reportUC.GetDefinitionByCode(c.Request.Context(), idStr)
		if err != nil {
			response.Error(c, errors.New("ERROR", "Report not found", http.StatusNotFound))
			return
		}
		response.Success(c, report)
		return
	}

	reportDef, err := h.reportUC.GetDefinition(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Report not found", http.StatusNotFound))
		return
	}

	response.Success(c, reportDef)
}

// Execute handles POST /api/v1/reports/:id/execute
func (h *ReportHandler) Execute(c *gin.Context) {
	idStr := c.Param("id")

	var input report.ExecuteInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// Allow empty body for reports without parameters
		input = report.ExecuteInput{}
	}

	// Try to parse as UUID first
	id, err := uuid.Parse(idStr)
	if err != nil {
		// Use as code
		input.ReportCode = idStr
	} else {
		input.ReportID = id
	}

	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			input.CreatedBy = &uid
		}
	}

	execution, err := h.reportUC.Execute(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to execute report", http.StatusInternalServerError))
		return
	}

	response.Success(c, execution)
}

// ListExecutions handles GET /api/v1/reports/:id/executions
func (h *ReportHandler) ListExecutions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid report ID", http.StatusBadRequest))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	output, err := h.reportUC.ListExecutions(c.Request.Context(), id, page, pageSize)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to list executions", http.StatusInternalServerError))
		return
	}

	response.SuccessWithMeta(c, output.Executions, response.NewMeta(page, pageSize, output.Total))
}

// GetExecution handles GET /api/v1/reports/executions/:id
func (h *ReportHandler) GetExecution(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid execution ID", http.StatusBadRequest))
		return
	}

	execution, err := h.reportUC.GetExecution(c.Request.Context(), id)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Execution not found", http.StatusNotFound))
		return
	}

	response.Success(c, execution)
}

// Download handles GET /api/v1/reports/executions/:id/download
func (h *ReportHandler) Download(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, errors.New("ERROR", "Invalid execution ID", http.StatusBadRequest))
		return
	}

	format := c.DefaultQuery("format", "xlsx")

	data, filename, err := h.reportUC.Export(c.Request.Context(), id, format)
	if err != nil {
		response.Error(c, errors.New("ERROR", "Failed to export report", http.StatusInternalServerError))
		return
	}

	// Set content type based on format
	contentType := "application/octet-stream"
	switch format {
	case "csv":
		contentType = "text/csv"
	case "xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "json":
		contentType = "application/json"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, data)
}
