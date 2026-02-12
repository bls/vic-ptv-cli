package cmd

import (
	"fmt"

	"github.com/bls/vic-ptv-cli/internal/api"
	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var (
	disruptionsRoute int
	disruptionsStop  int
)

var disruptionsCmd = &cobra.Command{
	Use:   "disruptions",
	Short: "Show current disruptions",
	Long:  `Show current service disruptions across the PTV network. Optionally filter by route or stop.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		routeChanged := cmd.Flags().Changed("route")
		stopChanged := cmd.Flags().Changed("stop")

		if routeChanged && stopChanged {
			return fmt.Errorf("specify either --route or --stop, not both")
		}

		var disruptions []Disruption
		if routeChanged {
			resp, err := client.DisruptionsByRoute(disruptionsRoute)
			if err != nil {
				return err
			}
			if flagJSON {
				return display.JSON(resp)
			}
			disruptions = resp.Disruptions.AllDisruptions()
		} else if stopChanged {
			resp, err := client.DisruptionsByStop(disruptionsStop)
			if err != nil {
				return err
			}
			if flagJSON {
				return display.JSON(resp)
			}
			disruptions = resp.Disruptions.AllDisruptions()
		} else {
			resp, err := client.Disruptions()
			if err != nil {
				return err
			}
			if flagJSON {
				return display.JSON(resp)
			}
			disruptions = resp.Disruptions.AllDisruptions()
		}

		display.DisruptionsList(disruptions)
		return nil
	},
}

func init() {
	disruptionsCmd.Flags().IntVar(&disruptionsRoute, "route", 0, "Filter by route ID")
	disruptionsCmd.Flags().IntVar(&disruptionsStop, "stop", 0, "Filter by stop ID")
	_ = disruptionsCmd.RegisterFlagCompletionFunc("route", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	_ = disruptionsCmd.RegisterFlagCompletionFunc("stop", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.AddCommand(disruptionsCmd)
}

// Disruption is re-exported for use in this file to avoid fully qualifying.
type Disruption = api.Disruption
