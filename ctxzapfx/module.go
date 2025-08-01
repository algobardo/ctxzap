package ctxzapfx

import (
	"github.com/algobardo/ctxzap"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides an fx.Module that automatically converts a *zap.Logger
// to a *ctxzap.Logger for dependency injection.
var Module = fx.Module("ctxzapfx",
	fx.Provide(NewLogger),
)

// NewLogger creates a new ctxzap.Logger from a zap.Logger.
// This function is used by the fx module to provide the wrapper logger.
func NewLogger(zapLogger *zap.Logger) *ctxzap.Logger {
	return ctxzap.New(zapLogger)
}