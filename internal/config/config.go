package config

import (
	"fmt"
	"international_site/internal/logger"
	"international_site/pkg/metrics"
	"international_site/pkg/tracing"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"gopkg.in/yaml.v3"
)

const (
	censorship = "***"
)

// App holds the application configuration.
type App struct {
	Server   Server        `yaml:"server"`
	Service  Service       `yaml:"service"`
	Database LTS           `yaml:"lts"`
	Logger   logger.Config `yaml:"logger"`
	Tracer   Tracer        `yaml:"tracer"`
}

// ToString returns a string representation of the configuration with sensitive data (passwords) masked
func (a *App) ToString() string {
	cfgCopy := *a

	cfgCopy.Database.Master.Host = censorship
	cfgCopy.Database.Master.Password = censorship
	cfgCopy.Database.Master.Username = censorship

	return fmt.Sprintf("%+v", cfgCopy)
}

// Server holds server-related configuration.
type Server struct {
	Host           string        `yaml:"host"`
	Port           string        `yaml:"port"`
	BaseURL        string        `yaml:"base_url"`
	ReadTimeout    time.Duration `yaml:"read_timeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
	Timezone       string        `yaml:"timezone"`
	UsernameLength int           `yaml:"min_username_length"`

	Metrics     *metrics.Metrics
	Tracer      opentracing.Tracer
	EventLogger tracing.EventLoggerProtocol
}

type Service struct {
	DefaultLang string
	SiteURL     string
}

// Tracer holds tracing configuration details.
type Tracer struct {
	ServiceName       string        `yaml:"service_name"`
	SamplerType       string        `yaml:"sampler_type"`
	SamplerParam      float64       `yaml:"sampler_param"`
	ReporterHost      string        `yaml:"reporter_host"`
	ReporterPort      string        `yaml:"reporter_port"`
	CollectorEndpoint string        `yaml:"collector_endpoint"`
	LogSpans          bool          `yaml:"reporter_log_spans"`
	BufferFlush       time.Duration `yaml:"reporter_buffer_flush_interval"`
	Disabled          bool          `yaml:"disabled"`
}

// LTS contains database connection settings.
type LTS struct {
	Master  LTSInstance `yaml:"master"`
	Replica LTSInstance `yaml:"replica"`
	Pool    LTSPool     `yaml:"pool"`
}

// LTSInstance represents a database instance.
type LTSInstance struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Sslmode  string `yaml:"sslmode"`
}

// LTSPool manages a pool of database connections.
type LTSPool struct {
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

// Load loads configuration from the specified path.
func Load(path string) (*App, error) {
	cfg := &App{}

	file, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err = yaml.Unmarshal(file, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}

// GetDSN returns the Data Source Name for the database instance.
func (d *LTSInstance) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		d.Host,
		d.Username,
		d.Password,
		d.Dbname,
		d.Port,
		d.Sslmode,
	)
}

// InitTracer initializes the tracer using the provided configuration.
//
//nolint:ireturn
func InitTracer(trCfg Tracer) (opentracing.Tracer, io.Closer, error) {
	if trCfg.Disabled {
		log.Println("tracer disabled")

		return opentracing.NoopTracer{}, io.NopCloser(nil), nil
	}

	localAgentHostPort := ""

	if trCfg.CollectorEndpoint == "" {
		localAgentHostPort = trCfg.ReporterHost + ":" + trCfg.ReporterPort
	}

	cfg := jaegercfg.Configuration{
		ServiceName: trCfg.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  trCfg.SamplerType,
			Param: trCfg.SamplerParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            trCfg.LogSpans,
			BufferFlushInterval: trCfg.BufferFlush,
			LocalAgentHostPort:  localAgentHostPort,
			CollectorEndpoint:   trCfg.CollectorEndpoint,
		},
	}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))

	if err != nil {
		return nil, nil, err
	}

	return tracer, closer, nil
}

// NewServer initializes the configs using the provided configuration.
func NewServer(
	cfg Server,
	tracer opentracing.Tracer,
	metrics *metrics.Metrics,
	eventLogger tracing.EventLoggerProtocol) *Server {
	return &Server{
		Host:         cfg.Host,
		Port:         cfg.Port,
		BaseURL:      cfg.BaseURL,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Tracer:       tracer,
		Metrics:      metrics,
		EventLogger:  eventLogger,
		Timezone:     cfg.Timezone,
	}
}
