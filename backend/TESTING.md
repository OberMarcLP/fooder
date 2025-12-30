# Testing Guide

## Overview

The Nom Database implements comprehensive testing across both backend (Go) and frontend (TypeScript/React) to ensure code quality, reliability, and maintainability.

## Testing Philosophy

- **Unit Tests**: Test individual functions and components in isolation
- **Integration Tests**: Test interactions between components and systems
- **Coverage Goals**: Aim for >70% code coverage for critical paths
- **Fast Feedback**: Tests should run quickly to enable rapid development

## Backend Testing (Go)

### Running Tests

```bash
# Run all tests
make test

# Or manually
cd backend
go test ./internal/...

# Run tests with verbose output
go test ./internal/... -v

# Run tests for specific package
go test ./internal/middleware/... -v

# Run specific test
go test ./internal/middleware/... -v -run TestMetrics_RecordRequest
```

### Test Coverage

Generate coverage reports:

```bash
# Generate coverage profile
go test ./internal/... -coverprofile=coverage.out -covermode=atomic

# View coverage summary
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

View coverage by package:

```bash
# Summary by package
go test ./internal/... -cover

# Detailed coverage
go test ./internal/middleware/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```


### Test Structure

Backend tests follow Go's standard testing conventions:

```go
package middleware

import (
	"testing"
)

func TestMetrics_RecordRequest(t *testing.T) {
	// Arrange
	metrics := &Metrics{...}

	// Act
	metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)

	// Assert
	if metrics.TotalRequests != 1 {
		t.Errorf("Expected TotalRequests to be 1, got %d", metrics.TotalRequests)
	}
}
```

### Current Test Coverage

As of the latest test run:

- **Middleware Package**: 54.5% coverage
  - Logging middleware: Fully tested
  - Metrics collection: Fully tested
  - Request ID middleware: Fully tested
  - Security middleware: Needs additional tests

- **Handlers Package**: 0.1% coverage
  - Basic validation tests exist
  - Needs integration tests with database

- **Services Package**: 17.4% coverage
  - Image processor: Well tested (78.6%)
  - Google Maps service: Needs tests
  - S3 service: Needs tests

### Test Files

- `internal/middleware/logging_test.go` - Logging middleware tests
- `internal/middleware/metrics_test.go` - Metrics collection tests
- `internal/handlers/restaurants_test.go` - Restaurant handler tests
- `internal/services/imageprocessor_test.go` - Image processing tests

