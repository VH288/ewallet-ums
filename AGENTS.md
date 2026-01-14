# AGENTS.md - Developer Guidelines for ewallet-ums

This document provides essential information for agentic coding assistants working on the ewallet-ums (User Management Service) Go microservice.

## Project Overview

ewallet-ums is a Go-based microservice providing user authentication and management functionality. It uses Gin for HTTP API, GORM for MySQL database operations, and gRPC for inter-service communication.

## Build/Lint/Test Commands

### Building
```bash
# Build the binary
go build -o ewallet-ums .

# Build with race detection
go build -race -o ewallet-ums .

# Clean build
go clean && go build -o ewallet-ums .
```

### Testing
**Note: Currently no test files exist in the project.**

```bash
# Run all tests (when tests are added)
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test file (when tests exist)
go test -v ./path/to/package

# Run a specific test function
go test -run TestFunctionName ./path/to/package
```

### Linting and Code Quality
```bash
# Format code (gofmt is built-in)
gofmt -w .

# Check formatting (shows diff if not formatted)
gofmt -d .

# Vet for potential issues
go vet ./...

# Install and run golangci-lint (if available)
golangci-lint run

# Mod tidy (clean up dependencies)
go mod tidy

# Mod verify (verify dependencies)
go mod verify
```

### Running the Application
```bash
# Run directly (requires .env file)
go run main.go

# Run with environment variables
PORT=8080 go run main.go

# Run built binary
./ewallet-ums
```

## Code Style Guidelines

### Project Structure
```
ewallet-ums/
├── cmd/                    # Application entry points and routing
├── helpers/               # Utility functions (config, logger, db, jwt, response)
├── internal/
│   ├── api/              # HTTP handlers
│   ├── interfaces/       # Interface definitions (repositories, services, handlers)
│   ├── models/           # Data structures and validation
│   ├── repository/       # Database layer implementations
│   └── services/         # Business logic layer
├── constants/            # Application constants
├── main.go              # Application entry point
├── go.mod               # Go module definition
└── .env                 # Environment configuration
```

### Naming Conventions

#### Interfaces
- Use "I" prefix: `IUserRepository`, `ILoginService`, `IRegisterHandler`
- Interface methods use PascalCase: `InsertNewUser()`, `GetUserByUsername()`

#### Structs and Types
- Use PascalCase for exported types: `User`, `LoginRequest`, `LoginResponse`
- Use camelCase for unexported types when needed

#### Functions and Methods
- Exported functions/methods: PascalCase - `SendResponseHTTP()`, `GenerateToken()`
- Unexported functions/methods: camelCase - `dependencyInject()`, `route()`

#### Variables
- camelCase for local variables: `userDetail`, `resp`, `req`
- PascalCase for exported struct fields: `ID`, `Username`, `FullName`
- Use descriptive names: prefer `userDetail` over `u`, `response` over `resp`

#### Constants
- PascalCase with descriptive prefixes when grouped: `SuccessMessage`, `ErrFailedBadRequest`
- Group related constants in const blocks

### Import Organization
```go
import (
    "context"
    "fmt"
    "time"

    "ewallet-ums/constants"
    "ewallet-ums/helpers"
    "ewallet-ums/internal/interfaces"
    "ewallet-ums/internal/models"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)
```
- Standard library imports first
- Blank line separator
- Local project imports (alphabetical order)
- Blank line separator
- Third-party imports (alphabetical order)

### Error Handling
- Always return errors from functions that can fail
- Use `fmt.Errorf()` for error wrapping: `fmt.Errorf("failed to get user: %v", err)`
- Handle errors immediately at call sites
- Log errors with context using logrus
- Return appropriate HTTP status codes with consistent error messages

### Struct Tags
- JSON tags: `json:"field_name"`
- GORM tags: `gorm:"column:column_name;type:varchar(100)"`
- Validator tags: `validate:"required"`
- Combined: `json:"username" gorm:"column:username;type:varchar(20)" validate:"required"`

### Validation
- Implement `Validate()` method on request/response structs
- Use `github.com/go-playground/validator/v10`
- Validate requests immediately after binding in API handlers
- Return appropriate error responses for validation failures

### Database Models
- Implement `TableName()` method for custom table names
- Use GORM conventions for field mapping
- Include proper GORM tags for column types and constraints
- Handle timestamps with `CreatedAt`, `UpdatedAt` fields

### API Layer (Gin Handlers)
```go
func (api *HandlerName) MethodName(c *gin.Context) {
    log := helpers.Logger
    req := models.RequestType{}

    // Bind and validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Error("failed to parse request: ", err)
        helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
        return
    }

    if err := req.Validate(); err != nil {
        log.Error("failed to validate request: ", err)
        helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
        return
    }

    // Call service
    resp, err := api.Service.Method(c.Request.Context(), req)
    if err != nil {
        log.Error("failed on service: ", err)
        helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
        return
    }

    helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}
```

### Service Layer
- Services should be stateless and depend on interfaces
- Accept `context.Context` as first parameter
- Return structured responses and errors
- Business logic should be in services, not handlers
- Use dependency injection pattern

### Repository Layer
- Implement interfaces defined in `internal/interfaces`
- Use GORM for database operations
- Accept `context.Context` for database operations
- Return domain models and errors

### Dependency Injection
- Use struct-based dependency injection in `cmd/http.go`
- Initialize all dependencies in `dependencyInject()` function
- Pass dependencies through constructor or struct fields

### Context Usage
- Always pass `context.Context` through the call chain
- Use `c.Request.Context()` in HTTP handlers
- Context is used for cancellation, timeouts, and request tracing

### Logging
- Use structured logging with logrus
- Access logger via `helpers.Logger`
- Log errors with context: `log.Error("failed to process: ", err)`
- Include relevant context in log messages

### Security Practices
- Hash passwords with bcrypt: `bcrypt.GenerateFromPassword()` and `bcrypt.CompareHashAndPassword()`
- Use JWT tokens for authentication
- Validate tokens in middleware
- Never log sensitive information (passwords, tokens)
- Use prepared statements through GORM

### Configuration
- Use environment variables for configuration
- Load config in `helpers/config.go`
- Provide sensible defaults with `helpers.GetEnv()`
- Use `.env` file for local development

### HTTP Response Format
- Use consistent response structure via `helpers.SendResponseHTTP()`
- Standard response format:
```json
{
    "message": "Success",
    "data": { ... }
}
```

### Middleware
- Implement middleware as methods on dependency struct
- Validate authentication tokens
- Check token expiration
- Set context values for downstream handlers

### gRPC Integration
- gRPC server runs alongside HTTP server
- Use `cmd.ServeGRPC()` to start gRPC server
- Implement both HTTP and gRPC handlers for services

### Testing Guidelines (When Adding Tests)
- Create test files as `*_test.go` in same package
- Use standard Go testing package
- Test functions should be named `TestXxx`
- Use table-driven tests for multiple test cases
- Mock dependencies using interfaces
- Test error conditions and edge cases
- Include integration tests for database operations

### Code Comments
- Add comments for exported functions and types
- Explain complex business logic
- Document function parameters and return values
- Use `// TODO:` for incomplete implementations
- Keep comments up-to-date with code changes

### Performance Considerations
- Use context for cancellation and timeouts
- Avoid memory leaks with proper error handling
- Use efficient database queries with GORM
- Consider connection pooling for database
- Profile performance-critical code paths

### Common Patterns
- Handler → Service → Repository pattern
- Interface-based dependency injection
- Structured error responses
- Request validation at API layer
- Context propagation throughout call stack

## Environment Variables

Required environment variables (defined in `.env`):
- `PORT`: HTTP server port (default: 8080)
- Database configuration (host, port, username, password, database name)
- JWT secrets and expiration times
- Other service-specific configuration

## Dependencies

Key dependencies (from go.mod):
- `github.com/gin-gonic/gin`: HTTP web framework
- `gorm.io/gorm`: ORM for database operations
- `gorm.io/driver/mysql`: MySQL driver for GORM
- `github.com/go-playground/validator/v10`: Input validation
- `github.com/golang-jwt/jwt/v5`: JWT token handling
- `golang.org/x/crypto/bcrypt`: Password hashing
- `github.com/sirupsen/logrus`: Structured logging
- `google.golang.org/grpc`: gRPC framework

## Development Workflow

1. Make changes to code
2. Run `go mod tidy` to clean dependencies
3. Run `gofmt -w .` to format code
4. Run `go vet ./...` to check for issues
5. Run `go build .` to verify compilation
6. Add tests for new functionality
7. Run tests: `go test ./...`
8. Commit changes following conventional commit format

## Notes for Agents

- Always check for compilation errors before making changes
- Follow existing patterns and conventions in the codebase
- Add appropriate error handling and logging
- Use interfaces for testability and dependency injection
- Keep functions small and focused on single responsibilities
- Document any deviations from these guidelines