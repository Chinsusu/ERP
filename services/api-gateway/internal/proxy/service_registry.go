package proxy

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ServiceStatus represents service health status
type ServiceStatus int

const (
	StatusHealthy ServiceStatus = iota
	StatusUnhealthy
	StatusUnknown
)

func (s ServiceStatus) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

// ServiceInfo holds service information
type ServiceInfo struct {
	Name    string
	URL     string
	Status  ServiceStatus
	LastCheck time.Time
}

// ServiceRegistry manages service URLs and health
type ServiceRegistry struct {
	services map[string]*ServiceInfo
	mu       sync.RWMutex
	logger   *zap.Logger
	client   *http.Client
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(logger *zap.Logger) *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*ServiceInfo),
		logger:   logger,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Register adds a service to the registry
func (r *ServiceRegistry) Register(name, host string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	url := fmt.Sprintf("http://%s", host)
	r.services[name] = &ServiceInfo{
		Name:   name,
		URL:    url,
		Status: StatusUnknown,
	}
}

// GetServiceURL returns the URL for a service
func (r *ServiceRegistry) GetServiceURL(name string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	svc, exists := r.services[name]
	if !exists {
		return "", fmt.Errorf("service %s not found", name)
	}

	return svc.URL, nil
}

// MarkHealthy marks a service as healthy
func (r *ServiceRegistry) MarkHealthy(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if svc, exists := r.services[name]; exists {
		svc.Status = StatusHealthy
		svc.LastCheck = time.Now()
	}
}

// MarkUnhealthy marks a service as unhealthy
func (r *ServiceRegistry) MarkUnhealthy(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if svc, exists := r.services[name]; exists {
		svc.Status = StatusUnhealthy
		svc.LastCheck = time.Now()
	}
}

// GetStatus returns status of a service
func (r *ServiceRegistry) GetStatus(name string) ServiceStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if svc, exists := r.services[name]; exists {
		return svc.Status
	}
	return StatusUnknown
}

// GetAllStatuses returns status of all services
func (r *ServiceRegistry) GetAllStatuses() map[string]string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	statuses := make(map[string]string)
	for name, svc := range r.services {
		statuses[name] = svc.Status.String()
	}
	return statuses
}

// HealthCheck performs health check on a service
func (r *ServiceRegistry) HealthCheck(name string) error {
	r.mu.RLock()
	svc, exists := r.services[name]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not found", name)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", svc.URL+"/health", nil)
	if err != nil {
		r.MarkUnhealthy(name)
		return err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		r.MarkUnhealthy(name)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		r.MarkHealthy(name)
		return nil
	}

	r.MarkUnhealthy(name)
	return fmt.Errorf("health check failed with status %d", resp.StatusCode)
}

// StartHealthChecks starts periodic health checks
func (r *ServiceRegistry) StartHealthChecks(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			r.mu.RLock()
			services := make([]string, 0, len(r.services))
			for name := range r.services {
				services = append(services, name)
			}
			r.mu.RUnlock()

			for _, name := range services {
				if err := r.HealthCheck(name); err != nil {
					r.logger.Warn("Health check failed",
						zap.String("service", name),
						zap.Error(err),
					)
				}
			}
		}
	}()
}
