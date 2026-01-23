# CHANGELOG

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project documentation structure
- README.md with project overview and quick start guide
- 01-ARCHITECTURE.md with comprehensive system architecture documentation
  - Microservices architecture overview
  - Service communication patterns (REST, gRPC, Event-driven)
  - Database per service strategy
  - API Gateway design
  - Security architecture (JWT + RBAC)
  - Docker Compose deployment architecture
  - Monitoring and observability strategy
  - Offline capability design
  - Disaster recovery plan

### Documentation
- Created comprehensive architecture diagrams using Mermaid
- Documented all 15 microservices with port allocations
- Defined security patterns and encryption strategies
- Specified resource allocation for Docker containers

---

## [1.0.0] - 2026-01-23

### Initial Release
- Project initialization
- Documentation framework setup

## [Unreleased] - 2026-01-23

### Added
- 02-SERVICE-SPECIFICATIONS.md - Comprehensive specifications for all 15 microservices
  - API Gateway routing and middleware configuration
  - Auth Service with JWT and RBAC implementation
  - User Service with department management
  - Master Data Service for materials and products (with INCI, CAS numbers)
  - Supplier Service with certification tracking (GMP, ISO, Organic)
  - Procurement Service with PR/PO workflows
  - WMS Service with lot management and FEFO logic
  - Manufacturing Service with BOM versioning, work orders, and QC
  - Sales Service with customer and order management
  - Summary of remaining services (Marketing, Finance, Reporting, Notification, AI, File)

### Documentation
- Detailed database schemas for each service
- REST API endpoints with HTTP methods
- gRPC internal communication methods
- Event publishing/subscription patterns
- Permission requirements per service
- Service dependency mapping with Mermaid diagram

### Phase 2 Completed - 2026-01-23

**Added:**
- 03-AUTH-SERVICE.md - Complete authentication & authorization documentation
  - JWT access & refresh token flows with sequence diagrams
  - RBAC permission system (format: service:resource:action)
  - User-role and role-permission management
  - Account security (bcrypt, lockout, token rotation)
  - Redis caching strategy for permissions
  - Comprehensive API endpoints and gRPC methods
  
- 04-USER-SERVICE.md - User & department management
  - Employee information management
  - Department hierarchy (tree structure)
  - User preferences and avatar handling
  - Employment status tracking
  - Emergency contact management
  
- 05-MASTER-DATA-SERVICE.md - Materials & products master data
  - Material management with INCI names & CAS numbers
  - Product master data with cosmetic licenses
  - Category hierarchy for materials and products
  - Units of measure with conversion logic
  - Material-supplier approved lists
  - Allergen tracking and storage requirements
  
- 13-API-GATEWAY.md - API Gateway implementation
  - Complete middleware chain (CORS, logging, rate limiting, auth)
  - Routing configuration for all 15 services
  - Circuit breaker pattern for fault tolerance
  - Rate limiting strategy (100 req/min per user)
  - Permission-based route protection
  - Error handling and monitoring metrics
