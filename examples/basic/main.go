package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/algobardo/ctxzap"
	"go.uber.org/zap"
)

func main() {
	// Create a production logger
	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()

	// Create a context-aware logger
	logger := ctxzap.New(zapLogger)

	// Example 1: Basic usage
	fmt.Println("=== Example 1: Basic Usage ===")
	basicExample(logger)

	// Example 2: HTTP middleware
	fmt.Println("\n=== Example 2: HTTP Middleware ===")
	httpExample(logger)

	// Example 3: Service layers
	fmt.Println("\n=== Example 3: Service Layers ===")
	serviceExample(logger)
}

func basicExample(logger *ctxzap.Logger) {
	ctx := context.Background()

	// Add fields to context
	ctx = ctxzap.WithFields(ctx,
		zap.String("request_id", "abc123"),
		zap.String("user_id", "user456"),
	)

	// Log with context - fields are automatically included
	logger.Info(ctx, "Processing user request",
		zap.String("action", "update_profile"),
	)

	// Add more fields later
	ctx = ctxzap.WithFields(ctx,
		zap.String("session_id", "session789"),
		zap.Bool("authenticated", true),
	)

	// All fields are included
	logger.Info(ctx, "User profile updated successfully")

	// Override existing field
	ctx = ctxzap.WithFields(ctx,
		zap.String("action", "send_notification"), // This overrides the previous "action"
	)

	logger.Info(ctx, "Sending notification")
}

func httpExample(logger *ctxzap.Logger) {
	// Create a simple HTTP handler with logging middleware
	handler := loggingMiddleware(logger)(http.HandlerFunc(handleRequest))

	// Simulate a request
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "/api/users/123", http.NoBody)
	req.Header.Set("X-Request-ID", "req-456")

	// Create a mock response writer
	handler.ServeHTTP(&mockResponseWriter{}, req)
}

func loggingMiddleware(logger *ctxzap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Extract request ID from header or generate one
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
			}

			// Add request metadata to context
			ctx := r.Context()
			ctx = ctxzap.WithFields(ctx,
				zap.String("request_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
			)

			// Log request start
			logger.Info(ctx, "Request started")

			// Pass context to next handler
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

			// Log request completion
			logger.Info(ctx, "Request completed",
				zap.Duration("duration", time.Since(start)),
				zap.Int("status", 200), // In real code, capture actual status
			)
		})
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract logger fields from context if needed
	fields := ctxzap.FieldsFromContext(ctx)
	fmt.Printf("Request has %d context fields\n", len(fields))

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func serviceExample(logger *ctxzap.Logger) {
	// Create services
	userService := &UserService{logger: logger}
	notificationService := &NotificationService{logger: logger}

	// Create a context with initial fields
	ctx := context.Background()
	ctx = ctxzap.WithFields(ctx,
		zap.String("request_id", "svc-123"),
		zap.String("correlation_id", "corr-456"),
	)

	// Simulate a business operation
	user := userService.GetUser(ctx, "user789")

	notificationService.SendWelcomeEmail(ctx, user)

	logger.Info(ctx, "Business operation completed successfully")
}

// UserService demonstrates service-layer logging
type UserService struct {
	logger *ctxzap.Logger
}

type User struct {
	ID    string
	Email string
	Name  string
}

func (s *UserService) GetUser(ctx context.Context, userID string) *User {
	// Add service-specific fields
	ctx = ctxzap.WithFields(ctx,
		zap.String("service", "user"),
		zap.String("operation", "get_user"),
		zap.String("user_id", userID),
	)

	s.logger.Debug(ctx, "Fetching user from database")

	// Simulate database query
	time.Sleep(5 * time.Millisecond)

	user := &User{
		ID:    userID,
		Email: "user@example.com",
		Name:  "John Doe",
	}

	s.logger.Info(ctx, "Successfully fetched user",
		zap.String("user_email", user.Email),
	)

	return user
}

// NotificationService demonstrates cross-service context propagation
type NotificationService struct {
	logger *ctxzap.Logger
}

func (s *NotificationService) SendWelcomeEmail(ctx context.Context, user *User) {
	// Add notification service fields
	ctx = ctxzap.WithFields(ctx,
		zap.String("service", "notification"),
		zap.String("operation", "send_welcome_email"),
		zap.String("recipient", user.Email),
	)

	s.logger.Info(ctx, "Sending welcome email")

	// Simulate email sending
	time.Sleep(10 * time.Millisecond)

	s.logger.Info(ctx, "Welcome email sent successfully",
		zap.String("template", "welcome_v2"),
	)
}

// mockResponseWriter is a simple mock for demonstration
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {}
