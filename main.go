package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Depado/smallblog/conf"
	"github.com/Depado/smallblog/filesystem"
	"github.com/Depado/smallblog/models"
	"github.com/Depado/smallblog/views"
)

func main() {
	var err error

	// Configuration loading
	if err = conf.Load("conf.yml"); err != nil {
		log.Fatal(err)
	}
	// Initial posts parsing
	if err = models.ParseDir("pages"); err != nil {
		log.Fatal(err)
	}
	// Watching filesystem
	go filesystem.Watch(conf.C.PagesDir)

	// Debug mode
	if !conf.C.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	// Router initialization
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/static", "./assets")

	// Routes Definition
	r.GET("/", views.Index)
	r.GET("/post/:slug", views.Post)
	r.GET("/post/:slug/raw", views.RawPost)

	// Run
	r.Run(fmt.Sprintf("%s:%d", conf.C.Host, conf.C.Port))
}
