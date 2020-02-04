package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Depado/smallblog/models"
)

// Post is the views for a single post.
func (br *BlogRouter) Post(c *gin.Context) {
	if page, ok := models.MPages[c.Param("slug")]; ok {
		data := br.GenerateCtx()
		data["post"] = page
		data["share_url"] = page.GetShare()
		c.HTML(http.StatusOK, "post.tmpl", data)
	} else {
		c.String(http.StatusNotFound, "404 not found")
	}
}

// RawPost is used to view the raw markdown file
func (br *BlogRouter) RawPost(c *gin.Context) {
	slug := c.Param("slug")
	if val, ok := models.MPages[slug]; ok {
		if _, err := c.Writer.Write([]byte(val.Raw)); err != nil {
			logrus.WithError(err).Error("Unable to write raw markdown")
			c.String(http.StatusInternalServerError, "500 internal server error")
		}
	} else {
		c.String(http.StatusNotFound, "404 not found")
	}
}
