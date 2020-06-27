package cmd

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"

	"github.com/Depado/smallblog/models"
)

// RunGenerate will run the generate command and generate a static build of the
// blog
func RunGenerate(output, input, title, description string) error {
	var err error
	var t *template.Template
	var fd *os.File

	pages := filepath.Join(output, "post")
	if err = createStructure(input, output); err != nil {
		return errors.Wrap(err, "create directory structure")
	}
	if err = models.ParseDir(input); err != nil {
		return errors.Wrap(err, "parse input directory")
	}
	if t, err = template.ParseGlob("templates/*.tmpl"); err != nil {
		return errors.Wrap(err, "parse template glob")
	}
	for _, v := range models.MPages {
		v.Slug = v.Slug + ".html"
		if fd, err = os.Create(filepath.Join(pages, v.Slug)); err != nil {
			return errors.Wrap(err, "create post file")
		}
		if err = t.ExecuteTemplate(
			fd, "post.tmpl",
			gin.H{"post": v, "gitalk": gin.H{}, "local": true},
		); err != nil {
			return errors.Wrap(err, "execute template post")
		}
	}

	if fd, err = os.Create(filepath.Join(output, "index.html")); err != nil {
		return errors.Wrap(err, "create index file")
	}
	data := gin.H{
		"posts":       models.SPages,
		"title":       title,
		"description": description,
		"local":       true,
		"author":      models.GetGlobalAuthor(),
	}
	if err = t.ExecuteTemplate(fd, "index.tmpl", data); err != nil {
		return errors.Wrap(err, "execute template index")
	}
	return nil
}

func createStructure(input, output string) error {
	var err error

	ouputp := filepath.Join(output, "post")
	if _, err = os.Stat(output); os.IsNotExist(err) {
		if err = os.Mkdir(output, os.ModePerm); err != nil {
			return errors.Wrap(err, "create output directory")
		}
	}
	if _, err = os.Stat(ouputp); os.IsNotExist(err) {
		if err = os.Mkdir(ouputp, os.ModePerm); err != nil {
			return errors.Wrap(err, "create pages directory")
		}
	}
	if err = copy.Copy("assets", filepath.Join(output, "/static")); err != nil {
		return errors.Wrap(err, "copy assets")
	}
	ad := filepath.Join(input, "assets")
	if _, err = os.Stat(ad); err == nil {
		if err = copy.Copy(ad, filepath.Join(output, "/assets")); err != nil {
			return errors.Wrap(err, "copy page assets")
		}
	}
	return nil
}
