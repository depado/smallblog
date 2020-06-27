package router

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Depado/smallblog/filesystem"
	"github.com/Depado/smallblog/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Run setups and runs the server
func (br *BlogRouter) SetupRoutes() error {
	var err error

	// Initial posts parsing
	if err = models.ParseDir(br.Pages); err != nil {
		return fmt.Errorf("unable to parse dir: %w", err)
	}

	// Watching filesystem
	go filesystem.Watch(br.Pages)

	// Debug mode
	if !viper.GetBool("server.debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	// Router initialization
	br.r = gin.Default()
	br.r.LoadHTMLGlob("templates/*.tmpl")
	br.r.Static("/static", "./assets")

	// Assets for pages
	ad := filepath.Join(br.Pages, "assets")
	if _, err := os.Stat(ad); err == nil {
		br.r.Static("/assets", ad)
	}

	// Routes Definition
	br.r.GET("/", br.Index)
	br.r.GET("/drafts", br.GetDrafts)
	br.r.GET("/rss", br.GetRSSFeed)
	br.r.GET("/tag/:tag", br.PostsByTag)
	br.r.GET("/post/:slug", br.Post)
	br.r.GET("/post/:slug/raw", br.RawPost)
	br.r.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-Agent: *\nDisallow: /post/*/raw\nDisallow: /tag/*")
	})

	return nil
}
