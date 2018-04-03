package models

import "github.com/spf13/viper"

var globalAuthor *Author

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

// GetGlobalAuthor retrieves the author configured in the configuration file
func GetGlobalAuthor() *Author {
	if globalAuthor == nil {
		globalAuthor = &Author{
			Name:    viper.GetString("blog.author.name"),
			Github:  viper.GetString("blog.author.github"),
			Site:    viper.GetString("blog.author.site"),
			Twitter: viper.GetString("blog.author.twitter"),
			Avatar:  viper.GetString("blog.author.avatar"),
		}
	}
	return globalAuthor
}
