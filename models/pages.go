package models

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/spf13/viper"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// MPages is the map containing all the articles. The key is the slug of the article
// for low complexity access when querying a slug.
var MPages map[string]*Page

// SPages is a sorted slice of pages. Sorted by date, it is used to render all
// the pages on the index page.
var SPages pageSlice

// Page is the main struct. It contains everything needed to render the article.
type Page struct {
	Raw         string
	Markdown    template.HTML
	Title       string
	Description string
	Author      *Author
	Banner      string
	Date        time.Time
	DateFmt     string
	Tags        []string
	File        string
	Slug        string
	Draft       bool
}

// NewPageFromFile parses a file, inserts it in the map and slice, and returns a *Page instance
func NewPageFromFile(fn string) (*Page, error) {
	var err error
	p := new(Page)

	if err = p.ParseFile(fn); err != nil {
		return nil, errors.Wrap(err, "NewPageFromFile")
	}
	if err = p.Insert(false); err != nil {
		return nil, errors.Wrap(err, "NewPageFromFile")
	}

	logrus.WithFields(logrus.Fields{
		"file": filepath.Base(fn),
		"url":  fmt.Sprintf("/post/%s", p.Slug),
	}).Info("New file added")

	return p, nil
}

type pageSlice []*Page

// Len is part of sort.Interface.
func (p pageSlice) Len() int {
	return len(p)
}

// Swap is part of sort.Interface.
func (p pageSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (p pageSlice) Less(i, j int) bool {
	return p[i].Date.After(p[j].Date)
}

// ParseFile parses a whole file and fills the Page struct.
func (p *Page) ParseFile(fn string) error {
	var err error

	p.File = filepath.Base(fn)
	h, b, err := SplitFile(fn)
	if err != nil {
		return err
	}
	p.Raw = string(b)
	p.ParseMarkdown(b)
	return p.ParseMetadata(h)
}

// UpdateFromFile parses the file once more
func (p *Page) UpdateFromFile(fn string) error {
	var err error
	old := p.Slug
	if err = p.ParseFile(fn); err != nil {
		return errors.Wrap(err, "UpdateFromFile")
	}
	if _, ok := MPages[p.Slug]; !ok {
		delete(MPages, old)
		MPages[p.Slug] = p
	}
	sort.Sort(SPages)

	logrus.WithFields(logrus.Fields{
		"file": filepath.Base(fn),
		"url":  fmt.Sprintf("/post/%s", p.Slug),
	}).Info("Modified")

	return nil
}

// Pop removes a Page (in case the file is deleted for example)
func (p *Page) Pop() {
	for i, n := range SPages {
		if n == p {
			SPages = append(SPages[:i], SPages[i+1:]...)
		}
	}
	delete(MPages, p.Slug)
	logrus.WithFields(logrus.Fields{
		"file": filepath.Base(p.File),
		"url":  fmt.Sprintf("/post/%s", p.Slug),
	}).Info("Removed")
}

// ParseMetadata parses the metadata on top of the markdown files. It will also
// raise errors when mandatory fields aren't present, or some slugs are duplicates.
func (p *Page) ParseMetadata(h []byte) error {
	var err error
	var t time.Time

	m := MetaData{}
	if err = yaml.Unmarshal(h, &m); err != nil {
		return err
	}
	if err = m.Validate(); err != nil {
		return err
	}
	if t, err = time.Parse("2006-01-02 15:04:05", m.Date); err != nil {
		return err
	}
	p.Slug = m.GenerateSlug()
	p.Description = m.Description
	p.Date = t
	p.DateFmt = t.Format("2006/01/02 15:04")
	p.Banner = m.Banner
	p.Tags = m.Tags
	p.Title = m.Title
	p.Draft = m.Draft

	// Fill in the author information if yaml header contains it. Otherwise
	// fallback to the GlobalAuthor, and if it's still empty, then nil.
	if m.Author != nil {
		p.Author = m.Author
	} else {
		p.Author = GetGlobalAuthor()
	}
	return nil
}

// ParseMarkdown will simply parse the markdown b and put it inside the Page structure.
func (p *Page) ParseMarkdown(b []byte) {
	p.Markdown = template.HTML(string(render(b)))
}

// Insert will try to insert the file in the MPages map and SPages slice. It will
// also validates that no pages have the same slug, and sort the SPages slice in case
// it's not a batch insertion. (A batch insertion means after all the inserts,
// SPages will be sorted manually)
func (p *Page) Insert(batch bool) error {
	if p.Draft && !viper.GetBool("blog.draft") {
		logrus.WithFields(logrus.Fields{"file": p.File, "state": "ignored"}).Info("Draft")
		return nil
	}
	if val, ok := MPages[p.Slug]; ok {
		return fmt.Errorf("Two pages have the same slug : %s and %s both have %s", p.File, val.File, p.Slug)
	}
	MPages[p.Slug] = p
	SPages = append(SPages, p)
	if !batch {
		sort.Sort(SPages)
	}
	return nil
}

// ParseDir cycles through a directory and parses each file one by one.
func ParseDir(dir string) error {
	var err error
	var files []os.FileInfo
	start := time.Now()

	if files, err = ioutil.ReadDir(dir); err != nil {
		return err
	}

	MPages = make(map[string]*Page, len(files))
	SPages = make(pageSlice, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		s := time.Now()
		p := Page{}
		if err = p.ParseFile(path.Join(dir, f.Name())); err != nil {
			logrus.WithFields(logrus.Fields{
				"took":  time.Since(s),
				"file":  p.File,
				"state": "ignored",
			}).WithError(err).Error("Couldn't process file")
			continue
		}
		if err = p.Insert(true); err != nil {
			logrus.WithFields(logrus.Fields{
				"took":  time.Since(s),
				"file":  p.File,
				"state": "ignored",
			}).WithError(err).Error("Couldn't insert file")
			continue
		}
		logrus.WithFields(logrus.Fields{
			"took": time.Since(s),
			"file": p.File,
			"url":  fmt.Sprintf("/post/%s", p.Slug),
		}).Debug("Page generated")
	}
	sort.Sort(SPages)
	logrus.WithFields(logrus.Fields{
		"files": len(files),
		"took":  time.Since(start),
	}).Info("Generated files")
	return nil
}
