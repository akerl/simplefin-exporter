package cmd

import (
	"fmt"

	"github.com/akerl/simplefin-exporter/config"
	"github.com/akerl/simplefin-exporter/fetcher"

	"github.com/akerl/metrics/server"
	"github.com/spf13/cobra"
)

func serveRunner(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("no config file provided")
	}
	configPath := args[0]

	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	cache := server.Cache{}

	f := fetcher.NewFetcher(conf, &cache)
	s := server.NewServer(conf.Port, &cache)

	f.RunAsync()
	return s.Run()
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run web server to serve metrics",
	RunE:  serveRunner,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
