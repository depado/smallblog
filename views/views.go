package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/Depado/smallblog/models"
)

// GetRSSFeed returns the RSS feed of the blog
func GetRSSFeed(c *gin.Context) {
	rss, err := models.RSS.ToRss()
	if err != nil {
		c.String(http.StatusInternalServerError, "500 internal server error")
		return
	}
	c.Data(200, "text/xml", []byte(rss))
}

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
		data := gin.H{
			"posts":       res,
			"title":       viper.GetString("blog.title"),
			"description": viper.GetString("blog.description"),
			"extra":       template.HTML(fmt.Sprintf(`Posts tagged with <span class="home-sm-tag">%s</span>`, tag)),
			"analytics":   gin.H{"tag": viper.GetString("analytics.tag"), "enabled": viper.GetBool("analytics.enabled")},
			"author":      models.GetGlobalAuthor(),
		}
		c.HTML(http.StatusOK, "index.tmpl", data)
	} else {
		c.String(http.StatusNotFound, "404 no posts found with this tag")
	}
}

// Post is the views for a single post.
func Post(c *gin.Context) {
	if page, ok := models.MPages[c.Param("slug")]; ok {
		data := gin.H{
			"post":        page,
			"gitalk":      models.GetGitalk(),
			"extra_style": models.GlobCSS,
			"analytics":   gin.H{"tag": viper.GetString("analytics.tag"), "enabled": viper.GetBool("analytics.enabled")},
			"share":       viper.GetBool("blog.share"),
			"share_url":   page.GetShare(),
		}
		c.HTML(http.StatusOK, "post.tmpl", data)
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
	data := gin.H{
		"posts":       models.SPages,
		"title":       viper.GetString("blog.title"),
		"description": viper.GetString("blog.description"),
		"analytics":   gin.H{"tag": viper.GetString("analytics.tag"), "enabled": viper.GetBool("analytics.enabled")},
		"author":      models.GetGlobalAuthor(),
	}
	c.HTML(http.StatusOK, "index.tmpl", data)
}
