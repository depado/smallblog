package cmd

import (
	"strings"

	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddLogFlags will add the log flags to the given command
func AddLogFlags(c *cobra.Command) {
	c.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	c.PersistentFlags().String("log.format", "text", "one of text or json")
	c.PersistentFlags().Bool("log.line", false, "enable filename and line in logs")
}

// AddBlogFlags will add the blog flags to the given command
func AddBlogFlags(c *cobra.Command) {
	c.PersistentFlags().String("blog.pages", "pages/", "directory in which articles are stored")
	c.PersistentFlags().String("blog.code.style", "monokai", "style of the code sections")
	c.PersistentFlags().String("blog.author.twitter", "", "Twitter handle of the author")
	c.PersistentFlags().String("blog.author.name", "", "name (or nickname) of the author")
	c.PersistentFlags().String("blog.author.github", "", "github username of the author")
	c.PersistentFlags().String("blog.author.site", "", "website of the author")
	c.PersistentFlags().String("blog.author.avatar", "", "URL to the author's avatar")
	c.PersistentFlags().Bool("blog.share", false, "add a Twitter share button on articles")
}

// BindPersistentFlags will bind all the persistant flags to the given command
// and call BindPFlags, effectively taking those flags into account when running
func BindPersistentFlags(c *cobra.Command) error {
	AddLogFlags(c)
	AddBlogFlags(c)
	return viper.BindPFlags(c.PersistentFlags())
}

// Initialize will be run when cobra finishes its initialization
func Initialize() {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("smallblog")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Configuration file
	if viper.GetString("conf") != "" {
		viper.SetConfigFile(viper.GetString("conf"))
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config/")
	}
	hasconf := viper.ReadInConfig() == nil

	// Set log level
	lvl := viper.GetString("log.level")
	if l, err := logrus.ParseLevel(lvl); err != nil {
		logrus.WithFields(logrus.Fields{"level": lvl, "fallback": "info"}).Warn("Invalid log level")
	} else {
		logrus.SetLevel(l)
	}

	// Set log format
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      true,
		})
	}
	// Defines if logrus should display filenames and line where the log ocured
	if viper.GetBool("log.line") {
		logrus.AddHook(filename.NewHook())
	}
	// Delays the log for once the logger has been setup
	if hasconf {
		logrus.WithField("file", viper.ConfigFileUsed()).Debug("Found configuration file")
	} else {
		logrus.Debug("No configuration file found")
	}
}
