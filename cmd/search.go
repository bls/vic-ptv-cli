package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bls/vic-ptv-cli/internal/display"
	"github.com/spf13/cobra"
)

var searchRouteTypes string

var searchCmd = &cobra.Command{
	Use:   "search <term>",
	Short: "Search for stops, routes, and outlets",
	Long:  `Search the PTV network for stops, routes, and outlets matching a search term.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		term := strings.Join(args, " ")

		var routeTypes []int
		if searchRouteTypes != "" {
			for _, s := range strings.Split(searchRouteTypes, ",") {
				s = strings.TrimSpace(s)
				rt, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("invalid route type %q: %w", s, err)
				}
				routeTypes = append(routeTypes, rt)
			}
		}

		resp, err := client.Search(term, routeTypes)
		if err != nil {
			return err
		}

		if flagJSON {
			return display.JSON(resp)
		}
		display.SearchResults(resp)
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchRouteTypes, "route-types", "", "Filter by route types (comma-separated: 0=train,1=tram,2=bus,3=vline_train,4=vline_coach)")
	rootCmd.AddCommand(searchCmd)
}
