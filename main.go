package main

import (
	"fmt"

	"github.com/Depado/smallblog/cmd"
	"github.com/Depado/smallblog/router"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Build number and versions injected at compile time
var (
	Version = "unknown"
	Build   = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "smallblog",
	Short: "Smallblog is a simple self-hosted no-db blog",
	Long:  "A simple blog engine which parses markdown files with front-matter in yaml.",
}

// Version command that will display the build number and version (if any)
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run: func(c *cobra.Command, args []string) { // nolint: unparam
		fmt.Printf("Build: %s\nVersion: %s\n", Build, Version)
	},
}

// Serve command that will actually run the server
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start serving the blog",
	Run: func(c *cobra.Command, args []string) {
		r := router.New(
			viper.GetString("blog.pages"),
			viper.GetString("server.host"),
			viper.GetInt("server.port"),
			viper.GetBool("server.debug"),
			viper.GetString("server.root_url"),
			viper.GetBool("gitalk.enabled"),
			viper.GetString("gitalk.token"),
			viper.GetString("gitalk.repo"),
			viper.GetString("gitalk.owner"),
			viper.GetStringSlice("gitalk.admins"),
			viper.GetBool("analytics.enabled"),
			viper.GetString("analytics.tag"),
			viper.GetBool("blog.share"),
		)
		if err := r.SetupRoutes(); err != nil {
			logrus.WithError(err).Fatal("Unable to setup routes")
		}
		r.Start()
	},
}

var generateCmd = &cobra.Command{
	Use:   "generate [output directory]",
	Short: "Generate a static site",
	Run: func(c *cobra.Command, args []string) {
		output := "build"
		if len(args) > 0 {
			output = args[0]
		}
		if err := cmd.RunGenerate(
			output,
			viper.GetString("blog.pages"),
			viper.GetString("blog.title"),
			viper.GetString("blog.description"),
		); err != nil {
			logrus.WithError(err).Fatal()
		}
	},
}

func main() {
	var err error

	cobra.OnInitialize(cmd.Initialize)
	if err = cmd.BindPersistentFlags(rootCmd); err != nil {
		logrus.WithError(err).Fatal()
	}
	if err = cmd.BindServeFlags(serveCmd); err != nil {
		logrus.WithError(err).Fatal()
	}

	// Adding commands to the root command
	rootCmd.AddCommand(versionCmd, serveCmd, generateCmd)

	// Adding new command to the root command
	if err = cmd.AddNewCommand(rootCmd); err != nil {
		logrus.WithError(err).Fatal()
	}
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal()
	}
}
