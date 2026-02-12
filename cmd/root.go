package cmd

import (
	"fmt"
	"os"

	"github.com/bls/vic-ptv-cli/internal/api"
	"github.com/bls/vic-ptv-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagDevID  string
	flagAPIKey string
	flagJSON   bool
)

var rootCmd = &cobra.Command{
	Use:   "ptv",
	Short: "CLI for Public Transport Victoria (PTV) Timetable API",
	Long: `A command-line interface for the Public Transport Victoria (PTV) Timetable API v3.

Query stops, routes, departures, disruptions, and fares from the PTV network
covering trains, trams, buses, V/Line, and more.

Data licensed from Public Transport Victoria under Creative Commons Attribution 3.0 Australia Licence.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagDevID, "dev-id", "", "PTV Developer ID")
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "PTV API Key")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output raw JSON")
}

// newClient creates a new API client from the current config.
func newClient() (*api.Client, error) {
	cfg, err := config.Load(flagDevID, flagAPIKey)
	if err != nil {
		config.PrintAuthHelp()
		return nil, fmt.Errorf("authentication required")
	}
	return api.NewClient(cfg.DevID, cfg.APIKey), nil
}
