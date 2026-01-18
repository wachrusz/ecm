package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"international_site/internal/logger"
)

// EventLoggerProtocol defines the interface for logging events and errors.
type EventLoggerProtocol interface {
	LogEvent(c *gin.Context, event ...interface{})
	LogError(c *gin.Context)
	StartSpanFromContext(c *gin.Context)
	FinishSpanFromContext(c *gin.Context)
}

// EventLogger implements EventLoggerProtocol using OpenTracing.
type EventLogger struct {
	Log *logger.Logger
}

// NewEventLogger returns new EventLogger.
func NewEventLogger(log *logger.Logger) *EventLogger {
	return &EventLogger{Log: log}
}

// LogEvent logs an event to the current span.
func (l *EventLogger) LogEvent(c *gin.Context, event ...interface{}) {
	if span := opentracing.SpanFromContext(c.Request.Context()); span != nil {
		span.LogKV("event", event)
	}

	sumMsg := ""

	for _, e := range event {
		msg, ok := e.(string)
		if !ok {
			continue
		}

		sumMsg += msg
	}

	if sumMsg != "" {
		l.Log.Info(sumMsg)
	}
}

// LogError marks the current span as an error.
func (l *EventLogger) LogError(c *gin.Context) {
	if span := opentracing.SpanFromContext(c.Request.Context()); span != nil {
		span.SetTag("error", true)
	}
}

// StartSpanFromContext starts a tracing span from context
func (l *EventLogger) StartSpanFromContext(c *gin.Context) {
	opentracing.StartSpanFromContext(c, "ws_details")
}

// FinishSpanFromContext finishes span using context
func (l *EventLogger) FinishSpanFromContext(c *gin.Context) {
	if span := opentracing.SpanFromContext(c); span != nil {
		span.Finish()
	}
}
