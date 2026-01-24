#!/bin/bash

# Start Auth Service Script

export PATH=$PATH:/usr/local/go/bin

# Service Configuration
export SERVICE_NAME=auth-service
export ENVIRONMENT=development
export PORT=8081
export GRPC_PORT=9081

# Database Configuration
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres123
export DB_NAME=auth_db
export DB_SSL_MODE=disable

# Redis Configuration
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=redis123
export REDIS_DB=0

# NATS Configuration
export NATS_URL=nats://localhost:4222

# JWT Configuration
export JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-minimum-32-characters-long
export JWT_ACCESS_TOKEN_EXPIRE=15m
export JWT_REFRESH_TOKEN_EXPIRE=7d

# Logging Configuration
export LOG_LEVEL=debug
export LOG_FORMAT=console

# Run the service
cd /opt/ERP/services/auth-service
go run cmd/main.go
