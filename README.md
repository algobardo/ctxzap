# CtxZap - Context-aware Zap Logger

[![Go Reference](https://pkg.go.dev/badge/github.com/algobardo/ctxzap.svg)](https://pkg.go.dev/github.com/algobardo/ctxzap)
[![Go Report Card](https://goreportcard.com/badge/github.com/algobardo/ctxzap)](https://goreportcard.com/report/github.com/algobardo/ctxzap)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/algobardo/ctxzap)](https://go.dev/)

CtxZap is a Go library that seamlessly integrates Uber's Zap logger with `context.Context`, allowing you to attach structured logging fields to contexts and automatically include them in all subsequent log entries.

## Features

- 🔒 **Thread-Safe**: Safe for concurrent use across goroutines
- 🎯 **Clean API**: Intuitive methods that follow Go idioms
- 🔧 **Flexible**: Compatible with existing Zap loggers and middleware

## Installation

```bash
go get github.com/algobardo/ctxzap
```

## Quick Start

```go
package main

import (
    "context"
    "github.com/algobardo/ctxzap"
    "go.uber.org/zap"
)

func main() {
    // Create a context-aware logger
    logger := ctxzap.New(zap.NewProduction())
    
    // Add fields to context
    ctx := context.Background()
    ctx = ctxzap.WithFields(ctx,
        zap.String("request_id", "abc123"),
        zap.String("user_id", "user456"),
    )
    
    // Log with context - fields are automatically included
    logger.Info(ctx, "Processing user request",
        zap.String("action", "update_profile"),
    )
    // Output: {"level":"info","msg":"Processing user request","request_id":"abc123","user_id":"user456","action":"update_profile"}
}
```

## API Reference

### Logger Creation

```go
// Create a new context-aware logger from existing zap logger
logger := ctxzap.New(zapLogger)
```

### Adding Fields to Context

```go
// Add fields to context
ctx = ctxzap.WithFields(ctx, 
    zap.String("key", "value"),
    zap.Int("count", 42),
)

// Fields are cumulative - add more fields later
ctx = ctxzap.WithFields(ctx, zap.Bool("authenticated", true))
```

### Logging with Context

```go
// All standard log levels are supported
logger.Debug(ctx, "Debug message", extraFields...)
logger.Info(ctx, "Info message", extraFields...)
logger.Warn(ctx, "Warning message", extraFields...)
logger.Error(ctx, "Error message", extraFields...)
```

### Extracting Fields

```go
// Get all fields from context (useful for middleware)
fields := ctxzap.FieldsFromContext(ctx)
```

## Comparison with Similar Libraries

### CtxZap vs Zax

| Feature | CtxZap | Zax |
|---------|---------|------|
| Store fields in context | ✅ | ✅ |
| Logger wrapper methods | ✅ | ❌ |
| Direct zap integration | ✅ | ✅ |
| Field deduplication | ✅ | ❌ |
| Zero allocation for empty ctx | ✅ | ❌ |

**When to use CtxZap over Zax:**
- You prefer methods like `logger.Info(ctx, msg)` over `logger.With(zax.Get(ctx)...).Info(msg)`

### CtxZap vs grpc-ecosystem/ctxzap

| Feature | CtxZap | grpc-ecosystem/ctxzap |
|---------|---------|----------------------|
| Store fields only | ✅ | ❌ (stores logger) |
| gRPC integration | ❌ | ✅ |
| Standalone usage | ✅ | ✅ |
| Field management | Advanced | Basic |

**When to use CtxZap over grpc-ecosystem/ctxzap:**
- You're not using gRPC or don't need gRPC-specific features
- You need more control over field management

## Advanced Usage

### HTTP Middleware Example

```go
func LoggingMiddleware(logger *ctxzap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Add request metadata to context
            ctx := r.Context()
            ctx = ctxzap.WithFields(ctx,
                zap.String("method", r.Method),
                zap.String("path", r.URL.Path),
                zap.String("request_id", generateRequestID()),
            )
            
            // Log request
            logger.Info(ctx, "Handling request")
            
            // Pass context to next handler
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Service Layer Example

```go
type UserService struct {
    logger *ctxzap.Logger
    db     *sql.DB
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // Add service-specific fields
    ctx = ctxzap.WithFields(ctx, 
        zap.String("service", "user"),
        zap.String("operation", "get_user"),
    )
    
    s.logger.Debug(ctx, "Fetching user from database")
    
    user, err := s.db.QueryUser(ctx, userID)
    if err != nil {
        s.logger.Error(ctx, "Failed to fetch user", zap.Error(err))
        return nil, err
    }
    
    s.logger.Info(ctx, "Successfully fetched user")
    return user, nil
}
```


## Best Practices

1. **Add fields early**: Add common fields (request ID, user ID) at the entry point of your application
2. **Use field keys consistently**: Establish naming conventions for field keys
3. **Don't store sensitive data**: Avoid putting passwords or tokens in context fields
4. **Keep contexts short-lived**: Don't store contexts longer than necessary

## Development

### Code Formatting and Linting

This project uses [golangci-lint](https://golangci-lint.run/) to ensure code quality and consistency. Before committing your changes, make sure to run the following:

#### Install golangci-lint

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Make sure it's in your PATH
export PATH=$PATH:$HOME/go/bin
```

#### Run Linting

```bash
# Run all linters
golangci-lint run

# Automatically fix some issues
golangci-lint run --fix
```

The project is configured with multiple linters including:
- `gofmt` - Go formatting
- `goimports` - Import formatting and organization
- `govet` - Go vet reports suspicious constructs
- `staticcheck` - Advanced static analysis
- And many more (see `.golangci.yml` for full configuration)

#### Pre-commit Hook (Optional)

To ensure code is always formatted before committing, you can set up a git pre-commit hook:

```bash
# Create pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
# Run golangci-lint before commit
golangci-lint run

# Exit with non-zero status if linting fails
if [ $? -ne 0 ]; then
    echo "Linting failed. Please fix the issues before committing."
    exit 1
fi
EOF

# Make it executable
chmod +x .git/hooks/pre-commit
```

## Contributing

Contributions are welcome! Please ensure your code passes all linting checks before submitting a Pull Request.

## License

MIT License - see LICENSE file for details