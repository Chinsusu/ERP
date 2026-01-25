package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/erp-cosmetics/auth-service/internal/delivery/http/handler"
	"github.com/erp-cosmetics/auth-service/internal/delivery/http/dto"
	"github.com/erp-cosmetics/auth-service/internal/infrastructure/persistence/postgres"
	"github.com/erp-cosmetics/auth-service/internal/usecase/auth"
	"github.com/erp-cosmetics/shared/pkg/database"
	"github.com/erp-cosmetics/shared/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthAPI_Login_Success(t *testing.T) {
	// Setup test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_NAME", "erp_test")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("NATS_URL", "nats://localhost:4223")
	os.Setenv("JWT_SECRET", "testsecret")

	// Initialize DB
	dsn := "host=localhost port=5433 user=user password=test dbname=erp_test sslmode=disable"
	db, err := database.Connect(database.NewDefaultConfig(dsn))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize JWT manager
	jwtManager := jwt.NewManager("testsecret", 15, 60)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	permRepo := postgres.NewPermissionRepository(db)
	tokenRepo := postgres.NewTokenRepository(db)
    
    // Initialize Auth Usecases
	loginUC := auth.NewLoginUseCase(userRepo, roleRepo, permRepo, tokenRepo, nil, jwtManager, nil)
    // NOTE: Passing nil for cache and eventPub for now to see if it build. 
    // Usually we should use real ones since they are in docker.

	// Setup Handler
	authHandler := handler.NewAuthHandler(loginUC, nil, nil)
	
	// Setup Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/v1/auth/login", authHandler.Login)

	// Test Case: Valid Credentials
	loginReq := dto.LoginRequest{
		Email:    "admin@company.vn",
		Password: "Admin@123",
	}
	body, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	
	assert.Equal(t, "success", resp["status"])
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
}
