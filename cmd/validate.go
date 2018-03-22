package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validate = &cobra.Command{
	Use:   "validate",
	Short: "Check that all files can be parsed",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validate")
	},
}
