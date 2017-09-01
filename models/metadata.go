package models

import (
	"fmt"

	"github.com/stvp/slug"
)

// meta is a struct that helps the unmarshalling of the yaml header in markdown files.
type meta struct {
	Title       string   `yaml:"title"` // Mandatory
	Description string   `yaml:"description"`
	Author      string   `yaml:"author"`
	Slug        string   `yaml:"slug"`
	Tags        []string `yaml:"tags"`
	Date        string   `yaml:"date"` // Mandatory
}

func (m *meta) Validate() error {
	if m.Date == "" {
		return fmt.Errorf("parser: The `date` field is mandatory")
	}
	if m.Title == "" {
		return fmt.Errorf("parser: The `title` field is mandatory")
	}
	return nil
}

func (m *meta) GenerateSlug() string {
	slug.Replacement = '-'
	if m.Slug == "" {
		return slug.Clean(m.Title)
	}
	return slug.Clean(m.Slug)
}
