package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	DevID  string
	APIKey string
}

// ConfigFilePath returns the path to the config file.
func ConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "vic-ptv-cli", "config.yaml")
}

// Load loads configuration from flags, environment, and config file.
// Priority: flags > env > config file.
func Load(flagDevID, flagAPIKey string) (*Config, error) {
	// 1. CLI flags (highest priority)
	if flagDevID != "" && flagAPIKey != "" {
		return &Config{DevID: flagDevID, APIKey: flagAPIKey}, nil
	}

	// 2. Environment variables
	envDevID := os.Getenv("PTV_DEV_ID")
	envAPIKey := os.Getenv("PTV_API_KEY")
	if envDevID != "" && envAPIKey != "" {
		return &Config{DevID: envDevID, APIKey: envAPIKey}, nil
	}

	// Mix flags and env if partially set
	devID := flagDevID
	if devID == "" {
		devID = envDevID
	}
	apiKey := flagAPIKey
	if apiKey == "" {
		apiKey = envAPIKey
	}

	// 3. Config file
	cfgPath := ConfigFilePath()
	if cfgPath != "" {
		viper.SetConfigFile(cfgPath)
		if err := viper.ReadInConfig(); err == nil {
			if devID == "" {
				devID = viper.GetString("devId")
			}
			if apiKey == "" {
				apiKey = viper.GetString("apiKey")
			}
		}
	}

	if devID != "" && apiKey != "" {
		return &Config{DevID: devID, APIKey: apiKey}, nil
	}

	return nil, fmt.Errorf("API credentials not configured")
}

// PrintAuthHelp prints instructions on how to configure API credentials.
func PrintAuthHelp() {
	fmt.Fprintln(os.Stderr, `PTV API credentials not configured.

To use this tool, you need a PTV Developer ID and API Key.

How to get credentials:
  1. Email PTV at APIKeyRequest@ptv.vic.gov.au
  2. Include your name and reason for requesting access
  3. You'll receive a Developer ID (numeric) and API Key

Configure credentials (pick one method):

  Option 1: Environment variables
    export PTV_DEV_ID=your_dev_id
    export PTV_API_KEY=your_api_key

  Option 2: Config file (~/.config/vic-ptv-cli/config.yaml)
    devId: your_dev_id
    apiKey: your_api_key

  Option 3: CLI flags
    ptv --dev-id your_dev_id --api-key your_api_key <command>`)
}
