package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Depado/smallblog/filesystem"
	"github.com/Depado/smallblog/models"
	"github.com/Depado/smallblog/views"
)

func main() {
	var err error

	if err = models.ParseDir("pages"); err != nil {
		log.Fatal(err)
	}
	go filesystem.Watch("pages")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// Loading templates
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/static", "./assets")
	// Routes Definition
	r.GET("/", views.Index)
	r.GET("/post/:slug", views.Post)
	r.GET("/post/:slug/raw", views.RawPost)
	// Run
	r.Run(":8080")
}
