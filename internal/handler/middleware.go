package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func (s *Server) prometheusMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.FullPath()

	c.Next()

	duration := time.Since(start).Seconds()
	statusCode := c.Writer.Status()

	s.config.Metrics.RequestsTotal.WithLabelValues(
		c.Request.Method,
		path,
		strconv.Itoa(statusCode),
	).Inc()

	s.config.Metrics.RequestLatency.WithLabelValues(
		c.Request.Method,
		path,
	).Observe(duration)
}

func (s *Server) tracingMiddleware(c *gin.Context) {
	tracer := opentracing.GlobalTracer()

	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request.Header),
	)

	span := tracer.StartSpan(
		c.Request.URL.Path,
		opentracing.ChildOf(spanCtx),
	)
	defer span.Finish()

	span.SetTag("http.method", c.Request.Method)
	span.SetTag("http.url", c.Request.URL.String())
	span.SetTag("component", "gin")

	c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

	c.Next()

	span.SetTag("http.status_code", c.Writer.Status())
	span.LogKV("event", "completed request")
}
