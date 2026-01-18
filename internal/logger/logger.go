package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a custom type for zap.Looger
type Logger struct {
	*zap.SugaredLogger
}

// Config holds logging configuration.
type Config struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	Stacktrace bool   `yaml:"stacktrace"`
}

// New creates a new zap.Logger based on the provided configuration.
func New(cfg Config) (*Logger, error) {
	zapCfg := zap.NewProductionConfig()

	levelMap := map[string]zapcore.Level{
		"debug": zap.DebugLevel,
		"info":  zap.InfoLevel,
		"warn":  zap.WarnLevel,
		"error": zap.ErrorLevel,
	}

	if !cfg.Stacktrace {
		zapCfg.DisableStacktrace = true
	}

	if level, ok := levelMap[cfg.Level]; ok {
		zapCfg.Level = zap.NewAtomicLevelAt(level)
	} else {
		return nil, fmt.Errorf("invalid log level: %s", cfg.Level)
	}

	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	logger, err := zapCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &Logger{SugaredLogger: logger.Sugar()}, nil
}

// WrapError logs an error with additional context and supports message formatting.
func (l *Logger) WrapError(msg string, err error, keysAndValues ...any) {
	l.With(zap.Error(err)).Errorw(msg, keysAndValues...)
}
