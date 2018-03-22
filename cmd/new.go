package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new article in the pages directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new")
	},
}
