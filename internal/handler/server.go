package handler

import (
	"international_site/internal/config"
	"international_site/internal/logger"
	"international_site/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ServerProtocol interface for a server
type ServerProtocol interface {
	ListenAndServe(host string, baseURL string, port string) error
}

// Server is implementation of ServerProtocol
type Server struct {
	config *config.Server
	logger *logger.Logger
	service service.Protocol
	router  *gin.Engine
	nowFunc func() time.Time
}

// New returns new Server instance
func New(
	service service.Protocol,
	config *config.Server,
	router *gin.Engine,
	logger *logger.Logger,
) *Server {
	return &Server{
		service: service,
		router:  router,
		config:  config,
		logger:  logger,
		nowFunc: nowFunc,
	}
}

func nowFunc() time.Time {
	return time.Now()
}

// ListenAndServe opens host to listen
func (s *Server) ListenAndServe() error {
	opentracing.SetGlobalTracer(s.config.Tracer)
	s.router.Use(s.prometheusMiddleware)
	s.registerRoutes(s.config.BaseURL)

	srv := &http.Server{
		Addr:         s.config.Host + ":" + s.config.Port,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	return srv.ListenAndServe()
}

func (s *Server) registerRoutes(baseURL string) {
	s.router.StaticFS("/swaggerFiles", http.Dir("./docs"))

	url := ginSwagger.URL("/swaggerFiles/swagger.json")

	r := s.router.Group(baseURL)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.Static("/swaggerFiles", "./docs")

	s.router.Static("/static", "./static")
	s.router.Static("/uploads", "./uploads")

	s.router.GET("/:locale/page/:slug", s.Page)

	s.router.GET("/:locale/about", s.AboutPage)
	s.router.GET("/:locale/certificates", s.CertificatesPage)

	s.router.GET("/:locale/products", s.ProductsPage)
	s.router.GET("/:locale/products/:category", s.ProductsByCategory)
	s.router.GET("/:locale/product/:id", s.ProductDetail)

	s.router.GET("/:locale/news", s.NewsPage)
	s.router.GET("/:locale/news/:id", s.NewsDetail)
	s.router.GET("/:locale/documents", s.DocumentsPage)

	s.router.GET("/:locale/contacts", s.ContactsPage)
	s.router.POST("/:locale/feedback", s.SubmitFeedback)

	s.router.GET("/:locale/search", s.SearchPage)
	s.router.GET("/:locale/sitemap.xml", s.Sitemap)
	s.router.GET("/:locale/privacy", s.PrivacyPage)

	api := s.router.Group("/api/:locale")
	{
		api.GET("/search", s.APISearch)
		api.GET("/products/filter", s.APIProductsFilter)
		api.POST("/feedback", s.APISubmitFeedback)
	}

	s.router.NoRoute(s.NotFoundPage)

	r.Use(s.tracingMiddleware)

	s.registerDebugHandlers(baseURL)
}

func (s *Server) registerDebugHandlers(baseURL string) {
	debugGroup := s.router.Group(baseURL + "/debug")
	{
		debugGroup.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	pprofGroup := s.router.Group("/debug/pprof")
	{
		pprofGroup.GET("/*any", gin.WrapH(http.DefaultServeMux))
	}
}
