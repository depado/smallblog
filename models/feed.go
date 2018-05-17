package models

import (
	"fmt"
	"time"

	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

// RSS is the main RSS feed
var RSS *feeds.Feed

func init() {
	now := time.Now()
	RSS = &feeds.Feed{
		Title:       viper.GetString("blog.title"),
		Description: viper.GetString("blog.description"),
		Author: &feeds.Author{
			Name: viper.GetString("blog.author"),
		},
		Link:    &feeds.Link{Href: viper.GetString("server.domain")},
		Created: now,
	}
}

// UpdateRSSFeed will update the items in the RSS feed to reflect the changes
func UpdateRSSFeed() {
	var items []*feeds.Item
	root := viper.GetString("server.domain")

	for _, v := range SPages {
		if v.Draft {
			continue
		}
		i := &feeds.Item{
			Title:       v.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%s", root, v.Slug)},
			Created:     v.Date,
			Author:      &feeds.Author{Name: v.Author.Name},
			Description: v.Description,
			Content:     v.Raw,
		}
		if v.Banner != "" {
			i.Enclosure = &feeds.Enclosure{Url: v.Banner, Type: "image"}
		}
		items = append(items, i)
	}
	RSS.Items = items
}
