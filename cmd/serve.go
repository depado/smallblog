package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BindServeFlags will bind all the necessary flags to the given cobra
// command
func BindServeFlags(c *cobra.Command) error {
	// Server
	c.Flags().String("server.host", "127.0.0.1", "host on which the server should listen")
	c.Flags().Int("server.port", 8080, "port on which the server should listen")
	c.Flags().Bool("server.debug", false, "debug mode for the server")
	c.Flags().String("server.root_url", "", "define the root URL of the service")

	// Analytics
	c.Flags().Bool("analytics.enabled", false, "enable Google Analytics feature")
	c.Flags().String("analytics.tag", "", "tag for the Google Analytics feature")
	// Gitalk
	c.Flags().Bool("gitalk.enabled", false, "enable the gitalk feature")
	c.Flags().String("gitalk.token", "", "read-only token for the gitalk app")
	c.Flags().String("gitalk.repo", "", "repository where the comments will be stored")
	c.Flags().String("gitalk.owner", "", "repository owner")
	c.Flags().StringArray("gitalk.admins", []string{}, "gitalk admins")
	return viper.BindPFlags(c.Flags()) // nolint: errcheck
}
