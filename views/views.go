package views

import (
	"net/http"

	"github.com/Depado/smallblog/models"
	"github.com/gin-gonic/gin"
)

// Post is the views for a single post.
func Post(c *gin.Context) {
	link := c.Param("link")
	if val, ok := models.MPages[link]; ok {
		c.HTML(http.StatusOK, "post.tmpl", gin.H{"post": val})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		c.Writer.Write([]byte("404 Article not found"))
	}
}

// Index is the view to list all posts.
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"posts": models.SPages})
}
