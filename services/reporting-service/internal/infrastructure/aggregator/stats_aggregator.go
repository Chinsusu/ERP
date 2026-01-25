package aggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// ServiceClient handles communication with other services
type ServiceClient struct {
	httpClient *http.Client
	baseURLs   map[string]string
	logger     *zap.Logger
}

// NewServiceClient creates a new service client
func NewServiceClient(baseURLs map[string]string, logger *zap.Logger) *ServiceClient {
	return &ServiceClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURLs: baseURLs,
		logger:   logger,
	}
}

// FetchData fetches data from a service
func (c *ServiceClient) FetchData(ctx context.Context, service, endpoint string, params map[string]string) (map[string]interface{}, error) {
	baseURL, ok := c.baseURLs[service]
	if !ok {
		return nil, fmt.Errorf("unknown service: %s", service)
	}

	url := baseURL + endpoint
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query params
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error("Failed to fetch data", 
			zap.String("service", service),
			zap.String("endpoint", endpoint),
			zap.Error(err),
		)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service returned status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// InventoryStats holds inventory statistics
type InventoryStats struct {
	TotalStockValue    float64 `json:"total_stock_value"`
	TotalItems         int     `json:"total_items"`
	LowStockItems      int     `json:"low_stock_items"`
	ExpiringIn30Days   int     `json:"expiring_30_days"`
	ExpiringIn7Days    int     `json:"expiring_7_days"`
	ExpiredItems       int     `json:"expired_items"`
	TotalWarehouses    int     `json:"total_warehouses"`
}

// SalesStats holds sales statistics
type SalesStats struct {
	TotalSales          float64 `json:"total_sales"`
	TotalOrders         int     `json:"total_orders"`
	PendingOrders       int     `json:"pending_orders"`
	TotalCustomers      int     `json:"total_customers"`
	AverageOrderValue   float64 `json:"average_order_value"`
	TopSellingProduct   string  `json:"top_selling_product"`
}

// ProductionStats holds production statistics
type ProductionStats struct {
	ActiveWorkOrders    int     `json:"active_work_orders"`
	CompletedToday      int     `json:"completed_today"`
	TotalOutput         float64 `json:"total_output"`
	QCPassRate          float64 `json:"qc_pass_rate"`
	PendingQC           int     `json:"pending_qc"`
}

// ProcurementStats holds procurement statistics
type ProcurementStats struct {
	PendingPRs          int     `json:"pending_prs"`
	OpenPOs             int     `json:"open_pos"`
	TotalPOValue        float64 `json:"total_po_value"`
	PendingDeliveries   int     `json:"pending_deliveries"`
	OverdueDeliveries   int     `json:"overdue_deliveries"`
}

// StatsAggregator aggregates statistics from various sources
type StatsAggregator struct {
	serviceClient *ServiceClient
	logger        *zap.Logger
}

// NewStatsAggregator creates a new stats aggregator
func NewStatsAggregator(serviceClient *ServiceClient, logger *zap.Logger) *StatsAggregator {
	return &StatsAggregator{
		serviceClient: serviceClient,
		logger:        logger,
	}
}

// GetInventoryStats returns inventory statistics
// In production, this would query WMS service or directly query database
func (a *StatsAggregator) GetInventoryStats(ctx context.Context) (*InventoryStats, error) {
	// Mock data for now - in production, query from WMS service or database
	return &InventoryStats{
		TotalStockValue:    1250000.00,
		TotalItems:         1500,
		LowStockItems:      23,
		ExpiringIn30Days:   45,
		ExpiringIn7Days:    12,
		ExpiredItems:       3,
		TotalWarehouses:    4,
	}, nil
}

// GetSalesStats returns sales statistics
func (a *StatsAggregator) GetSalesStats(ctx context.Context) (*SalesStats, error) {
	return &SalesStats{
		TotalSales:         850000.00,
		TotalOrders:        156,
		PendingOrders:      23,
		TotalCustomers:     89,
		AverageOrderValue:  5448.72,
		TopSellingProduct:  "Lipstick Matte Collection",
	}, nil
}

// GetProductionStats returns production statistics
func (a *StatsAggregator) GetProductionStats(ctx context.Context) (*ProductionStats, error) {
	return &ProductionStats{
		ActiveWorkOrders:   12,
		CompletedToday:     5,
		TotalOutput:        15000,
		QCPassRate:         98.5,
		PendingQC:          8,
	}, nil
}

// GetProcurementStats returns procurement statistics
func (a *StatsAggregator) GetProcurementStats(ctx context.Context) (*ProcurementStats, error) {
	return &ProcurementStats{
		PendingPRs:         15,
		OpenPOs:            28,
		TotalPOValue:       320000.00,
		PendingDeliveries:  12,
		OverdueDeliveries:  2,
	}, nil
}
