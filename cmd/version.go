package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/akerl/simplefin-exporter/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of simplefin-exporter",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("%s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
