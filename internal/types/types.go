package types

import (
	"international_site/internal/storage/models"
	"time"
)

type Page struct {
	FeaturedProducts []models.Product         `json:"featured_products"`
	LatestNews       []models.News            `json:"latest_news"`
	Categories       []models.ProductCategory `json:"categories"`
	Contacts         []models.Contact         `json:"contacts"`
	Locale           string                   `json:"locale"`
	Title            string                   `json:"title"`
}

// DTO для передачи данных в шаблоны
type ProductDTO struct {
	ID               uint      `json:"id"`
	SKU              string    `json:"sku"`
	Name             string    `json:"name"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	CategoryID       uint      `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	ImageURL         string    `json:"image_url"`
	FileURL          string    `json:"file_url"`
	Specs            []SpecDTO `json:"specs"`
	CreatedAt        time.Time `json:"created_at"`
}

type CategoryDTO struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Slug         string `json:"slug"`
	ParentID     *uint  `json:"parent_id"`
	ProductCount int    `json:"product_count"`
}

type NewsDTO struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Excerpt   string    `json:"excerpt"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type DocumentDTO struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileURL     string    `json:"file_url"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

type ContactDTO struct {
	ID    uint   `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type FeedbackRequest struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Company string `json:"company"`
	Message string `json:"message" binding:"required"`
}

type SearchResult struct {
	Type        string `json:"type"`
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Image       string `json:"image,omitempty"`
	Date        string `json:"date,omitempty"`
}

type SpecDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
