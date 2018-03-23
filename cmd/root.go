package cmd

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "smallblog",
	Short: "Smallblog is a simple self-hosted no-db blog",
	Long:  "A simple blog engine which parses markdown files with front-matter in yaml.",
}

// Execute executes the commands
func Execute() {
	rootCmd.AddCommand(serve, validate, newCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal()
	}
}

func init() {
	cobra.OnInitialize(initialize)

	// Global flags
	rootCmd.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	rootCmd.PersistentFlags().String("log.format", "text", "one of text or json")
	rootCmd.PersistentFlags().Bool("log.line", false, "enable filename and line in logs")

	// Blog
	rootCmd.PersistentFlags().String("blog.title", "", "your blog's title")
	rootCmd.PersistentFlags().String("blog.description", "", "your blog's description")
	rootCmd.PersistentFlags().String("blog.pages", "pages/", "directory in which articles are stored")
	rootCmd.PersistentFlags().String("blog.code.style", "monokai", "style of the code sections")

	// Flag binding
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initialize() {
	// Environment variables
	viper.SetEnvPrefix("smallblog")
	viper.AutomaticEnv()

	// Configuration file
	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config/")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("No configuration file found")
	}
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
