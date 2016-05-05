package models

import (
	"html/template"
	"reflect"
	"testing"
	"time"

	"github.com/russross/blackfriday"
)

func TestParseFile(t *testing.T) {

}

func TestParseMetadata(t *testing.T) {
	var err error

	payload := map[string]Page{
		`title: Example Article
description: Because everyone needs examples.
slug: example
author: Depado
date: 2016-05-04 15:31:00
tags:
    - inspiration`: Page{
			Title:       "Example Article",
			Description: "Because everyone needs examples.",
			Slug:        "example",
			Author:      "Depado",
			Date:        time.Date(2016, time.May, 4, 15, 31, 0, 0, time.UTC),
			DateFmt:     "2016/05/04 15:31",
			Tags:        []string{"inspiration"},
		},
	}
	mustErr := []string{
		`date: 2016-05-04 15:31:00`, // Missing title
		`title: Example Article`,    // Missing date
		`description: Because everyone needs examples.
slug: example
author: Depado`, // Missing date
		`date: xxxx-05-04 15:31:00`, // Impossible date
	}
	for in, expected := range payload {
		p := Page{}
		if err = p.ParseMetadata([]byte(in)); err != nil {
			t.Log(err)
			t.FailNow()
		}
		if !reflect.DeepEqual(p, expected) {
			t.Logf("Expected :\n%+vGot :\n%+v", expected, p)
			t.FailNow()
		}
	}
	for _, in := range mustErr {
		p := Page{}
		if err = p.ParseMetadata([]byte(in)); err == nil {
			t.Logf("%s must fail", in)
			t.FailNow()
		}
	}
}

func TestParseMarkdown(t *testing.T) {
	// Just validating that p.Markdown is set to the right value after a call to ParseMarkdown
	payload := map[string]template.HTML{
		`**Hello**`: template.HTML(string(blackfriday.MarkdownCommon([]byte(`**Hello**`)))),
	}
	for in, expected := range payload {
		p := Page{}
		p.ParseMarkdown([]byte(in))
		if p.Markdown != expected {
			t.Logf("\nExpected :\n%+v\nGot :\n%+v", expected, p.Markdown)
			t.FailNow()
		}
	}
}
