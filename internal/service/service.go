package service

import (
	"html/template"
	"international_site/internal/config"
	"international_site/internal/logger"
	"international_site/internal/storage/lts"
	"international_site/internal/storage/models"
	"international_site/internal/types"
	"time"
)

type Protocol interface {
	GetPageBySlug(slug, locale string) (*models.Page, error)
	GetProductCategories(locale string) ([]models.ProductCategory, error)
	GetCategoryBySlug(slug, locale string) (*models.ProductCategory, error)
	GetCategoryByID(locale string, id uint) (*models.ProductCategory, error)
	GetProducts(locale string, offset, limit int) ([]models.Product, int, error)
	GetProductsByCategory(locale string, categoryID uint, offset, limit int) ([]models.Product, int, error)
	GetProductByID(locale string, id uint) (*models.Product, error)
	GetRelatedProducts(locale string, categoryID, excludeID uint, limit int) ([]models.Product, error)
	GetNews(locale string, offset, limit int) ([]models.News, int, error)
	GetNewsByID(locale string, id uint) (*models.News, error)
	GetRecentNews(locale string, limit int) ([]models.News, error)
	GetDocumentsByType(docType, locale string) ([]models.Document, error)
	GetContactsByType(contactType, locale string) ([]models.Contact, error)
	SearchProducts(locale, query string, offset, limit int) ([]models.Product, int, error)
	SearchAll(locale, query string) ([]types.SearchResult, error)
	SearchAPI(locale, query string, limit int) ([]types.SearchResult, error)
	FilterProducts(locale string, categoryID uint, search, sortBy, sortOrder string, offset, limit int) ([]models.Product, int, error)
	GetProductsSorted(locale, search, sortBy, sortOrder string, offset, limit int) ([]models.Product, int, error)
	SaveFeedback(feedback types.FeedbackRequest) (uint, error)
	GetTranslation(key, locale string) string
	GenerateSitemap(locale string) (string, error)
	GetAvailableLanguages() []string
}

type Instance struct {
	logger    *logger.Logger
	lts       lts.Protocol
	cfg       *config.Service
	NowFunc   func() time.Time
	templates map[string]*template.Template
}

func New(
	logger *logger.Logger,
	lts lts.Protocol,
	cfg *config.Service,
	now func() time.Time,
) *Instance {
	return &Instance{
		logger:  logger,
		lts:     lts,
		cfg:     cfg,
		NowFunc: now,
	}
}
