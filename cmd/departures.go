package cmd

import (
	"fmt"
	"strconv"

	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var (
	departuresRouteType int
	departuresLimit     int
)

var departuresCmd = &cobra.Command{
	Use:   "departures <stop_id>",
	Short: "Show upcoming departures from a stop",
	Long: `Show upcoming departures from a stop. Requires a stop ID and route type.

Route types: 0=Train, 1=Tram, 2=Bus, 3=V/Line Train, 4=V/Line Coach`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		stopID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid stop_id %q: must be a number", args[0])
		}

		if !cmd.Flags().Changed("route-type") {
			return fmt.Errorf("--route-type is required (0=train, 1=tram, 2=bus, 3=vline_train, 4=vline_coach)")
		}

		resp, err := client.Departures(departuresRouteType, stopID, departuresLimit)
		if err != nil {
			return err
		}

		if flagJSON {
			return display.JSON(resp)
		}
		display.DeparturesList(resp)
		return nil
	},
}

func init() {
	departuresCmd.Flags().IntVar(&departuresRouteType, "route-type", -1, "Route type (0=train, 1=tram, 2=bus, 3=vline_train, 4=vline_coach)")
	departuresCmd.Flags().IntVar(&departuresLimit, "limit", 5, "Maximum number of departures to show")
	rootCmd.AddCommand(departuresCmd)
}
