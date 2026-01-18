package handler

import (
	"international_site/internal/types"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) Page(c *gin.Context) {
	locale := getLocale(c)
	slug := c.Param("slug")

	page, err := s.service.GetPageBySlug(slug, locale)
	if err != nil {
		abortWithError(c, err)
		return
	}

	if page == nil {
		c.Status(404)
		return
	}

	c.JSON(200, page)
}

func (s *Server) AboutPage(c *gin.Context) {
	s.Page(c)
}

func (s *Server) CertificatesPage(c *gin.Context) {
	s.Page(c)
}

func (s *Server) PrivacyPage(c *gin.Context) {
	s.Page(c)
}

func (s *Server) ProductsPage(c *gin.Context) {
	locale := getLocale(c)
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, total, err := s.service.GetProducts(locale, offset, limit)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"items": products,
		"total": total,
	})
}

func (s *Server) ProductsByCategory(c *gin.Context) {
	locale := getLocale(c)
	slug := c.Param("slug")
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	category, err := s.service.GetCategoryBySlug(slug, locale)
	if err != nil || category == nil {
		c.Status(404)
		return
	}

	products, total, err := s.service.GetProductsByCategory(locale, category.ID, offset, limit)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"category": category,
		"items":    products,
		"total":    total,
	})
}

func (s *Server) ProductDetail(c *gin.Context) {
	locale := getLocale(c)
	id, _ := strconv.Atoi(c.Param("id"))

	product, err := s.service.GetProductByID(locale, uint(id))
	if err != nil || product == nil {
		c.Status(404)
		return
	}

	c.JSON(200, product)
}

func (s *Server) NewsPage(c *gin.Context) {
	locale := getLocale(c)
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := s.service.GetNews(locale, offset, limit)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"items": items,
		"total": total,
	})
}

func (s *Server) NewsDetail(c *gin.Context) {
	locale := getLocale(c)
	id, _ := strconv.Atoi(c.Param("id"))

	item, err := s.service.GetNewsByID(locale, uint(id))
	if err != nil || item == nil {
		c.Status(404)
		return
	}

	c.JSON(200, item)
}

func (s *Server) DocumentsPage(c *gin.Context) {
	locale := getLocale(c)
	docType := c.Param("type")

	items, err := s.service.GetDocumentsByType(docType, locale)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, items)
}

func (s *Server) ContactsPage(c *gin.Context) {
	locale := getLocale(c)
	contactType := c.Param("type")

	items, err := s.service.GetContactsByType(contactType, locale)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, items)
}

func (s *Server) SearchPage(c *gin.Context) {
	locale := getLocale(c)
	query := c.Query("q")

	results, err := s.service.SearchAll(locale, query)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, results)
}

func (s *Server) APISearch(c *gin.Context) {
	locale := getLocale(c)
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	results, err := s.service.SearchAPI(locale, query, limit)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, results)
}

func (s *Server) APIProductsFilter(c *gin.Context) {
	locale := getLocale(c)
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	search := c.Query("search")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	items, total, err := s.service.FilterProducts(
		locale,
		uint(categoryID),
		search,
		sortBy,
		sortOrder,
		offset,
		limit,
	)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"items": items,
		"total": total,
	})
}

func (s *Server) SubmitFeedback(c *gin.Context) {
	s.APISubmitFeedback(c)
}

func (s *Server) APISubmitFeedback(c *gin.Context) {
	var req types.FeedbackRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, err := s.service.SaveFeedback(req)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{"id": id})
}

func (s *Server) Sitemap(c *gin.Context) {
	locale := getLocale(c)

	xml, err := s.service.GenerateSitemap(locale)
	if err != nil {
		abortWithError(c, err)
		return
	}

	c.Data(200, "application/xml", []byte(xml))
}

func (s *Server) NotFoundPage(c *gin.Context) {
	c.Status(404)
}
