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

# Build with specific Go version (requires Go 1.25.1+)
go build -o ewallet-ums .
```

### Testing
**Note: Currently no test files exist in the project. Use these commands when tests are added.**

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test -v ./internal/services

# Run a specific test function
go test -run TestFunctionName ./internal/services

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Generate mocks (when tests are added)
go generate ./...
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

# Download dependencies
go mod download
```

### Running the Application
```bash
# Run directly (requires .env file)
go run main.go

# Run with environment variables
PORT=8080 go run main.go

# Run built binary
./ewallet-ums

# Run with specific config
DB_HOST=localhost DB_PORT=3306 ./ewallet-ums
```

## Code Style Guidelines

### Project Structure
```
ewallet-ums/
├── cmd/                    # Application entry points and routing
│   ├── proto/             # gRPC protobuf definitions and generated files
│   ├── dependency.go      # Dependency injection setup
│   ├── grpc.go           # gRPC server setup
│   ├── http.go           # HTTP server setup
│   ├── middleware.go     # HTTP middleware
│   └── route.go          # HTTP route definitions
├── helpers/               # Utility functions (config, logger, db, jwt, response)
├── internal/
│   ├── api/              # HTTP handlers
│   ├── interfaces/       # Interface definitions (repositories, services, handlers)
│   ├── models/           # Data structures and validation
│   ├── repository/       # Database layer implementations
│   └── services/         # Business logic layer
├── constants/            # Application constants
├── main.go              # Application entry point
├── go.mod               # Go module definition (Go 1.25.1)
├── go.sum               # Go module checksums
└── .env                 # Environment configuration (never commit)
```

### Naming Conventions

#### Interfaces
- Use "I" prefix for all interfaces: `IUserRepository`, `ILoginService`, `IRegisterHandler`
- Interface methods use PascalCase: `InsertNewUser()`, `GetUserByUsername()`

#### Structs and Types
- Use PascalCase for exported types: `User`, `LoginRequest`, `LoginResponse`
- Use camelCase for unexported types when needed
- Handler structs use PascalCase: `LoginHandler`, `RegisterHandler`

#### Functions and Methods
- Exported functions/methods: PascalCase - `SendResponseHTTP()`, `GenerateToken()`
- Unexported functions/methods: camelCase - `dependencyInject()`, `route()`
- Constructor functions: `NewUserRepository()`, `NewLoginService()`

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
- Standard library imports first (alphabetical)
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
- Use structured error responses

### Struct Tags
- JSON tags: `json:"field_name"`
- GORM tags: `gorm:"column:column_name;type:varchar(100)"`
- Validator tags: `validate:"required"`
- Combined: `json:"username" gorm:"column:username;type:varchar(20)" validate:"required"`
- Omit sensitive fields: `json:"password,omitempty"`

### Validation
- Implement `Validate()` method on request/response structs
- Use `github.com/go-playground/validator/v10`
- Validate requests immediately after binding in API handlers
- Return appropriate error responses for validation failures
- Use struct tags for validation rules

### Database Models
- Implement `TableName()` method for custom table names
- Use GORM conventions for field mapping
- Include proper GORM tags for column types and constraints
- Handle timestamps with `CreatedAt`, `UpdatedAt` fields (GORM automatic)
- Use appropriate data types for MySQL compatibility

### API Layer (Gin Handlers)
```go
func (api *LoginHandler) Login(c *gin.Context) {
    log := helpers.Logger
    req := models.LoginRequest{}

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
    resp, err := api.LoginService.Login(c.Request.Context(), req)
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
- Services coordinate between repositories and business logic

### Repository Layer
- Implement interfaces defined in `internal/interfaces`
- Use GORM for database operations
- Accept `context.Context` for database operations
- Return domain models and errors
- Handle database-specific operations and queries

### Dependency Injection
- Use struct-based dependency injection in `cmd/dependency.go`
- Initialize all dependencies in `dependencyInject()` function
- Pass dependencies through constructor or struct fields
- Follow clean architecture principles

### Context Usage
- Always pass `context.Context` through the call chain
- Use `c.Request.Context()` in HTTP handlers
- Context is used for cancellation, timeouts, and request tracing
- Include context in all service and repository calls

### Logging
- Use structured logging with logrus
- Access logger via `helpers.Logger`
- Log errors with context: `log.Error("failed to process: ", err)`
- Include relevant context in log messages
- Log at appropriate levels (Error, Warn, Info, Debug)

### Security Practices
- Hash passwords with bcrypt: `bcrypt.GenerateFromPassword()` and `bcrypt.CompareHashAndPassword()`
- Use JWT tokens for authentication with proper expiration
- Validate tokens in middleware
- Never log sensitive information (passwords, tokens, secrets)
- Use prepared statements through GORM
- Validate all user inputs

### Configuration
- Use environment variables for configuration via `.env` file
- Load config in `helpers/config.go`
- Provide sensible defaults with `helpers.GetEnv(key, defaultValue)`
- Use `.env` file for local development (never commit)
- Support runtime environment variable overrides

### HTTP Response Format
- Use consistent response structure via `helpers.SendResponseHTTP()`
- Standard response format:
```json
{
    "message": "Success",
    "data": { ... }
}
```
- Error responses follow same format with appropriate status codes

### Middleware
- Implement middleware as methods on dependency struct
- Validate authentication tokens in middleware
- Check token expiration and validity
- Set context values for downstream handlers
- Handle CORS, logging, and security headers

### gRPC Integration
- gRPC server runs alongside HTTP server
- Use `cmd.ServeGRPC()` to start gRPC server
- Implement both HTTP and gRPC handlers for services
- Protobuf definitions in `cmd/proto/`
- Generated Go code from protobuf files

### Testing Guidelines (When Adding Tests)
- Create test files as `*_test.go` in same package
- Use standard Go testing package
- Test functions should be named `TestXxx`
- Use table-driven tests for multiple test cases
- Mock dependencies using interfaces and go.uber.org/mock
- Test error conditions and edge cases
- Include integration tests for database operations
- Use `go generate` for mock generation

### Code Comments
- Add comments for exported functions and types
- Explain complex business logic
- Document function parameters and return values
- Use `// TODO:` for incomplete implementations
- Keep comments up-to-date with code changes
- Comment security-sensitive code

### Performance Considerations
- Use context for cancellation and timeouts
- Avoid memory leaks with proper error handling
- Use efficient database queries with GORM
- Consider connection pooling for database
- Profile performance-critical code paths
- Use appropriate data types and structures

### Common Patterns
- Handler → Service → Repository pattern
- Interface-based dependency injection
- Structured error responses
- Request validation at API layer
- Context propagation throughout call stack
- Clean architecture with separation of concerns

## Environment Variables

Required environment variables (defined in `.env`):
- `PORT`: HTTP server port (default: 8080)
- Database configuration: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- JWT configuration: `JWT_SECRET`, `JWT_EXPIRATION`, `REFRESH_TOKEN_EXPIRATION`
- Other service-specific configuration

**Security Note**: Never commit `.env` files to version control.

## Dependencies

Key dependencies (from go.mod):
- `github.com/gin-gonic/gin`: HTTP web framework
- `github.com/go-playground/validator/v10`: Input validation
- `github.com/golang-jwt/jwt/v5`: JWT token handling
- `github.com/joho/godotenv`: Environment variable loading
- `github.com/sirupsen/logrus`: Structured logging
- `golang.org/x/crypto/bcrypt`: Password hashing
- `google.golang.org/grpc`: gRPC framework
- `gorm.io/gorm`: ORM for database operations
- `gorm.io/driver/mysql`: MySQL driver for GORM
- `go.uber.org/mock`: Mock generation for testing

## Development Workflow

1. Set up environment: copy `.env.example` to `.env` and configure
2. Install dependencies: `go mod download`
3. Run tests (when available): `go test ./...`
4. Format code: `gofmt -w .`
5. Vet code: `go vet ./...`
6. Build application: `go build -o ewallet-ums .`
7. Run application: `./ewallet-ums`
8. Add tests for new functionality
9. Commit changes following conventional commit format

## Notes for Agents

- Always check for compilation errors before making changes
- Follow existing patterns and conventions in the codebase
- Add appropriate error handling and logging
- Use interfaces for testability and dependency injection
- Keep functions small and focused on single responsibilities
- Document any deviations from these guidelines
- Respect security practices: never expose secrets or sensitive data
- Use the existing dependency injection pattern
- Follow the Handler → Service → Repository architecture</content>
<parameter name="filePath">/home/dodoko/Work/Dodo/Go/Go Project/Project/ecom/ewallet-ums/AGENTS.md