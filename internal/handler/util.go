package handler

import "github.com/gin-gonic/gin"

func getLocale(c *gin.Context) string {
	locale := c.Param("locale")
	if locale == "" {
		locale = "en"
	}
	return locale
}

func abortWithError(c *gin.Context, err error) {
	c.JSON(500, gin.H{"error": err.Error()})
}
