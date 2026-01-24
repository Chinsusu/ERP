package router

import (
	"strings"
	"time"

	"github.com/erp-cosmetics/api-gateway/internal/config"
	"github.com/erp-cosmetics/api-gateway/internal/health"
	"github.com/erp-cosmetics/api-gateway/internal/middleware"
	"github.com/erp-cosmetics/api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Router handles HTTP routing
type Router struct {
	cfg            *config.Config
	logger         *zap.Logger
	redis          *redis.Client
	registry       *proxy.ServiceRegistry
	proxyHandler   *proxy.Handler
	healthHandler  *health.Handler
	rateLimiter    *middleware.RateLimiter
	cbManager      *middleware.CircuitBreakerManager
}

// NewRouter creates a new router
func NewRouter(
	cfg *config.Config,
	logger *zap.Logger,
	redisClient *redis.Client,
) *Router {
	// Create service registry
	registry := proxy.NewServiceRegistry(logger)

	// Register services from routes
	registeredServices := make(map[string]bool)
	for _, route := range cfg.Routes {
		serviceName := extractServiceName(route.Service)
		if !registeredServices[serviceName] {
			registry.Register(serviceName, route.Service)
			registeredServices[serviceName] = true
		}
	}

	// Create handlers
	proxyHandler := proxy.NewHandler(registry, logger, cfg.DefaultTimeout)
	healthHandler := health.NewHandler(registry, "1.0.0")

	// Create rate limiter
	rateLimiter := middleware.NewRateLimiter(redisClient, middleware.RateLimiterConfig{
		Enabled:        cfg.RateLimitEnabled,
		RequestsPerMin: cfg.RateLimitPerMin,
		Burst:          cfg.RateLimitBurst,
		IPLimitPerMin:  cfg.IPRateLimitPerMin,
	})

	// Create circuit breaker manager
	cbManager := middleware.NewCircuitBreakerManager(middleware.CircuitBreakerConfig{
		Enabled:   cfg.CircuitBreakerEnabled,
		Threshold: cfg.CircuitBreakerThreshold,
		Timeout:   cfg.CircuitBreakerTimeout,
	})

	return &Router{
		cfg:           cfg,
		logger:        logger,
		redis:         redisClient,
		registry:      registry,
		proxyHandler:  proxyHandler,
		healthHandler: healthHandler,
		rateLimiter:   rateLimiter,
		cbManager:     cbManager,
	}
}

// Setup configures the Gin router
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// Global middlewares (order matters!)
	engine.Use(middleware.RequestID())
	engine.Use(middleware.CORS(middleware.DefaultCORSConfig()))
	engine.Use(middleware.Logger(r.logger))
	engine.Use(middleware.Recovery(r.logger))

	// Health check endpoints (no auth, no rate limit)
	engine.GET("/health", r.healthHandler.Health)
	engine.GET("/ready", r.healthHandler.Ready)
	engine.GET("/live", r.healthHandler.Live)
	engine.GET("/health/:service", r.healthHandler.ServiceHealth)

	// Setup route groups based on configuration
	for _, route := range r.cfg.Routes {
		r.setupRoute(engine, route)
	}

	return engine
}

// setupRoute configures a single route
func (r *Router) setupRoute(engine *gin.Engine, route config.RouteConfig) {
	serviceName := extractServiceName(route.Service)

	// Build middleware chain
	handlers := []gin.HandlerFunc{}

	// Rate limiting (before auth so unauthenticated requests are also limited)
	handlers = append(handlers, r.rateLimiter.Middleware())

	// Authentication
	if route.AuthRequired {
		handlers = append(handlers, middleware.Auth(middleware.AuthConfig{
			JWTSecret: r.cfg.JWTSecret,
			Redis:     r.redis,
		}))
	}

	// Circuit breaker
	handlers = append(handlers, r.cbManager.Middleware(serviceName))

	// Proxy handler
	handlers = append(handlers, r.proxyHandler.Proxy(serviceName))

	// Register route for all methods
	group := engine.Group(route.Prefix)
	{
		group.Any("", handlers...)
		group.Any("/*path", handlers...)
	}
}

// extractServiceName extracts service name from host:port
func extractServiceName(service string) string {
	parts := strings.Split(service, ":")
	return parts[0]
}

// StartHealthChecks starts periodic health checks
func (r *Router) StartHealthChecks() {
	r.registry.StartHealthChecks(30 * time.Second)
}
