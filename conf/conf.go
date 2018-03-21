package conf

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Load is in charge of correctly loading the configuration from flags, env
// configuration file and default values
func Load() {
	// Defining the flags as the primary and most important source

	// Logger
	pflag.String("log.level", "info", "one of debug, info, warn, error or fatal")
	pflag.String("log.format", "text", "one of text or json")
	pflag.Bool("log.line", false, "enable filename and line in logs")

	// Server
	pflag.String("server.host", "127.0.0.1", "host on which the server should listen")
	pflag.Int("server.port", 8080, "port on which the server should listen")
	pflag.Bool("server.debug", false, "debug mode for the server")

	// Gitalk
	pflag.Bool("gitalk.enabled", false, "enable the gitalk feature")
	pflag.String("gitalk.client_id", "", "client ID of the gitalk app")
	pflag.String("gitalk.client_secret", "", "client secret of the gitalk app")
	pflag.String("gitalk.repo", "", "repository where the comments will be stored")
	pflag.String("gitalk.owner", "", "repository owner")
	pflag.StringArray("gitalk.admins", []string{}, "gitalk admins")

	// Blog
	pflag.String("blog.title", "", "your blog's title")
	pflag.String("blog.description", "", "your blog's description")
	pflag.String("blog.pages", "pages/", "directory in which articles are stored")
	pflag.String("blog.code.style", "monokai", "style of the code sections")

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		logrus.WithError(err).Fatal("Couldn't bind flags")
	}

	// Environment
	viper.SetEnvPrefix("smallblog")
	viper.AutomaticEnv()

	// Configuration file
	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config/")
	viper.ReadInConfig()

	// Defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("server.host", "127.0.0.1")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("blog.pages", "pages/")
	viper.SetDefault("blog.code.style", "monokai")

	// Parsing flags
	pflag.Parse()
	setupLogger()
}

func setupLogger() {
	lvl := viper.GetString("log.level")
	l, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.WithField("level", lvl).Warn("Invalid log level, fallback to 'info'")
	} else {
		logrus.SetLevel(l)
	}
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if viper.GetBool("log.line") {
		logrus.AddHook(filename.NewHook())
	}
}

// GetGitalk returns a gitalk struct
func GetGitalk() Gitalk {
	return Gitalk{
		Enabled:      viper.GetBool("gitalk.enabled"),
		ClientID:     viper.GetString("gitalk.client_id"),
		ClientSecret: viper.GetString("gitalk.client_secret"),
		Repo:         viper.GetString("gitalk.repo"),
		Owner:        viper.GetString("gitalk.owner"),
		Admin:        viper.GetStringSlice("gitalk.admins"),
	}
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
