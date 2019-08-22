package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Depado/smallblog/router"
)

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start serving the blog",
	Run: func(cmd *cobra.Command, args []string) {
		router.Run()
	},
}

func init() {
	// Server
	serve.Flags().String("server.host", "127.0.0.1", "host on which the server should listen")
	serve.Flags().Int("server.port", 8080, "port on which the server should listen")
	serve.Flags().Bool("server.debug", false, "debug mode for the server")
	serve.Flags().String("server.domain", "", "domain of the blog used for RSS feed and share functionalities")
	serve.Flags().Bool("server.tls", false, "whether https is activated for the domain")

	// Gitalk
	serve.Flags().Bool("analytics.enabled", false, "enable Google Analytics feature")
	serve.Flags().String("analytics.tag", "", "tag for the Google Analytics feature")
	serve.Flags().Bool("gitalk.enabled", false, "enable the gitalk feature")
	serve.Flags().String("gitalk.client_id", "", "client ID of the gitalk app")
	serve.Flags().String("gitalk.client_secret", "", "client secret of the gitalk app")
	serve.Flags().String("gitalk.repo", "", "repository where the comments will be stored")
	serve.Flags().String("gitalk.owner", "", "repository owner")
	serve.Flags().StringArray("gitalk.admins", []string{}, "gitalk admins")
	viper.BindPFlags(serve.Flags()) // nolint: errcheck
}
