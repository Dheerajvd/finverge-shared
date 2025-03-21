package logger

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey string

const LoggerKey ctxKey = "logger"

var BaseLogger *zap.Logger

// InitLogger initializes Zap
func InitLogger(logDir string) {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	logFile := filepath.Join(logDir, "app.log")

	// Set up log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile, // Log file path
		MaxSize:    10,      // Max file size (in MB)
		MaxBackups: 7,       // Max number of old log files to retain
		MaxAge:     30,      // Max age of log files (in days)
		Compress:   true,    // Compress rotated log files
	}

	config := zapcore.EncoderConfig{
		TimeKey:      "timestamp",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),    // Log format: JSON
		zapcore.AddSync(lumberjackLogger), // Write logs to file with rotation
		zap.DebugLevel,                    // Log level
	)

	BaseLogger = zap.New(core, zap.AddCaller())
}

// GetLogger retrieves the logger from context
func GetLogger(ctx context.Context) *zap.Logger {
	if ctxLogger, ok := ctx.Value(LoggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return BaseLogger
}

// WithContext attaches a logger to context
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// CreateChildLogger returns a child logger with a component field
func CreateChildLogger(ctx context.Context, component string) context.Context {
	logger := GetLogger(ctx).With(zap.String("component", component))
	return WithContext(ctx, logger)
}

// LogChild logs a message immediately instead of storing it in memory
func LogChild(ctx context.Context, msg string, fields ...zap.Field) {
	logger := GetLogger(ctx)
	logger.Info(msg, fields...)
}
