package report

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/erp-cosmetics/reporting-service/internal/domain/entity"
	"github.com/erp-cosmetics/reporting-service/internal/domain/repository"
	"github.com/erp-cosmetics/reporting-service/internal/infrastructure/export"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// UseCase defines report use case interface
type UseCase interface {
	GetDefinition(ctx context.Context, id uuid.UUID) (*entity.ReportDefinition, error)
	GetDefinitionByCode(ctx context.Context, code string) (*entity.ReportDefinition, error)
	ListDefinitions(ctx context.Context, page, pageSize int) (*ListOutput, error)
	ListByType(ctx context.Context, reportType string) ([]*entity.ReportDefinition, error)
	Execute(ctx context.Context, input *ExecuteInput) (*entity.ReportExecution, error)
	GetExecution(ctx context.Context, id uuid.UUID) (*entity.ReportExecution, error)
	ListExecutions(ctx context.Context, reportID uuid.UUID, page, pageSize int) (*ExecutionListOutput, error)
	Export(ctx context.Context, executionID uuid.UUID, format string) ([]byte, string, error)
}

// ExecuteInput for executing a report
type ExecuteInput struct {
	ReportID   uuid.UUID              `json:"report_id"`
	ReportCode string                 `json:"report_code"`
	Parameters map[string]interface{} `json:"parameters"`
	Format     string                 `json:"format"` // CSV, XLSX, JSON
	CreatedBy  *uuid.UUID
}

// ListOutput for listing reports
type ListOutput struct {
	Reports  []*entity.ReportDefinition `json:"reports"`
	Total    int64                      `json:"total"`
	Page     int                        `json:"page"`
	PageSize int                        `json:"page_size"`
}

// ExecutionListOutput for listing executions
type ExecutionListOutput struct {
	Executions []*entity.ReportExecution `json:"executions"`
	Total      int64                     `json:"total"`
	Page       int                       `json:"page"`
	PageSize   int                       `json:"page_size"`
}

type useCase struct {
	reportRepo    repository.ReportDefinitionRepository
	executionRepo repository.ReportExecutionRepository
	csvExporter   *export.CSVExporter
	excelExporter *export.ExcelExporter
	logger        *zap.Logger
}

// NewUseCase creates new report use case
func NewUseCase(
	reportRepo repository.ReportDefinitionRepository,
	executionRepo repository.ReportExecutionRepository,
	logger *zap.Logger,
) UseCase {
	return &useCase{
		reportRepo:    reportRepo,
		executionRepo: executionRepo,
		csvExporter:   export.NewCSVExporter(),
		excelExporter: export.NewExcelExporter(),
		logger:        logger,
	}
}

func (uc *useCase) GetDefinition(ctx context.Context, id uuid.UUID) (*entity.ReportDefinition, error) {
	return uc.reportRepo.GetByID(ctx, id)
}

func (uc *useCase) GetDefinitionByCode(ctx context.Context, code string) (*entity.ReportDefinition, error) {
	return uc.reportRepo.GetByCode(ctx, code)
}

func (uc *useCase) ListDefinitions(ctx context.Context, page, pageSize int) (*ListOutput, error) {
	offset := (page - 1) * pageSize
	reports, total, err := uc.reportRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ListOutput{
		Reports:  reports,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (uc *useCase) ListByType(ctx context.Context, reportType string) ([]*entity.ReportDefinition, error) {
	return uc.reportRepo.ListByType(ctx, reportType)
}

func (uc *useCase) Execute(ctx context.Context, input *ExecuteInput) (*entity.ReportExecution, error) {
	// Get report definition
	var report *entity.ReportDefinition
	var err error

	if input.ReportID != uuid.Nil {
		report, err = uc.reportRepo.GetByID(ctx, input.ReportID)
	} else if input.ReportCode != "" {
		report, err = uc.reportRepo.GetByCode(ctx, input.ReportCode)
	} else {
		return nil, fmt.Errorf("report_id or report_code required")
	}

	if err != nil {
		return nil, err
	}

	// Create execution record
	paramsJSON, _ := json.Marshal(input.Parameters)
	execution := &entity.ReportExecution{
		ReportID:   report.ID,
		Parameters: paramsJSON,
		Status:     entity.ExecutionStatusPending,
		FileFormat: input.Format,
		CreatedBy:  input.CreatedBy,
	}

	if err := uc.executionRepo.Create(ctx, execution); err != nil {
		return nil, err
	}

	// Execute report synchronously for now
	// In production, this would be async via message queue
	go uc.executeReport(context.Background(), execution, report)

	return execution, nil
}

func (uc *useCase) executeReport(ctx context.Context, execution *entity.ReportExecution, report *entity.ReportDefinition) {
	execution.Start()
	uc.executionRepo.Update(ctx, execution)

	// Simulate report execution
	// In production, this would:
	// 1. Parse query template with parameters
	// 2. Execute query against appropriate database
	// 3. Format and export results

	// Mock data for demo
	mockData := [][]interface{}{
		{"MAT001", "Raw Material A", "Raw Materials", "Main Warehouse", 1000.0, "KG", 500.0, "OK"},
		{"MAT002", "Raw Material B", "Raw Materials", "Main Warehouse", 250.0, "KG", 300.0, "LOW"},
		{"MAT003", "Packaging Box", "Packaging", "Packaging WH", 5000.0, "PCS", 2000.0, "OK"},
	}

	execution.Complete(len(mockData))

	// Generate preview
	preview, _ := json.Marshal(mockData[:min(10, len(mockData))])
	execution.ResultPreview = preview

	uc.executionRepo.Update(ctx, execution)

	uc.logger.Info("Report executed",
		zap.String("execution_id", execution.ID.String()),
		zap.String("report_code", report.Code),
		zap.Int("rows", execution.RowCount),
	)
}

func (uc *useCase) GetExecution(ctx context.Context, id uuid.UUID) (*entity.ReportExecution, error) {
	return uc.executionRepo.GetByID(ctx, id)
}

func (uc *useCase) ListExecutions(ctx context.Context, reportID uuid.UUID, page, pageSize int) (*ExecutionListOutput, error) {
	offset := (page - 1) * pageSize
	executions, total, err := uc.executionRepo.GetByReportID(ctx, reportID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	return &ExecutionListOutput{
		Executions: executions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

func (uc *useCase) Export(ctx context.Context, executionID uuid.UUID, format string) ([]byte, string, error) {
	execution, err := uc.executionRepo.GetByID(ctx, executionID)
	if err != nil {
		return nil, "", err
	}

	if !execution.IsComplete() {
		return nil, "", fmt.Errorf("execution not complete")
	}

	// Parse preview data for export
	var data [][]interface{}
	if err := json.Unmarshal(execution.ResultPreview, &data); err != nil {
		return nil, "", err
	}

	// Get column headers from report definition
	headers := []string{"Col1", "Col2", "Col3", "Col4", "Col5", "Col6", "Col7", "Col8"}

	var exportData []byte
	var filename string

	switch format {
	case entity.FormatCSV:
		exportData, err = uc.csvExporter.Export(headers, data)
		filename = "report.csv"
	case entity.FormatXLSX:
		exportData, err = uc.excelExporter.Export("Report", headers, data)
		filename = "report.xlsx"
	default:
		exportData, err = json.Marshal(data)
		filename = "report.json"
	}

	if err != nil {
		return nil, "", err
	}

	return exportData, filename, nil
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
