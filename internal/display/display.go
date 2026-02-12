package display

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/bls/vic-ptv-cli/internal/api"
)

// RouteTypeName returns a human-readable name for a route type ID.
func RouteTypeName(routeType int) string {
	switch routeType {
	case 0:
		return "Train"
	case 1:
		return "Tram"
	case 2:
		return "Bus"
	case 3:
		return "V/Line Train"
	case 4:
		return "V/Line Coach"
	default:
		return fmt.Sprintf("Unknown (%d)", routeType)
	}
}

// JSON outputs any value as indented JSON.
func JSON(v interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// SearchResults displays search results as a table.
func SearchResults(resp *api.SearchResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tNAME\tID\tROUTE TYPE")
	for _, s := range resp.Stops {
		fmt.Fprintf(w, "Stop\t%s (%s)\t%d\t%s\n", s.StopName, s.StopSuburb, s.StopID, RouteTypeName(s.RouteType))
	}
	for _, r := range resp.Routes {
		name := r.RouteName
		if r.RouteNumber != "" {
			name = r.RouteNumber + " - " + r.RouteName
		}
		fmt.Fprintf(w, "Route\t%s\t%d\t%s\n", name, r.RouteID, RouteTypeName(r.RouteType))
	}
	for _, o := range resp.Outlets {
		fmt.Fprintf(w, "Outlet\t%s (%s)\t-\t-\n", o.OutletName, o.OutletSuburb)
	}
	w.Flush()
}

// DeparturesList displays departures as a table.
func DeparturesList(resp *api.DeparturesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SCHEDULED\tESTIMATED\tROUTE\tDIRECTION\tPLATFORM")
	loc := time.Now().Location()
	for _, d := range resp.Departures {
		scheduled := "-"
		if d.ScheduledDepartureUTC != nil {
			scheduled = d.ScheduledDepartureUTC.In(loc).Format("15:04")
		}
		estimated := "-"
		if d.EstimatedDepartureUTC != nil {
			estimated = d.EstimatedDepartureUTC.In(loc).Format("15:04")
		}

		routeName := fmt.Sprintf("Route %d", d.RouteID)
		if r, ok := resp.Routes[fmt.Sprintf("%d", d.RouteID)]; ok {
			routeName = r.RouteName
			if r.RouteNumber != "" {
				routeName = r.RouteNumber + " " + r.RouteName
			}
		}

		dirName := "-"
		if dir, ok := resp.Directions[fmt.Sprintf("%d", d.DirectionID)]; ok {
			dirName = dir.DirectionName
		}

		platform := d.PlatformNumber
		if platform == "" {
			platform = "-"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", scheduled, estimated, routeName, dirName, platform)
	}
	w.Flush()
}

// StopDetail displays stop details.
func StopDetail(resp *api.StopResponse) {
	s := resp.Stop
	fmt.Printf("Stop: %s\n", s.StopName)
	fmt.Printf("ID: %d\n", s.StopID)
	fmt.Printf("Route Type: %s\n", RouteTypeName(s.RouteType))
	if s.StationType != "" {
		fmt.Printf("Station Type: %s\n", s.StationType)
	}
	if s.StationDescription != "" {
		fmt.Printf("Description: %s\n", s.StationDescription)
	}
	if s.StopAmenities != nil {
		a := s.StopAmenities
		fmt.Println("\nAmenities:")
		fmt.Printf("  Toilet: %s\n", boolYesNo(a.Toilet))
		fmt.Printf("  Taxi Rank: %s\n", boolYesNo(a.TaxiRank))
		fmt.Printf("  CCTV: %s\n", boolYesNo(a.CCTV))
		if a.CarParking != "" {
			fmt.Printf("  Car Parking: %s\n", a.CarParking)
		}
	}
	if s.StopAccessibility != nil {
		a := s.StopAccessibility
		fmt.Println("\nAccessibility:")
		fmt.Printf("  Wheelchair: %s\n", boolYesNo(a.Wheelchair))
		fmt.Printf("  Lift Access: %s\n", boolYesNo(a.LiftAccess))
		fmt.Printf("  Escalator: %s\n", boolYesNo(a.Escalator))
		fmt.Printf("  Stairs: %s\n", boolYesNo(a.Stairs))
		fmt.Printf("  Lighting: %s\n", boolYesNo(a.Lighting))
		fmt.Printf("  Hearing Loop: %s\n", boolYesNo(a.Hearing))
	}
}

func boolYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// RoutesList displays routes as a table.
func RoutesList(resp *api.RoutesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNUMBER\tNAME\tTYPE")
	for _, r := range resp.Routes {
		num := r.RouteNumber
		if num == "" {
			num = "-"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", r.RouteID, num, r.RouteName, RouteTypeName(r.RouteType))
	}
	w.Flush()
}

// RouteDetail displays route details.
func RouteDetail(resp *api.RouteResponse) {
	r := resp.Route
	fmt.Printf("Route: %s\n", r.RouteName)
	fmt.Printf("ID: %d\n", r.RouteID)
	if r.RouteNumber != "" {
		fmt.Printf("Number: %s\n", r.RouteNumber)
	}
	fmt.Printf("Type: %s\n", RouteTypeName(r.RouteType))
	if r.RouteGTFSID != "" {
		fmt.Printf("GTFS ID: %s\n", r.RouteGTFSID)
	}
	if r.RouteServiceStatus != nil {
		fmt.Printf("Service Status: %s\n", r.RouteServiceStatus.Description)
	}
}

// DisruptionsList displays disruptions as a table.
func DisruptionsList(disruptions []api.Disruption) {
	if len(disruptions) == 0 {
		fmt.Println("No current disruptions.")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tTYPE\tTITLE")
	for _, d := range disruptions {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", d.DisruptionID, d.DisruptionStatus, d.DisruptionType, truncate(d.Title, 60))
	}
	w.Flush()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// FareEstimate displays fare estimate results.
func FareEstimate(resp *api.FareEstimateResponse) {
	if resp.FareEstimate == nil {
		fmt.Println("No fare estimate available.")
		return
	}

	fe := resp.FareEstimate
	if fe.IsJourneyInFreeTramZone {
		fmt.Println("This journey is within the Free Tram Zone - no fare required!")
		return
	}
	if fe.IsEarlyBird {
		fmt.Println("Note: Early Bird fare may apply (free travel on selected trains before 7am)")
	}

	if len(fe.PassengerFares) == 0 {
		fmt.Println("No fare data available.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PASSENGER TYPE\t2 HOUR\tDAILY\tWEEKLY\tWEEKEND CAP")
	for _, f := range fe.PassengerFares {
		fmt.Fprintf(w, "%s\t$%.2f\t$%.2f\t$%.2f\t$%.2f\n",
			f.PassengerType, f.Fare2Hour, f.FareDaily, f.FareWeekly, f.FareWeekend)
	}
	w.Flush()
}

// RouteTypesList displays route types as a table.
func RouteTypesList(resp *api.RouteTypesResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME")
	for _, rt := range resp.RouteTypes {
		fmt.Fprintf(w, "%d\t%s\n", rt.RouteTypeID, rt.RouteTypeName)
	}
	w.Flush()
}
