package cmd

import (
	"fmt"
	"strconv"

	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var fareCmd = &cobra.Command{
	Use:   "fare <min_zone> <max_zone>",
	Short: "Estimate fare between zones",
	Long:  `Estimate the fare for travel between two myki zones.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		minZone, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid min_zone %q: must be a number", args[0])
		}
		maxZone, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid max_zone %q: must be a number", args[1])
		}

		resp, err := client.FareEstimate(minZone, maxZone)
		if err != nil {
			return err
		}

		if flagJSON {
			return display.JSON(resp)
		}
		display.FareEstimate(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fareCmd)
}
