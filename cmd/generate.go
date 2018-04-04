package cmd

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Depado/smallblog/models"
)

var generate = &cobra.Command{
	Use:   "generate [output directory]",
	Short: "Generate a static site",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		output := "build"
		if len(args) > 0 {
			output = args[0]
		}
		createStructure(output)
		pages := filepath.Join(output, "post")
		if err = models.ParseDir(viper.GetString("blog.pages")); err != nil {
			logrus.WithError(err).Fatal("Couldn't parse directory")
		}
		t, err := template.ParseGlob("templates/*.tmpl")
		for _, v := range models.MPages {
			var fd *os.File
			v.Slug = v.Slug + ".html"
			if fd, err = os.Create(filepath.Join(pages, v.Slug)); err != nil {
				logrus.WithError(err).Fatal("Couldn't create file")
			}
			t.ExecuteTemplate(fd, "post.tmpl", gin.H{"post": v, "gitalk": models.GetGitalk(), "local": true})
		}
		var fd *os.File
		if fd, err = os.Create(filepath.Join(output, "index.html")); err != nil {
			logrus.WithError(err).Fatal("Couldn't create file")
		}
		data := gin.H{
			"posts":       models.SPages,
			"title":       viper.GetString("blog.title"),
			"description": viper.GetString("blog.description"),
			"local":       true,
			"author":      models.GetGlobalAuthor(),
		}
		if err = t.ExecuteTemplate(fd, "index.tmpl", data); err != nil {
			logrus.WithError(err).Fatal("Couldn't create index")
		}
	},
}

func createStructure(output string) {
	var err error

	pages := filepath.Join(output, "post")
	if _, err = os.Stat(output); os.IsNotExist(err) {
		if err = os.Mkdir(output, os.ModePerm); err != nil {
			logrus.WithError(err).Fatal("Couldn't create output directory")
		}
	}
	if _, err = os.Stat(pages); os.IsNotExist(err) {
		if err = os.Mkdir(pages, os.ModePerm); err != nil {
			logrus.WithError(err).Fatal("Couldn't create output directory")
		}
	}
	if err = copy.Copy("assets", filepath.Join(output, "/static")); err != nil {
		logrus.WithError(err).Fatal("Couldn't copy assets")
	}
	if err = copy.Copy("assets", filepath.Join(output, "/static")); err != nil {
		logrus.WithError(err).Fatal("Couldn't copy assets")
	}
	ad := filepath.Join(viper.GetString("blog.pages"), "assets")
	if _, err := os.Stat(ad); err == nil {
		if err := copy.Copy(ad, filepath.Join(output, "/assets")); err != nil {
			logrus.WithError(err).Fatal("Couldn't copy pages assets")
		}
	}
}
