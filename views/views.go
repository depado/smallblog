package views

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Depado/smallblog/conf"
	"github.com/Depado/smallblog/models"
)

// Post is the views for a single post.
func Post(c *gin.Context) {
	slug := c.Param("slug")
	if val, ok := models.MPages[slug]; ok {
		c.HTML(http.StatusOK, "post.tmpl", gin.H{"post": val})
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		c.Writer.Write([]byte("404 Post not found"))
	}
}

// RawPost is used to view the raw markdown file
func RawPost(c *gin.Context) {
	slug := c.Param("slug")
	if val, ok := models.MPages[slug]; ok {
		c.Writer.Write([]byte(val.Raw))
	} else {
		c.AbortWithStatus(http.StatusNotFound)
		c.Writer.Write([]byte("404 Post not found"))
	}
}

// Index is the view to list all posts.
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"posts": models.SPages, "info": conf.C})
}
