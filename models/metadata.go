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

// Author represents the author of a single article
type Author struct {
	Name    string `yaml:"name"`
	Twitter string `yaml:"twitter"`
	Site    string `yaml:"site"`
	Github  string `yaml:"github"`
	Avatar  string `yaml:"avatar"`
}

// IsEmpty checks if all the fields of an Author are blank
func (a Author) IsEmpty() bool {
	return a.Name == "" && a.Twitter == "" && a.Site == "" && a.Github == "" && a.Avatar == ""
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

// GenerateSlug will generate a slug if necessary, based on the title of the
// article. Otherwise it will return the slug set in the metadata
func (m *MetaData) GenerateSlug() string {
	slug.Replacement = '-'
	if m.Slug == "" {
		return slug.Clean(m.Title)
	}
	return slug.Clean(m.Slug)
}
