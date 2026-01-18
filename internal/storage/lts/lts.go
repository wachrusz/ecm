package lts

import (
	"fmt"
	"international_site/internal/config"
	"international_site/internal/storage/models"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Protocol defines the methods for long-term storage.
type Protocol interface {
	GetPageBySlug(slug, locale string) (*models.Page, error)
	GetCategories(locale string) ([]models.ProductCategory, error)
	GetCategoryBySlug(slug, locale string) (*models.ProductCategory, error)
	GetCategoryByID(id uint, locale string) (*models.ProductCategory, error)
	GetProductCountByCategory(categoryID uint) (int, error)
	GetProducts(locale string, offset, limit int) ([]models.Product, int, error)
	GetProductsByCategory(locale string, categoryID uint, offset, limit int) ([]models.Product, int, error)
	GetProductByID(id uint, locale string) (*models.Product, error)
	GetRelatedProducts(locale string, categoryID, excludeID uint, limit int) ([]models.Product, error)
	SearchProducts(locale, query string, offset, limit int) ([]models.Product, int, error)
	FilterProducts(locale string, categoryID uint, search, sortBy, sortOrder string, offset, limit int) ([]models.Product, int, error)
	GetNews(locale string, offset, limit int) ([]models.News, int, error)
	GetNewsByID(id uint, locale string) (*models.News, error)
	SearchNews(locale, query string, offset, limit int) ([]models.News, int, error)
	GetDocumentsByType(docType, locale string) ([]models.Document, error)
	SearchDocuments(locale, query string) ([]models.Document, error)
	GetContactsByType(contactType, locale string) ([]models.Contact, error)
	SearchPages(locale, query string) ([]models.Page, error)
	SaveFeedback(feedback models.Feedback) (uint, error)
}

// Instance implements the LongTermStorageProtocol for Postgres.
type Instance struct {
	db *gorm.DB
}

// New creates a new PostgresStorage instance.
func New(cfg config.LTSInstance) (*Instance, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	return &Instance{db: db}, nil
}

func (i *Instance) GetPageBySlug(slug, locale string) (*models.Page, error) {
	var page models.Page

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Where("slug = ?", slug).
		First(&page).Error

	if err != nil {
		return nil, err
	}

	return &page, nil
}

func (i *Instance) GetCategories(locale string) ([]models.ProductCategory, error) {
	var categories []models.ProductCategory

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Children.Translations", "language_code = ?", locale).
		Where("parent_id IS NULL").
		Order("sort_order ASC").
		Find(&categories).Error

	return categories, err
}

func (i *Instance) GetCategoryBySlug(slug, locale string) (*models.ProductCategory, error) {
	// Ищем по slug из перевода (предполагаем, что slug хранится в поле Description)
	var categories []models.ProductCategory

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Children.Translations", "language_code = ?", locale).
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	for _, cat := range categories {
		for _, trans := range cat.Translations {
			nameSlug := strings.ToLower(strings.ReplaceAll(trans.Name, " ", "-"))
			if nameSlug == slug {
				return &cat, nil
			}
		}
	}

	return nil, fmt.Errorf("category not found")
}

func (i *Instance) GetCategoryByID(id uint, locale string) (*models.ProductCategory, error) {
	var category models.ProductCategory

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Children.Translations", "language_code = ?", locale).
		Where("id = ?", id).
		First(&category).Error

	return &category, err
}

func (i *Instance) GetProductCountByCategory(categoryID uint) (int, error) {
	var count int64

	err := i.db.Model(&models.Product{}).
		Debug().
		Where("category_id = ?", categoryID).
		Count(&count).Error

	return int(count), err
}

func (i *Instance) GetProducts(locale string, offset, limit int) ([]models.Product, int, error) {
	var products []models.Product
	var total int64

	i.db.Model(&models.Product{}).Count(&total)

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Specs.Translations", "language_code = ?", locale).
		Preload("Category.Translations", "language_code = ?", locale).
		Offset(offset).
		Limit(limit).
		Order("sort_order ASC, created_at DESC").
		Find(&products).Error

	return products, int(total), err
}

func (i *Instance) GetProductsByCategory(locale string, categoryID uint, offset, limit int) ([]models.Product, int, error) {
	var products []models.Product
	var total int64

	i.db.Model(&models.Product{}).
		Debug().
		Where("category_id = ?", categoryID).
		Count(&total)

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Specs.Translations", "language_code = ?", locale).
		Preload("Category.Translations", "language_code = ?", locale).
		Where("category_id = ?", categoryID).
		Offset(offset).
		Limit(limit).
		Order("sort_order ASC, created_at DESC").
		Find(&products).Error

	return products, int(total), err
}

func (i *Instance) GetProductByID(id uint, locale string) (*models.Product, error) {
	var product models.Product

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Specs.Translations", "language_code = ?", locale).
		Preload("Category.Translations", "language_code = ?", locale).
		Where("id = ?", id).
		First(&product).Error

	return &product, err
}

func (i *Instance) GetRelatedProducts(locale string, categoryID, excludeID uint, limit int) ([]models.Product, error) {
	var products []models.Product

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Category.Translations", "language_code = ?", locale).
		Where("category_id = ? AND id != ?", categoryID, excludeID).
		Limit(limit).
		Order("RANDOM()").
		Find(&products).Error

	return products, err
}

func (i *Instance) SearchProducts(locale, query string, offset, limit int) ([]models.Product, int, error) {
	var products []models.Product
	var total int64

	query = "%" + strings.ToLower(query) + "%"

	subQuery := i.db.Model(&models.ProductTranslation{}).
		Debug().
		Select("product_id").
		Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(short_description) LIKE ?) AND language_code = ?",
			query, query, query, locale)

	i.db.Model(&models.Product{}).
		Debug().
		Where("id IN (?)", subQuery).
		Count(&total)

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Specs.Translations", "language_code = ?", locale).
		Preload("Category.Translations", "language_code = ?", locale).
		Where("id IN (?)", subQuery).
		Offset(offset).
		Limit(limit).
		Order("sort_order ASC, created_at DESC").
		Find(&products).Error

	return products, int(total), err
}

func (i *Instance) FilterProducts(locale string, categoryID uint, search, sortBy, sortOrder string, offset, limit int) ([]models.Product, int, error) {
	var products []models.Product
	var total int64

	query := i.db.Model(&models.Product{})

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		searchQuery := "%" + strings.ToLower(search) + "%"
		subQuery := i.db.Model(&models.ProductTranslation{}).
			Debug().
			Select("product_id").
			Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?) AND language_code = ?",
				searchQuery, searchQuery, locale)
		query = query.Where("id IN (?)", subQuery)
	}

	query.Count(&total)

	switch sortBy {
	case "name":
		query = query.
			Debug().
			Joins("LEFT JOIN product_translations ON products.id = product_translations.product_id AND product_translations.language_code = ?", locale).
			Order("product_translations.name " + sortOrder)
	case "created_at":
		query = query.Order("created_at " + sortOrder)
	default:
		query = query.Order("sort_order " + sortOrder)
	}

	err := query.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Preload("Specs.Translations", "language_code = ?", locale).
		Preload("Category.Translations", "language_code = ?", locale).
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	return products, int(total), err
}

func (i *Instance) GetNews(locale string, offset, limit int) ([]models.News, int, error) {
	var news []models.News
	var total int64

	i.db.Model(&models.News{}).
		Where("published = ?", true).
		Count(&total)

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Where("published = ?", true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&news).Error

	return news, int(total), err
}

func (i *Instance) GetNewsByID(id uint, locale string) (*models.News, error) {
	var news models.News

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Where("id = ? AND published = ?", id, true).
		First(&news).Error

	return &news, err
}

func (i *Instance) SearchNews(locale, query string, offset, limit int) ([]models.News, int, error) {
	var news []models.News
	var total int64

	query = "%" + strings.ToLower(query) + "%"

	subQuery := i.db.Model(&models.NewsTranslation{}).
		Debug().
		Select("news_id").
		Where("(LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(excerpt) LIKE ?) AND language_code = ?",
			query, query, query, locale)

	i.db.Model(&models.News{}).
		Debug().
		Where("id IN (?) AND published = ?", subQuery, true).
		Count(&total)

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Where("id IN (?) AND published = ?", subQuery, true).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&news).Error

	return news, int(total), err
}

func (i *Instance) GetDocumentsByType(docType, locale string) ([]models.Document, error) {
	var documents []models.Document

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Where("type = ?", docType).
		Order("created_at DESC").
		Find(&documents).Error

	return documents, err
}

func (i *Instance) SearchDocuments(locale, query string) ([]models.Document, error) {
	var documents []models.Document

	query = "%" + strings.ToLower(query) + "%"

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Joins("LEFT JOIN document_translations ON documents.id = document_translations.document_id AND document_translations.language_code = ?", locale).
		Where("LOWER(document_translations.title) LIKE ? OR LOWER(document_translations.description) LIKE ?", query, query).
		Order("created_at DESC").
		Find(&documents).Error

	return documents, err
}

func (i *Instance) GetContactsByType(contactType, locale string) ([]models.Contact, error) {
	var contacts []models.Contact

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Where("type = ?", contactType).
		Order("sort_order ASC").
		Find(&contacts).Error

	return contacts, err
}

func (i *Instance) SearchPages(locale, query string) ([]models.Page, error) {
	var pages []models.Page

	query = "%" + strings.ToLower(query) + "%"

	err := i.db.
		Debug().
		Preload("Translations", "language_code = ?", locale).
		Joins("LEFT JOIN page_translations ON pages.id = page_translations.page_id AND page_translations.language_code = ?", locale).
		Where("LOWER(page_translations.title) LIKE ? OR LOWER(page_translations.content) LIKE ?", query, query).
		Order("created_at DESC").
		Find(&pages).Error

	return pages, err
}

func (i *Instance) SaveFeedback(feedback models.Feedback) (uint, error) {
	err := i.db.Create(&feedback).Error
	return feedback.ID, err
}
