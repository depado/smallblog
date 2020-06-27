package router

import (
	"net/http"

	"github.com/Depado/smallblog/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetRSSFeed returns the RSS feed of the blog
func (br *BlogRouter) GetRSSFeed(c *gin.Context) {
	rss, err := models.RSS.ToRss()
	if err != nil {
		logrus.WithError(err).Error("Unable to generate RSS feed")
		c.String(http.StatusInternalServerError, "500 internal server error")
		return
	}

	c.Data(200, "text/xml", []byte(rss))
}
