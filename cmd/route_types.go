package cmd

import (
	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var routeTypesCmd = &cobra.Command{
	Use:   "route-types",
	Short: "List route types and their IDs",
	Long:  `List all PTV route types (train, tram, bus, etc.) and their numeric IDs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		resp, err := client.RouteTypes()
		if err != nil {
			return err
		}

		if flagJSON {
			return display.JSON(resp)
		}
		display.RouteTypesList(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(routeTypesCmd)
}
