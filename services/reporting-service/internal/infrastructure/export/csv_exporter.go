package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
)

// CSVExporter exports data to CSV format
type CSVExporter struct{}

// NewCSVExporter creates a new CSV exporter
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{}
}

// Export exports data to CSV
func (e *CSVExporter) Export(headers []string, data [][]interface{}) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write headers
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write headers: %w", err)
	}

	// Write data rows
	for _, row := range data {
		record := make([]string, len(row))
		for i, val := range row {
			record[i] = formatValue(val)
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("csv writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// formatValue converts interface to string
func formatValue(val interface{}) string {
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%v", val)
}
