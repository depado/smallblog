package models

import (
	"fmt"

	"github.com/stvp/slug"
)

// MetaData is a struct that helps the unmarshalling of the yaml header in markdown files.
type MetaData struct {
	Title       string   `yaml:"title"` // Mandatory
	Description string   `yaml:"description"`
	Banner      string   `yaml:"banner"`
	Author      *Author  `yaml:"author,omitempty"`
	Slug        string   `yaml:"slug"`
	Tags        []string `yaml:"tags"`
	Date        string   `yaml:"date"` // Mandatory
	Draft       bool     `yaml:"draft"`
}

// Validate validates that the metada is valid
func (m *MetaData) Validate() error {
	if m.Date == "" {
		return fmt.Errorf("parser: The `date` field is mandatory")
	}
	if m.Title == "" {
		return fmt.Errorf("parser: The `title` field is mandatory")
	}
	return nil
}

// HandleSlug will generate (and store) the slug if necessary
func (m *MetaData) HandleSlug() {
	slug.Replacement = '-'
	if m.Slug == "" {
		m.Slug = slug.Clean(m.Title)
	}
	m.Slug = slug.Clean(m.Slug)
}

// GenerateSlug will generate a slug if necessary, based on the title of the
// article. Otherwise it will return the slug set in the metadata
func (m *MetaData) GenerateSlug() string {
	slug.Replacement = '-'
	if m.Slug == "" {
		return slug.Clean(m.Title)
	}
	return slug.Clean(m.Slug)
}
