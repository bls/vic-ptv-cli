package cmd

import (
	"fmt"

	"github.com/bls/vic-ptv-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Long:  `Display the current configuration including credential status and config file location.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := config.ConfigFilePath()
		fmt.Printf("Config file: %s\n", cfgPath)

		cfg, err := config.Load(flagDevID, flagAPIKey)
		if err != nil {
			fmt.Println("Developer ID: not set")
			fmt.Println("API Key: not set")
			fmt.Println("\nStatus: not configured")
			fmt.Println("\nRun 'ptv' with no arguments for setup instructions.")
			return nil
		}

		// Mask the credentials for display
		maskedDevID := cfg.DevID
		maskedKey := "****" + cfg.APIKey[max(0, len(cfg.APIKey)-4):]

		fmt.Printf("Developer ID: %s\n", maskedDevID)
		fmt.Printf("API Key: %s\n", maskedKey)
		fmt.Println("\nStatus: configured")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
