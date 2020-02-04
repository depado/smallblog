package router

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"

	"github.com/Depado/smallblog/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// PostsByTag searches for posts containing tag
func (br *BlogRouter) PostsByTag(c *gin.Context) {
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
			"extra":       template.HTML(fmt.Sprintf(`Posts tagged with <span class="btn tag">%s</span>`, tag)),
			"analytics":   gin.H{"tag": viper.GetString("analytics.tag"), "enabled": viper.GetBool("analytics.enabled")},
			"author":      models.GetGlobalAuthor(),
		}
		c.HTML(http.StatusOK, "index.tmpl", data)
	} else {
		c.String(http.StatusNotFound, "404 no posts found with this tag")
	}
}

// Index is the view to list all posts.
func (br *BlogRouter) Index(c *gin.Context) {
	data := gin.H{
		"posts":       models.SPages,
		"title":       viper.GetString("blog.title"),
		"description": viper.GetString("blog.description"),
		"analytics":   gin.H{"tag": viper.GetString("analytics.tag"), "enabled": viper.GetBool("analytics.enabled")},
		"author":      models.GetGlobalAuthor(),
	}
	c.HTML(http.StatusOK, "index.tmpl", data)
}

// GetDrafts gets the unsorted drafts
func (br *BlogRouter) GetDrafts(c *gin.Context) {
	o := models.PageSlice{}
	for _, v := range models.MPages {
		if v.Draft {
			o = append(o, v)
		}
	}
	sort.Sort(o)
	data := gin.H{
		"posts":       o,
		"title":       viper.GetString("blog.title"),
		"description": viper.GetString("blog.description"),
		"extra":       template.HTML(`These articles are drafts and may be incomplete`),
		"analytics":   gin.H{"tag": viper.GetString("analytics.tag"), "enabled": viper.GetBool("analytics.enabled")},
		"author":      models.GetGlobalAuthor(),
	}
	c.HTML(http.StatusOK, "index.tmpl", data)
}
