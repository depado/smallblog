package models

import (
	"github.com/spf13/viper"
)

var gitalk *Gitalk

// GetGitalk returns a gitalk struct
func GetGitalk() *Gitalk {
	if gitalk == nil {
		gitalk = &Gitalk{
			Enabled:      viper.GetBool("gitalk.enabled"),
			ClientID:     viper.GetString("gitalk.client_id"),
			ClientSecret: viper.GetString("gitalk.client_secret"),
			Repo:         viper.GetString("gitalk.repo"),
			Owner:        viper.GetString("gitalk.owner"),
			Admin:        viper.GetStringSlice("gitalk.admins"),
		}
	}
	return gitalk
}

// Gitalk is a struct holding all the information necessary to make gitalk work
type Gitalk struct {
	Enabled      bool     `yaml:"enabled"`
	ClientID     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Repo         string   `yaml:"repo"`
	Owner        string   `yaml:"owner"`
	Admin        []string `yaml:"admin"`
}
