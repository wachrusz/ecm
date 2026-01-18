package service

import (
	"fmt"
	"html/template"
	"international_site/internal/config"
	"international_site/internal/storage/lts"
	"international_site/internal/storage/models"
	"international_site/internal/types"
	"strings"
	"time"
)

func NewService(storage lts.Protocol, cfg *config.Service) *Instance {
	s := &Instance{
		lts: storage,
		cfg: cfg,
	}

	s.loadTemplates()
	return s
}

func (i *Instance) loadTemplates() {
	i.templates = make(map[string]*template.Template)
}

func (i *Instance) GetPageBySlug(slug, locale string) (*models.Page, error) {
	return i.lts.GetPageBySlug(slug, locale)
}

func (i *Instance) GetProductCategories(locale string) ([]models.ProductCategory, error) {
	return i.lts.GetCategories(locale)
}

func (i *Instance) GetCategoryBySlug(slug, locale string) (*models.ProductCategory, error) {
	return i.lts.GetCategoryBySlug(slug, locale)
}

func (i *Instance) GetCategoryByID(locale string, id uint) (*models.ProductCategory, error) {
	return i.lts.GetCategoryByID(id, locale)
}

func (i *Instance) GetProducts(locale string, offset, limit int) ([]models.Product, int, error) {
	return i.lts.GetProducts(locale, offset, limit)
}

func (i *Instance) GetProductsByCategory(locale string, categoryID uint, offset, limit int) ([]models.Product, int, error) {
	return i.lts.GetProductsByCategory(locale, categoryID, offset, limit)
}

func (i *Instance) GetProductByID(locale string, id uint) (*models.Product, error) {
	return i.lts.GetProductByID(id, locale)
}

func (i *Instance) GetRelatedProducts(locale string, categoryID, excludeID uint, limit int) ([]models.Product, error) {
	return i.lts.GetRelatedProducts(locale, categoryID, excludeID, limit)
}

func (i *Instance) GetNews(locale string, offset, limit int) ([]models.News, int, error) {
	return i.lts.GetNews(locale, offset, limit)
}

func (i *Instance) GetNewsByID(locale string, id uint) (*models.News, error) {
	return i.lts.GetNewsByID(id, locale)
}

func (i *Instance) GetRecentNews(locale string, limit int) ([]models.News, error) {
	news, _, err := i.lts.GetNews(locale, 0, limit)
	return news, err
}

func (i *Instance) GetDocumentsByType(docType, locale string) ([]models.Document, error) {
	return i.lts.GetDocumentsByType(docType, locale)
}

func (i *Instance) GetContactsByType(contactType, locale string) ([]models.Contact, error) {
	return i.lts.GetContactsByType(contactType, locale)
}

func (i *Instance) SearchProducts(locale, query string, offset, limit int) ([]models.Product, int, error) {
	return i.lts.SearchProducts(locale, query, offset, limit)
}

func (i *Instance) SearchAll(locale, query string) ([]types.SearchResult, error) {
	var results []types.SearchResult

	// Ищем продукты
	products, _, err := i.lts.SearchProducts(locale, query, 0, 10)
	if err == nil {
		for _, p := range products {
			var productName string
			for _, trans := range p.Translations {
				if trans.LanguageCode == locale {
					productName = trans.Name
					break
				}
			}
			if productName == "" && len(p.Translations) > 0 {
				productName = p.Translations[0].Name
			}

			var shortDesc string
			for _, trans := range p.Translations {
				if trans.LanguageCode == locale {
					shortDesc = trans.ShortDesc
					break
				}
			}

			results = append(results, types.SearchResult{
				Type:        "product",
				ID:          p.ID,
				Title:       productName,
				Description: shortDesc,
				URL:         fmt.Sprintf("/%s/product/%d", locale, p.ID),
				Image:       p.ImageURL,
			})
		}
	}

	// Ищем новости
	news, _, err := i.lts.SearchNews(locale, query, 0, 10)
	if err == nil {
		for _, n := range news {
			var newsTitle, excerpt string
			for _, trans := range n.Translations {
				if trans.LanguageCode == locale {
					newsTitle = trans.Title
					excerpt = trans.Excerpt
					break
				}
			}

			results = append(results, types.SearchResult{
				Type:        "news",
				ID:          n.ID,
				Title:       newsTitle,
				Description: excerpt,
				URL:         fmt.Sprintf("/%s/news/%d", locale, n.ID),
				Image:       n.ImageURL,
				Date:        n.CreatedAt.Format("02.01.2006"),
			})
		}
	}

	// Ищем страницы
	pages, err := i.lts.SearchPages(locale, query)
	if err == nil {
		for _, p := range pages {
			var pageTitle string
			for _, trans := range p.Translations {
				if trans.LanguageCode == locale {
					pageTitle = trans.Title
					break
				}
			}

			results = append(results, types.SearchResult{
				Type:        "page",
				ID:          p.ID,
				Title:       pageTitle,
				Description: "",
				URL:         fmt.Sprintf("/%s/page/%s", locale, p.Slug),
			})
		}
	}

	// Ищем документы
	docs, err := i.lts.SearchDocuments(locale, query)
	if err == nil {
		for _, d := range docs {
			var docTitle, docDesc string
			for _, trans := range d.Translations {
				if trans.LanguageCode == locale {
					docTitle = trans.Title
					docDesc = trans.Description
					break
				}
			}

			results = append(results, types.SearchResult{
				Type:        "document",
				ID:          d.ID,
				Title:       docTitle,
				Description: docDesc,
				URL:         d.FileURL,
			})
		}
	}

	return results, nil
}

func (i *Instance) SearchAPI(locale, query string, limit int) ([]types.SearchResult, error) {
	results, err := i.SearchAll(locale, query)
	if err != nil {
		return nil, err
	}

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

func (i *Instance) FilterProducts(locale string, categoryID uint, search, sortBy, sortOrder string, offset, limit int) ([]models.Product, int, error) {
	return i.lts.FilterProducts(locale, categoryID, search, sortBy, sortOrder, offset, limit)
}

func (i *Instance) GetProductsSorted(locale, search, sortBy, sortOrder string, offset, limit int) ([]models.Product, int, error) {
	return i.FilterProducts(locale, 0, search, sortBy, sortOrder, offset, limit)
}

func (i *Instance) SaveFeedback(feedback types.FeedbackRequest) (uint, error) {
	model := models.Feedback{
		Name:      feedback.Name,
		Email:     feedback.Email,
		Phone:     feedback.Phone,
		Company:   feedback.Company,
		Message:   feedback.Message,
		CreatedAt: time.Now(),
		Processed: false,
	}

	return i.lts.SaveFeedback(model)
}

func (i *Instance) GetTranslation(key, locale string) string {
	// TODO: Реализовать получение переводов из БД
	return key
}

func (i *Instance) GenerateSitemap(locale string) (string, error) {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	staticPages := []struct {
		path     string
		priority string
	}{
		{"", "1.0"},
		{"about", "0.8"},
		{"products", "0.9"},
		{"news", "0.7"},
		{"documents", "0.6"},
		{"contacts", "0.7"},
	}

	for _, page := range staticPages {
		builder.WriteString(fmt.Sprintf(`
	<url>
		<loc>%s/%s/%s</loc>
		<priority>%s</priority>
		<changefreq>weekly</changefreq>
	</url>`, i.cfg.SiteURL, locale, page.path, page.priority))
	}

	// Продукты
	products, _, err := i.lts.GetProducts(locale, 0, 1000)
	if err == nil {
		for _, product := range products {
			builder.WriteString(fmt.Sprintf(`
	<url>
		<loc>%s/%s/product/%d</loc>
		<priority>0.6</priority>
		<changefreq>monthly</changefreq>
	</url>`, i.cfg.SiteURL, locale, product.ID))
		}
	}

	// Новости
	news, _, err := i.lts.GetNews(locale, 0, 100)
	if err == nil {
		for _, item := range news {
			builder.WriteString(fmt.Sprintf(`
	<url>
		<loc>%s/%s/news/%d</loc>
		<priority>0.5</priority>
		<changefreq>yearly</changefreq>
		<lastmod>%s</lastmod>
	</url>`, i.cfg.SiteURL, locale, item.ID, item.CreatedAt.Format("2006-01-02")))
		}
	}

	builder.WriteString(`</urlset>`)
	return builder.String(), nil
}

func (i *Instance) GetAvailableLanguages() []string {
	// TODO: Получать из БД
	return []string{"ru", "en", "pl"}
}
