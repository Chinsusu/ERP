package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig holds circuit breaker configuration
type CircuitBreakerConfig struct {
	Enabled   bool
	Threshold int           // Number of failures before opening
	Timeout   time.Duration // How long to stay open
}

// CircuitBreakerManager manages circuit breakers for multiple services
type CircuitBreakerManager struct {
	breakers map[string]*gobreaker.CircuitBreaker
	mu       sync.RWMutex
	config   CircuitBreakerConfig
}

// NewCircuitBreakerManager creates a new circuit breaker manager
func NewCircuitBreakerManager(config CircuitBreakerConfig) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
		config:   config,
	}
}

// GetBreaker returns circuit breaker for a service
func (m *CircuitBreakerManager) GetBreaker(serviceName string) *gobreaker.CircuitBreaker {
	m.mu.RLock()
	cb, exists := m.breakers[serviceName]
	m.mu.RUnlock()

	if exists {
		return cb
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double check
	if cb, exists := m.breakers[serviceName]; exists {
		return cb
	}

	settings := gobreaker.Settings{
		Name:        serviceName,
		MaxRequests: 1, // Number of requests allowed in half-open state
		Interval:    60 * time.Second,
		Timeout:     m.config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= uint32(m.config.Threshold)
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// Log state changes - could add metrics here
		},
	}

	cb = gobreaker.NewCircuitBreaker(settings)
	m.breakers[serviceName] = cb
	return cb
}

// Middleware checks circuit breaker state before allowing request
func (m *CircuitBreakerManager) Middleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.config.Enabled {
			c.Next()
			return
		}

		cb := m.GetBreaker(serviceName)

		state := cb.State()
		if state == gobreaker.StateOpen {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error":       "Service temporarily unavailable",
				"code":        "SERVICE_UNAVAILABLE",
				"service":     serviceName,
				"retry_after": int(m.config.Timeout.Seconds()),
			})
			return
		}

		c.Next()

		// Note: The circuit breaker state is managed internally by gobreaker
		// based on the Execute function. For middleware usage, we rely on
		// the proxy layer to report failures through the Execute pattern.
	}
}

// Execute wraps a function with circuit breaker
func (m *CircuitBreakerManager) Execute(serviceName string, fn func() (interface{}, error)) (interface{}, error) {
	if !m.config.Enabled {
		return fn()
	}

	cb := m.GetBreaker(serviceName)
	return cb.Execute(fn)
}
