package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/Depado/smallblog/conf"
	"github.com/Depado/smallblog/filesystem"
	"github.com/Depado/smallblog/models"
	"github.com/Depado/smallblog/views"
)

func main() {
	var err error

	// Configuration loading
	conf.Load()

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

	// Routes Definition
	r.GET("/", views.Index)
	r.GET("/tag/:tag", views.PostsByTag)
	r.GET("/post/:slug", views.Post)
	r.GET("/post/:slug/raw", views.RawPost)
	r.GET("/robots.txt", func(c *gin.Context) { c.String(http.StatusOK, "User-Agent: *\nDisallow: /post/*/raw") })

	// Run
	logrus.WithFields(logrus.Fields{
		"host": viper.GetString("server.host"),
		"port": viper.GetInt("server.port"),
	}).Info("Starting server")
	r.Run(fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port")))
}
