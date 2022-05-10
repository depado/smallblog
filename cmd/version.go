package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build number and versions injected at compile time
var (
	Version = "unknown"
	Build   = "unknown"
	Time    = "unknown"
	Packer  = ""
)

// Descriptive help text for version command
var versionHelp = `
This command will output the build number, version number and build date of smallblog.
The build number corresponds to the sha1 commit the binary was built against,
while the version number corresponds to the latest tag the binary was built on.
Finally the build date corresponds to the date the binary was built.

If both values are "unknown" make sure to build smallblog with "make".
`

// VersionCmd is a command that will display the build number and version (if any)
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build, version and build date",
	Long:  versionHelp,
	Run: func(cmd *cobra.Command, args []string) {
		packed := ""
		if Packer != "" {
			packed = fmt.Sprintf("Packer: %s\n", Packer)
		}

		fmt.Printf("Build: %s\nVersion: %s\nBuild Date: %s\n%s", Build, Version, Time, packed)
	},
}
