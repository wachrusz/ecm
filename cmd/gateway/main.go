package main

import (
	"fmt"
	"international_site/internal/config"
	"international_site/internal/handler"
	"international_site/internal/logger"
	"international_site/internal/service"
	"international_site/internal/storage/lts"
	"international_site/pkg/metrics"
	"international_site/pkg/tracing"
	"log"
	"net/http"
	_ "net/http/pprof" // nolint: gosec
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
	"go.uber.org/zap"
)

func getConfigPath() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	log.Printf("loaded config: app_configs/config_%s.yaml", env)

	return fmt.Sprintf("app_configs/config_%s.yaml", env)
}

func nowFunc() time.Time {
	return time.Now()
}

// NOTE: don't hardcode host due cors policy on server!

// @version 1.0
// @BasePath /api/v1
// @schemes http https ws wss
func main() {
	cfg, err := config.Load(getConfigPath())

	if err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.Logger)

	if err != nil {
		logger.Panic("panic", zap.Error(err))
	}

	logger.Warn("Loaded config", zap.String("config", cfg.ToString()))

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Panic("panic", zap.Error(err))
		}
	}()

	lts, err := lts.New(cfg.Database.Master)

	if err != nil {
		logger.Panic("panic", zap.Error(err))
	}

	metrics := metrics.NewMetrics()
	eventLogger := tracing.NewEventLogger(logger)
	tracer, closer, err := config.InitTracer(cfg.Tracer)

	defer closer.Close()

	if err != nil {
		logger.Panic("panic", zap.Error(err))
	}

	serverCfg := config.NewServer(cfg.Server, tracer, metrics, eventLogger)
	router := gin.Default()
	service := service.New(logger, lts, &cfg.Service, nowFunc)

	server := handler.New(service, serverCfg, router, logger)

	if err := server.ListenAndServe(); err != nil {
		logger.Panic("panic", zap.Error(err))
	}
}

// nolint: revive
func CheckOrigin(r *http.Request) bool {
	return true
}
