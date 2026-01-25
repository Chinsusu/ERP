package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ExecutionStatus constants
const (
	ExecutionStatusPending   = "PENDING"
	ExecutionStatusRunning   = "RUNNING"
	ExecutionStatusCompleted = "COMPLETED"
	ExecutionStatusFailed    = "FAILED"
	ExecutionStatusCancelled = "CANCELLED"
)

// ExportFormat constants
const (
	FormatCSV  = "CSV"
	FormatXLSX = "XLSX"
	FormatPDF  = "PDF"
	FormatJSON = "JSON"
)

// ReportExecution represents a report execution instance
type ReportExecution struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ReportID      uuid.UUID      `gorm:"type:uuid;not null" json:"report_id"`
	Parameters    datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"parameters"`
	Status        string         `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	Progress      int            `gorm:"default:0" json:"progress"`
	StartedAt     *time.Time     `gorm:"type:timestamp" json:"started_at,omitempty"`
	CompletedAt   *time.Time     `gorm:"type:timestamp" json:"completed_at,omitempty"`
	RowCount      int            `gorm:"default:0" json:"row_count"`
	ResultPreview datatypes.JSON `gorm:"type:jsonb" json:"result_preview,omitempty"`
	FilePath      string         `gorm:"type:text" json:"file_path,omitempty"`
	FileFormat    string         `gorm:"type:varchar(20)" json:"file_format,omitempty"`
	FileSize      int64          `gorm:"type:bigint" json:"file_size,omitempty"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message,omitempty"`
	CreatedBy     *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Report *ReportDefinition `gorm:"foreignKey:ReportID" json:"report,omitempty"`
}

// TableName specifies table name
func (ReportExecution) TableName() string {
	return "report_executions"
}

// BeforeCreate sets defaults
func (e *ReportExecution) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// Start marks execution as started
func (e *ReportExecution) Start() {
	now := time.Now()
	e.Status = ExecutionStatusRunning
	e.StartedAt = &now
	e.Progress = 0
}

// Complete marks execution as completed
func (e *ReportExecution) Complete(rowCount int) {
	now := time.Now()
	e.Status = ExecutionStatusCompleted
	e.CompletedAt = &now
	e.Progress = 100
	e.RowCount = rowCount
}

// Fail marks execution as failed
func (e *ReportExecution) Fail(err string) {
	now := time.Now()
	e.Status = ExecutionStatusFailed
	e.CompletedAt = &now
	e.ErrorMessage = err
}

// UpdateProgress updates execution progress
func (e *ReportExecution) UpdateProgress(progress int) {
	if progress > 100 {
		progress = 100
	}
	e.Progress = progress
}

// IsComplete checks if execution is complete
func (e *ReportExecution) IsComplete() bool {
	return e.Status == ExecutionStatusCompleted
}

// IsFailed checks if execution failed
func (e *ReportExecution) IsFailed() bool {
	return e.Status == ExecutionStatusFailed
}

// HasExportFile checks if export file exists
func (e *ReportExecution) HasExportFile() bool {
	return e.FilePath != "" && e.FileSize > 0
}
