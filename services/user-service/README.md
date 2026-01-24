# User Service

User and Department Management Service for ERP Cosmetics System.

## Status

🚧 **In Development** - Foundation Complete

**Completed**:
- ✅ Database schema (3 tables)
- ✅ Migrations with seed data
- ✅ Domain entities

**To Do**:
- ⏳ Repository implementations
- ⏳ Use cases
- ⏳ HTTP API handlers
- ⏳ gRPC service
- ⏳ Auth Service integration

## Overview

User Service manages users, departments, and organizational hierarchy. It integrates with Auth Service for authentication and role management.

## Database Schema

### Tables

**departments** - Hierarchical organizational structure
- Uses materialized path pattern for efficient tree queries
- Fields: id, code, name, parent_id, manager_id, level, path, status

**users** - User information
- Fields: id, email, employee_code, first_name, last_name, phone, avatar_url, department_id, manager_id, status
- Employee code auto-generated: `EMP{YYYYMMDD}{sequence}`

**user_profiles** - Extended user information
- Fields: user_id, date_of_birth, address, emergency_contact, join_date

## API Endpoints

### Users
- `GET /api/v1/users` - List users (pagination, filters)
- `GET /api/v1/users/:id` - Get user details
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Soft delete
- `GET /api/v1/users/:id/roles` - Get user roles
- `POST /api/v1/users/:id/roles` - Assign role
- `DELETE /api/v1/users/:id/roles/:role_id` - Remove role
- `PATCH /api/v1/users/:id/status` - Change status

### Departments
- `GET /api/v1/departments` - List as tree
- `GET /api/v1/departments/:id` - Get details
- `POST /api/v1/departments` - Create
- `PUT /api/v1/departments/:id` - Update
- `DELETE /api/v1/departments/:id` - Soft delete
- `GET /api/v1/departments/:id/users` - List users

## gRPC Methods

```protobuf
service UserService {
  rpc GetUser(GetUserRequest) returns (UserResponse);
  rpc GetUsersByIds(GetUsersByIdsRequest) returns (UsersResponse);
  rpc GetUsersByDepartment(GetUsersByDepartmentRequest) returns (UsersResponse);
  rpc SearchUsers(SearchUsersRequest) returns (UsersResponse);
}
```

## Business Rules

1. **Email**: Must be unique
2. **Employee Code**: Auto-generated if not provided
3. **Department Hierarchy**: Supports nested structure
4. **Soft Delete**: Never hard delete
5. **Auth Integration**: Sync user creation with Auth Service

## Events Published

- `user.created`
- `user.updated`
- `user.deleted`
- `user.status_changed`
- `department.created`
- `department.updated`
- `department.deleted`

## Implementation Guide

### 1. Repository Layer

Create in `internal/infrastructure/persistence/postgres/`:

**user_repo.go**
```go
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
    return &userRepository{db: db}
}

// Implement: Create, GetByID, GetByEmail, Update, Delete, List, etc.
```

**department_repo.go**
```go
// Implement: Create, GetByID, GetTree, Update, Delete, GetUsers, etc.
// Special: GetTree() should return hierarchical structure
```

**user_profile_repo.go**
```go
// Implement: Create, GetByUserID, Update
```

### 2. Use Cases

Create in `internal/usecase/user/` and `internal/usecase/department/`:

**create_user.go**
```go
type CreateUserUseCase struct {
    userRepo    repository.UserRepository
    profileRepo repository.UserProfileRepository
    authClient  client.AuthClient  // gRPC to Auth Service
    eventPub    event.Publisher
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // 1. Validate input
    // 2. Generate employee code if needed
    // 3. Create user in database
    // 4. Create user profile
    // 5. Call Auth Service to create credentials
    // 6. Publish user.created event
    // 7. Return user
}
```

### 3. HTTP Handlers

Create in `internal/delivery/http/handler/`:

**user_handler.go**
```go
type UserHandler struct {
    createUC *user.CreateUserUseCase
    updateUC *user.UpdateUserUseCase
    // ... other use cases
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, err)
        return
    }
    
    user, err := h.createUC.Execute(c.Request.Context(), &req)
    if err != nil {
        response.Error(c, err)
        return
    }
    
    response.Success(c, user)
}
```

### 4. Auth Service Integration

Create `internal/infrastructure/client/auth_client.go`:

```go
type AuthClient struct {
    conn *grpc.ClientConn
}

func (c *AuthClient) CreateUserCredentials(ctx context.Context, email, password string) error {
    // Call Auth Service gRPC to create user credentials
}

func (c *AuthClient) AssignRole(ctx context.Context, userID, roleID string) error {
    // Call Auth Service gRPC to assign role
}
```

### 5. Main Application

Create `cmd/main.go`:

```go
func main() {
    // 1. Load config
    // 2. Initialize logger
    // 3. Connect to PostgreSQL
    // 4. Connect to Redis (if needed)
    // 5. Connect to NATS
    // 6. Initialize Auth Service gRPC client
    // 7. Initialize repositories
    // 8. Initialize use cases
    // 9. Initialize HTTP handlers
    // 10. Setup router
    // 11. Start HTTP server
    // 12. Start gRPC server
    // 13. Graceful shutdown
}
```

## Running Migrations

```bash
# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres123
export DB_NAME=user_db

# Run migrations
for file in migrations/*up.sql; do
  echo "Applying $file..."
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$file"
done
```

## Testing

Create test files in `test/`:
- Integration tests for repositories
- Unit tests for use cases
- API tests for HTTP handlers

## Configuration

Port: 8082 (HTTP), 9082 (gRPC)
Database: `user_db`

## Dependencies

- Auth Service (gRPC): For user credentials and role management
- PostgreSQL: Primary database
- NATS: Event publishing
- Redis: (Optional) Caching

## Next Steps

1. Implement repository layer
2. Implement use cases
3. Create HTTP handlers and DTOs
4. Implement gRPC service
5. Create Auth Service client
6. Write tests
7. Create Dockerfile and Makefile
8. Documentation

## Reference

See Auth Service implementation as template:
- `/opt/ERP/services/auth-service/`

---

**Status**: Foundation ready, implementation in progress
