package export

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter exports data to Excel format
type ExcelExporter struct{}

// NewExcelExporter creates a new Excel exporter
func NewExcelExporter() *ExcelExporter {
	return &ExcelExporter{}
}

// Export exports data to Excel (XLSX)
func (e *ExcelExporter) Export(sheetName string, headers []string, data [][]interface{}) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	// Create sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet: %w", err)
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
		// Auto-fit column width (approximate)
		colName, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, colName, colName, float64(len(header)+5))
	}

	// Write data
	for rowIdx, row := range data {
		for colIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, val)
		}
	}

	// Auto-filter
	if len(headers) > 0 && len(data) > 0 {
		lastCell, _ := excelize.CoordinatesToCellName(len(headers), len(data)+1)
		f.AutoFilter(sheetName, "A1:"+lastCell, nil)
	}

	// Write to buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write excel: %w", err)
	}

	return buf.Bytes(), nil
}

// ExportMultiSheet exports multiple sheets to Excel
func (e *ExcelExporter) ExportMultiSheet(sheets map[string]SheetData) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	first := true
	for sheetName, sheetData := range sheets {
		if first {
			f.SetSheetName("Sheet1", sheetName)
			first = false
		} else {
			f.NewSheet(sheetName)
		}

		// Write headers
		for i, header := range sheetData.Headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
		}

		// Write data
		for rowIdx, row := range sheetData.Data {
			for colIdx, val := range row {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
				f.SetCellValue(sheetName, cell, val)
			}
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write excel: %w", err)
	}

	return buf.Bytes(), nil
}

// SheetData represents data for a single sheet
type SheetData struct {
	Headers []string
	Data    [][]interface{}
}
