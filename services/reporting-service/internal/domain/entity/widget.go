package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// WidgetType constants
const (
	WidgetTypeKPI       = "KPI"
	WidgetTypeLineChart = "LINE_CHART"
	WidgetTypeBarChart  = "BAR_CHART"
	WidgetTypePieChart  = "PIE_CHART"
	WidgetTypeTable     = "TABLE"
	WidgetTypeGauge     = "GAUGE"
	WidgetTypeMap       = "MAP"
)

// Widget represents a dashboard widget
type Widget struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DashboardID     uuid.UUID      `gorm:"type:uuid;not null" json:"dashboard_id"`
	WidgetType      string         `gorm:"type:varchar(50);not null" json:"widget_type"`
	Title           string         `gorm:"type:varchar(200);not null" json:"title"`
	Subtitle        string         `gorm:"type:varchar(300)" json:"subtitle,omitempty"`
	Icon            string         `gorm:"type:varchar(50)" json:"icon,omitempty"`
	DataSource      string         `gorm:"type:varchar(100)" json:"data_source,omitempty"`
	QueryParams     datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"query_params,omitempty"`
	RefreshInterval int            `gorm:"default:300" json:"refresh_interval"` // Seconds
	Config          datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"config,omitempty"`
	PositionX       int            `gorm:"default:0" json:"position_x"`
	PositionY       int            `gorm:"default:0" json:"position_y"`
	Width           int            `gorm:"default:4" json:"width"`
	Height          int            `gorm:"default:2" json:"height"`
	IsVisible       bool           `gorm:"default:true" json:"is_visible"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies table name
func (Widget) TableName() string {
	return "widgets"
}

// BeforeCreate sets defaults
func (w *Widget) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// IsChartWidget checks if widget is a chart
func (w *Widget) IsChartWidget() bool {
	return w.WidgetType == WidgetTypeLineChart ||
		w.WidgetType == WidgetTypeBarChart ||
		w.WidgetType == WidgetTypePieChart
}

// IsKPIWidget checks if widget is a KPI
func (w *Widget) IsKPIWidget() bool {
	return w.WidgetType == WidgetTypeKPI
}

// IsTableWidget checks if widget is a table
func (w *Widget) IsTableWidget() bool {
	return w.WidgetType == WidgetTypeTable
}

// GetGridPosition returns grid position string
func (w *Widget) GetGridPosition() string {
	return ""
}

// WidgetConfig represents widget-specific configuration
type WidgetConfig struct {
	// Chart config
	XField     string   `json:"xField,omitempty"`
	YField     string   `json:"yField,omitempty"`
	Color      string   `json:"color,omitempty"`
	Colors     []string `json:"colors,omitempty"`
	ShowLegend bool     `json:"showLegend,omitempty"`

	// KPI config
	Icon      string `json:"icon,omitempty"`
	Format    string `json:"format,omitempty"` // number, currency, percentage
	Field     string `json:"field,omitempty"`
	Threshold int    `json:"threshold,omitempty"`

	// Table config
	Limit   int      `json:"limit,omitempty"`
	Columns []string `json:"columns,omitempty"`
}
