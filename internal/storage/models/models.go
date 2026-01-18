package models

import (
	"time"
)

// Language - языки сайта
type Language struct {
	Code string `json:"code" gorm:"column:code;primaryKey;size:10"`
	Name string `json:"name" gorm:"column:name;size:100"`
}

// Page - страница сайта
type Page struct {
	ID           uint              `json:"id" gorm:"column:page_id;primaryKey;autoIncrement"`
	Slug         string            `json:"slug" gorm:"column:slug;uniqueIndex;size:255"`
	Template     string            `json:"template" gorm:"column:template;size:100"`
	CreatedAt    time.Time         `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Translations []PageTranslation `json:"translations" gorm:"foreignKey:PageID;references:ID"`
}

// PageTranslation - переводы страниц
type PageTranslation struct {
	PageID       uint   `json:"page_id" gorm:"column:page_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Title        string `json:"title" gorm:"column:title;size:255"`
	Content      string `json:"content" gorm:"column:content;type:text"`
	MetaTitle    string `json:"meta_title" gorm:"column:meta_title;size:255"`
	MetaDesc     string `json:"meta_description" gorm:"column:meta_description;size:500"`
}

// ProductCategory - категории продукции
type ProductCategory struct {
	ID           uint                         `json:"id" gorm:"column:category_id;primaryKey;autoIncrement"`
	ParentID     *uint                        `json:"parent_id" gorm:"column:parent_id;index"`
	SortOrder    int                          `json:"sort_order" gorm:"column:sort_order;default:0"`
	CreatedAt    time.Time                    `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Translations []ProductCategoryTranslation `json:"translations" gorm:"foreignKey:CategoryID;references:ID"`
	Children     []ProductCategory            `json:"children" gorm:"foreignKey:ParentID;references:ID"`
	Products     []Product                    `json:"products" gorm:"foreignKey:CategoryID;references:ID"`
}

// ProductCategoryTranslation
type ProductCategoryTranslation struct {
	CategoryID   uint   `json:"category_id" gorm:"column:category_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Name         string `json:"name" gorm:"column:name;size:255"`
	Description  string `json:"description" gorm:"column:description;type:text"`
}

// Product - продукт
type Product struct {
	ID           uint                 `json:"id" gorm:"column:product_id;primaryKey;autoIncrement"`
	CategoryID   uint                 `json:"category_id" gorm:"column:category_id;index"`
	SKU          string               `json:"sku" gorm:"column:sku;uniqueIndex;size:100"`
	ImageURL     string               `json:"image_url" gorm:"column:image_url;size:500"`
	FileURL      string               `json:"file_url" gorm:"column:file_url;size:500"`
	SortOrder    int                  `json:"sort_order" gorm:"column:sort_order;default:0"`
	CreatedAt    time.Time            `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time            `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Category     ProductCategory      `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	Translations []ProductTranslation `json:"translations" gorm:"foreignKey:ProductID;references:ID"`
	Specs        []ProductSpec        `json:"specs" gorm:"foreignKey:ProductID;references:ID"`
}

// ProductTranslation
type ProductTranslation struct {
	ProductID    uint   `json:"product_id" gorm:"column:product_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Name         string `json:"name" gorm:"column:name;size:255"`
	Description  string `json:"description" gorm:"column:description;type:text"`
	ShortDesc    string `json:"short_description" gorm:"column:short_description;size:500"`
}

// ProductSpec - характеристики продукта
type ProductSpec struct {
	ID           uint                     `json:"id" gorm:"column:product_spec_id;primaryKey;autoIncrement"`
	ProductID    uint                     `json:"product_id" gorm:"column:product_id;index"`
	SortOrder    int                      `json:"sort_order" gorm:"column:sort_order;default:0"`
	CreatedAt    time.Time                `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Product      Product                  `json:"product" gorm:"foreignKey:ProductID;references:ID"`
	Translations []ProductSpecTranslation `json:"translations" gorm:"foreignKey:SpecID;references:ID"`
}

// ProductSpecTranslation
type ProductSpecTranslation struct {
	SpecID       uint   `json:"spec_id" gorm:"column:spec_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Name         string `json:"name" gorm:"column:name;size:255"`
	Value        string `json:"value" gorm:"column:value;size:500"`
}

// News - новости
type News struct {
	ID           uint              `json:"id" gorm:"column:news_id;primaryKey;autoIncrement"`
	ImageURL     string            `json:"image_url" gorm:"column:image_url;size:500"`
	Published    bool              `json:"published" gorm:"column:published;default:false"`
	CreatedAt    time.Time         `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time         `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	Translations []NewsTranslation `json:"translations" gorm:"foreignKey:NewsID;references:ID"`
}

// NewsTranslation
type NewsTranslation struct {
	NewsID       uint   `json:"news_id" gorm:"column:news_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Title        string `json:"title" gorm:"column:title;size:255"`
	Content      string `json:"content" gorm:"column:content;type:text"`
	Excerpt      string `json:"excerpt" gorm:"column:excerpt;size:500"`
}

// Document - документы (ГОСТы, сертификаты)
type Document struct {
	ID           uint                  `json:"id" gorm:"column:document_id;primaryKey;autoIncrement"`
	FileURL      string                `json:"file_url" gorm:"column:file_url;size:500"`
	Type         string                `json:"type" gorm:"column:type;size:50"` // gost, certificate, reference
	CreatedAt    time.Time             `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Translations []DocumentTranslation `json:"translations" gorm:"foreignKey:DocumentID;references:ID"`
}

// DocumentTranslation
type DocumentTranslation struct {
	DocumentID   uint   `json:"document_id" gorm:"column:document_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Title        string `json:"title" gorm:"column:title;size:255"`
	Description  string `json:"description" gorm:"column:description;type:text"`
}

// Contact - контактная информация
type Contact struct {
	ID           uint                 `json:"id" gorm:"column:contact_id;primaryKey;autoIncrement"`
	Type         string               `json:"type" gorm:"column:type;size:50"` // phone, email, address, map
	Value        string               `json:"value" gorm:"column:value;size:500"`
	SortOrder    int                  `json:"sort_order" gorm:"column:sort_order;default:0"`
	CreatedAt    time.Time            `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Translations []ContactTranslation `json:"translations" gorm:"foreignKey:ContactID;references:ID"`
}

// ContactTranslation
type ContactTranslation struct {
	ContactID    uint   `json:"contact_id" gorm:"column:contact_id;primaryKey"`
	LanguageCode string `json:"language_code" gorm:"column:language_code;primaryKey;size:10"`
	Label        string `json:"label" gorm:"column:label;size:255"` // например "Главный офис"
}

// Feedback - форма обратной связи
type Feedback struct {
	ID        uint      `json:"id" gorm:"column:feedback_id;primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"column:name;size:100"`
	Email     string    `json:"email" gorm:"column:email;size:100"`
	Phone     string    `json:"phone" gorm:"column:phone;size:20"`
	Company   string    `json:"company" gorm:"column:company;size:200"`
	Message   string    `json:"message" gorm:"column:message;type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	Processed bool      `json:"processed" gorm:"column:processed;default:false"`
}
