package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/Depado/smallblog/models"
)

// PostsByTag searches for posts containing tag
func PostsByTag(c *gin.Context) {
	tag := c.Param("tag")
	res := []*models.Page{}
	for _, v := range models.SPages {
		for _, b := range v.Tags {
			if b == tag {
				res = append(res, v)
			}
		}
	}
	if len(res) > 0 {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"posts":       res,
			"title":       viper.GetString("blog.title"),
			"description": viper.GetString("blog.description"),
			"extra":       template.HTML(fmt.Sprintf(`Posts tagged with <span class="home-sm-tag">%s</span>`, tag)),
		})
	} else {
		c.String(http.StatusNotFound, "404 no posts found with this tag")
	}
}

// Post is the views for a single post.
func Post(c *gin.Context) {
	slug := c.Param("slug")
	if val, ok := models.MPages[slug]; ok {
		c.HTML(http.StatusOK, "post.tmpl", gin.H{"post": val, "gitalk": models.GetGitalk()})
	} else {
		c.String(http.StatusNotFound, "404 not found")
	}
}

// RawPost is used to view the raw markdown file
func RawPost(c *gin.Context) {
	slug := c.Param("slug")
	if val, ok := models.MPages[slug]; ok {
		c.Writer.Write([]byte(val.Raw))
	} else {
		c.String(http.StatusNotFound, "404 not found")
	}
}

// Index is the view to list all posts.
func Index(c *gin.Context) {
	a := models.GetGlobalAuthor()
	data := gin.H{
		"posts":       models.SPages,
		"title":       viper.GetString("blog.title"),
		"description": viper.GetString("blog.description"),
	}
	if !a.IsEmpty() {
		data["author"] = a
	}
	c.HTML(http.StatusOK, "index.tmpl", data)
}
