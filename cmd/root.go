package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "simplefin-exporter",
	Short:         "Serve metrics for simplefin data",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute function is the entrypoint for the CLI
func Execute() error {
	return rootCmd.Execute()
}
