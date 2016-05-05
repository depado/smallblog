package models

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"

	"github.com/russross/blackfriday"
	"github.com/stvp/slug"
	"gopkg.in/yaml.v2"
)

// MPages is the map containing all the articles. The key is the slug of the article
// for low complexity access when querying a slug.
var MPages map[string]*Page

// SPages is a sorted slice of pages. Sorted by date, it is used to render all
// the pages on the index page.
var SPages pageSlice

// meta is a struct that helps the unmarshalling of the yaml header in markdown files.
type meta struct {
	Title       string   `yaml:"title"` // Mandatory
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Slug        string   `yaml:"slug"`
	Tags        []string `yaml:"tags"`
	Date        string   `yaml:"date"` // Mandatory
}

// Page is the main struct. It contains everything needed to render the article.
type Page struct {
	Raw         string
	Markdown    template.HTML
	Title       string
	Description string
	Author      string
	Date        time.Time
	DateFmt     string
	Tags        []string
	File        string
	Slug        string
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
	var file *os.File

	p.File = filepath.Base(fn)
	if file, err = os.Open(fn); err != nil {
		return err
	}
	defer file.Close()
	header := ""
	body := ""
	in := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == "" && in {
			in = false
			continue
		}
		if in {
			header += scanner.Text() + "\n"
		} else {
			body += scanner.Text() + "\n"
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	p.Raw = body
	p.ParseMarkdown([]byte(body))
	if err = p.ParseMetadata([]byte(header)); err != nil {
		return err
	}
	return nil
}

// ParseMetadata parses the metadata on top of the markdown files. It will also
// raise errors when mandatory fields aren't present, or some slugs are duplicates.
func (p *Page) ParseMetadata(h []byte) error {
	var err error
	var t time.Time

	m := meta{}
	if err = yaml.Unmarshal(h, &m); err != nil {
		return err
	}
	if m.Date == "" {
		return fmt.Errorf("Parser: The `date` field is mandatory.")
	}
	if m.Title == "" {
		return fmt.Errorf("Parser: The `title` field is mandatory")
	}
	slug.Replacement = '-'
	if m.Slug == "" {
		p.Slug = slug.Clean(m.Title)
	} else {
		p.Slug = slug.Clean(m.Slug)
	}
	if t, err = time.Parse("2006-01-02 15:04:05", m.Date); err != nil {
		return err
	}
	p.Description = m.Description
	p.Date = t
	p.DateFmt = t.Format("2006/01/02 15:04")
	p.Tags = m.Tags
	p.Author = m.Author
	p.Title = m.Title
	return nil
}

// ParseMarkdown will simply parse the markdown b and put it inside the Page structure.
func (p *Page) ParseMarkdown(b []byte) {
	p.Markdown = template.HTML(string(blackfriday.MarkdownCommon(b)))
}

func (p *Page) Insert(batch bool) error {
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
		s := time.Now()
		p := Page{}
		if err = p.ParseFile(path.Join(dir, f.Name())); err != nil {
			fmt.Printf("[SB] [ ERR] [% 9v] [%s] Could not process file : %s\n", time.Since(s), p.File, err)
			continue
		}
		if err = p.Insert(true); err != nil {
			fmt.Printf("[SB] [ ERR] [% 9v] [%s] [Ignored] Could not insert file : %s\n", time.Since(s), p.File, err)
			continue
		}
		fmt.Printf("[SB] [INFO] [% 9v] [%s] [/post/%s] %s\n", time.Since(s), p.File, p.Slug, p.Title)
	}
	sort.Sort(SPages)
	fmt.Printf("[SB] [INFO] Generated %d pages in %v\n", len(files), time.Since(start))
	return nil
}
