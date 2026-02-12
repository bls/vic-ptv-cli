package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var routesType string

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "List routes",
	Long:  `List all PTV routes, optionally filtered by route type.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		var routeTypes []int
		if routesType != "" {
			for _, s := range strings.Split(routesType, ",") {
				s = strings.TrimSpace(s)
				rt, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("invalid route type %q: %w", s, err)
				}
				routeTypes = append(routeTypes, rt)
			}
		}

		resp, err := client.Routes(routeTypes)
		if err != nil {
			return err
		}

		if flagJSON {
			return display.JSON(resp)
		}
		display.RoutesList(resp)
		return nil
	},
}

func init() {
	routesCmd.Flags().StringVar(&routesType, "type", "", "Filter by route type (comma-separated: 0=train,1=tram,2=bus,3=vline_train,4=vline_coach)")
	rootCmd.AddCommand(routesCmd)
}
