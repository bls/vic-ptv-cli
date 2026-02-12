package cmd

import (
	"fmt"
	"strconv"

	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var routeCmd = &cobra.Command{
	Use:   "route <route_id>",
	Short: "Show route details",
	Long:  `Show detailed information about a specific route.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		routeID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid route_id %q: must be a number", args[0])
		}

		resp, err := client.Route(routeID)
		if err != nil {
			return err
		}

		if flagJSON {
			return display.JSON(resp)
		}
		display.RouteDetail(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(routeCmd)
}
