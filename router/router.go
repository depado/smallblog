package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Depado/smallblog/filesystem"
	"github.com/Depado/smallblog/models"
	"github.com/Depado/smallblog/views"
)

// Run setups and runs the server
func Run() {
	var err error

	// Initial posts parsing
	if err = models.ParseDir(viper.GetString("blog.pages")); err != nil {
		log.Fatal(err)
	}

	// Watching filesystem
	go filesystem.Watch(viper.GetString("blog.pages"))

	// Debug mode
	if !viper.GetBool("server.debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Router initialization
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/static", "./assets")

	// Assets for pages
	ad := filepath.Join(viper.GetString("blog.pages"), "assets")
	if _, err := os.Stat(ad); err == nil {
		r.Static("/assets", ad)
	}

	// Routes Definition
	r.GET("/", views.Index)
	r.GET("/drafts", views.GetDrafts)
	r.GET("/rss", views.GetRSSFeed)
	r.GET("/tag/:tag", views.PostsByTag)
	r.GET("/post/:slug", views.Post)
	r.GET("/post/:slug/raw", views.RawPost)
	r.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-Agent: *\nDisallow: /post/*/raw\nDisallow: /tag/*")
	})

	// Run
	logrus.WithFields(logrus.Fields{
		"host": viper.GetString("server.host"),
		"port": viper.GetInt("server.port"),
	}).Info("Starting server")

	if err = r.Run(fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))); err != nil {
		logrus.WithError(err).Fatal("Couldn't start server")
	}
}
