# Claude Instructions

## Important Commands to Run

After completing any set of todos, please run the following commands to ensure code quality:

### 1. Run Linting
```bash
golangci-lint run ./...
```

### 2. Run Tests
```bash
go test ./...
```

Both commands must pass without errors before considering the work complete.