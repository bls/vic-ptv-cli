package api

import "time"

// Status is the common status field in API responses.
type Status struct {
	Version string `json:"version"`
	Health  int    `json:"health"`
}

// SearchResponse is the response from GET /v3/search/{search_term}.
type SearchResponse struct {
	Stops   []ResultStop   `json:"stops"`
	Routes  []ResultRoute  `json:"routes"`
	Outlets []ResultOutlet `json:"outlets"`
	Status  Status         `json:"status"`
}

// ResultStop is a stop result from search.
type ResultStop struct {
	StopID       int     `json:"stop_id"`
	StopName     string  `json:"stop_name"`
	StopSuburb   string  `json:"stop_suburb"`
	RouteType    int     `json:"route_type"`
	StopLatitude float64 `json:"stop_latitude"`
	StopLongitude float64 `json:"stop_longitude"`
}

// ResultRoute is a route result from search.
type ResultRoute struct {
	RouteID     int    `json:"route_id"`
	RouteName   string `json:"route_name"`
	RouteNumber string `json:"route_number"`
	RouteType   int    `json:"route_type"`
	RouteGTFSID string `json:"route_gtfs_id"`
}

// ResultOutlet is an outlet result from search.
type ResultOutlet struct {
	OutletName     string  `json:"outlet_name"`
	OutletBusiness string  `json:"outlet_business"`
	OutletSuburb   string  `json:"outlet_suburb"`
	OutletLatitude float64 `json:"outlet_latitude"`
	OutletLongitude float64 `json:"outlet_longitude"`
}

// DeparturesResponse is the response from GET /v3/departures/...
type DeparturesResponse struct {
	Departures []Departure          `json:"departures"`
	Stops      map[string]StopInfo  `json:"stops"`
	Routes     map[string]RouteInfo `json:"routes"`
	Runs       map[string]RunInfo   `json:"runs"`
	Directions map[string]Direction `json:"directions"`
	Status     Status               `json:"status"`
}

// Departure is a single departure.
type Departure struct {
	StopID                int        `json:"stop_id"`
	RouteID               int        `json:"route_id"`
	RunID                 int        `json:"run_id"`
	RunRef                string     `json:"run_ref"`
	DirectionID           int        `json:"direction_id"`
	DisruptionIDs         []int      `json:"disruption_ids"`
	ScheduledDepartureUTC *time.Time `json:"scheduled_departure_utc"`
	EstimatedDepartureUTC *time.Time `json:"estimated_departure_utc"`
	AtPlatform            bool       `json:"at_platform"`
	PlatformNumber        string     `json:"platform_number"`
	Flags                 string     `json:"flags"`
	DepartureSequence     int        `json:"departure_sequence"`
}

// StopInfo is expanded stop info in departures response.
type StopInfo struct {
	StopID       int    `json:"stop_id"`
	StopName     string `json:"stop_name"`
	StopSuburb   string `json:"stop_suburb"`
	RouteType    int    `json:"route_type"`
}

// RouteInfo is expanded route info in departures response.
type RouteInfo struct {
	RouteID     int    `json:"route_id"`
	RouteName   string `json:"route_name"`
	RouteNumber string `json:"route_number"`
	RouteType   int    `json:"route_type"`
}

// RunInfo is expanded run info in departures response.
type RunInfo struct {
	RunID       int    `json:"run_id"`
	RunRef      string `json:"run_ref"`
	RouteID     int    `json:"route_id"`
	RouteType   int    `json:"route_type"`
	DirectionID int    `json:"direction_id"`
}

// Direction is a direction entry.
type Direction struct {
	DirectionID   int    `json:"direction_id"`
	DirectionName string `json:"direction_name"`
	RouteID       int    `json:"route_id"`
	RouteType     int    `json:"route_type"`
}

// StopResponse is the response from GET /v3/stops/{stop_id}/route_type/{route_type}.
type StopResponse struct {
	Stop   StopDetails `json:"stop"`
	Status Status      `json:"status"`
}

// StopDetails contains detailed stop information.
type StopDetails struct {
	StopID             int           `json:"stop_id"`
	StopName           string        `json:"stop_name"`
	StationType        string        `json:"station_type"`
	StationDescription string        `json:"station_description"`
	RouteType          int           `json:"route_type"`
	StopLocation       *StopLocation `json:"stop_location"`
	StopAmenities      *StopAmenity  `json:"stop_amenities"`
	StopAccessibility  *StopAccess   `json:"stop_accessibility"`
}

// StopLocation has lat/lon for a stop.
type StopLocation struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// StopAmenity describes amenities at a stop.
type StopAmenity struct {
	Toilet    bool `json:"toilet"`
	TaxiRank  bool `json:"taxi_rank"`
	CarParking string `json:"car_parking"`
	CCTV      bool `json:"cctv"`
}

// StopAccess describes accessibility at a stop.
type StopAccess struct {
	Lighting    bool `json:"lighting"`
	Stairs      bool `json:"stairs"`
	Escalator   bool `json:"escalator"`
	LiftAccess  bool `json:"lift_access"`
	Hearing     bool `json:"hearing_loop"`
	Wheelchair  bool `json:"wheelchair"`
}

// RoutesResponse is the response from GET /v3/routes.
type RoutesResponse struct {
	Routes []RouteWithStatus `json:"routes"`
	Status Status            `json:"status"`
}

// RouteResponse is the response from GET /v3/routes/{route_id}.
type RouteResponse struct {
	Route  RouteWithStatus `json:"route"`
	Status Status          `json:"status"`
}

// RouteWithStatus is a route with service status info.
type RouteWithStatus struct {
	RouteID            int                `json:"route_id"`
	RouteName          string             `json:"route_name"`
	RouteNumber        string             `json:"route_number"`
	RouteType          int                `json:"route_type"`
	RouteGTFSID        string             `json:"route_gtfs_id"`
	RouteServiceStatus *RouteServiceStatus `json:"route_service_status"`
}

// RouteServiceStatus is the service status of a route.
type RouteServiceStatus struct {
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

// DisruptionsResponse is the response from GET /v3/disruptions.
type DisruptionsResponse struct {
	Disruptions DisruptionCategories `json:"disruptions"`
	Status      Status               `json:"status"`
}

// DisruptionCategories groups disruptions by transport mode.
type DisruptionCategories struct {
	MetroTrain    []Disruption `json:"metro_train"`
	MetroTram     []Disruption `json:"metro_tram"`
	MetroBus      []Disruption `json:"metro_bus"`
	VLineTrain    []Disruption `json:"regional_train"`
	VLineCoach    []Disruption `json:"regional_coach"`
	VLineBus      []Disruption `json:"regional_bus"`
	SchoolBus     []Disruption `json:"school_bus"`
	Telebus       []Disruption `json:"telebus"`
	NightBus      []Disruption `json:"night_bus"`
	Ferry         []Disruption `json:"ferry"`
	Interstate    []Disruption `json:"interstate"`
	SkyBus        []Disruption `json:"skybus"`
	TaxiAndRideshare []Disruption `json:"taxi"`
	General       []Disruption `json:"general"`
}

// AllDisruptions returns all disruptions from all categories as a flat list.
func (dc DisruptionCategories) AllDisruptions() []Disruption {
	var all []Disruption
	all = append(all, dc.MetroTrain...)
	all = append(all, dc.MetroTram...)
	all = append(all, dc.MetroBus...)
	all = append(all, dc.VLineTrain...)
	all = append(all, dc.VLineCoach...)
	all = append(all, dc.VLineBus...)
	all = append(all, dc.SchoolBus...)
	all = append(all, dc.Telebus...)
	all = append(all, dc.NightBus...)
	all = append(all, dc.Ferry...)
	all = append(all, dc.Interstate...)
	all = append(all, dc.SkyBus...)
	all = append(all, dc.TaxiAndRideshare...)
	all = append(all, dc.General...)
	return all
}

// Disruption is a single disruption.
type Disruption struct {
	DisruptionID     int        `json:"disruption_id"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	DisruptionStatus string     `json:"disruption_status"`
	DisruptionType   string     `json:"disruption_type"`
	FromDate         *time.Time `json:"from_date"`
	ToDate           *time.Time `json:"to_date"`
	PublishedOn      *time.Time `json:"published_on"`
	LastUpdated      *time.Time `json:"last_updated"`
}

// RouteDisruptionsResponse is the response from GET /v3/disruptions/route/{route_id}.
type RouteDisruptionsResponse struct {
	Disruptions DisruptionCategories `json:"disruptions"`
	Status      Status               `json:"status"`
}

// StopDisruptionsResponse is the response from GET /v3/disruptions/stop/{stop_id}.
type StopDisruptionsResponse struct {
	Disruptions DisruptionCategories `json:"disruptions"`
	Status      Status               `json:"status"`
}

// FareEstimateResponse is the response from GET /v3/fare_estimate/...
type FareEstimateResponse struct {
	FareEstimate *FareEstimateResult `json:"fare_estimate"`
	Status       Status              `json:"status"`
}

// FareEstimateResult contains fare estimate details.
type FareEstimateResult struct {
	IsEarlyBird    bool              `json:"is_early_bird"`
	IsJourneyInFreeTramZone bool     `json:"is_journey_in_free_tram_zone"`
	IsTHSOnlyZone  bool              `json:"is_ths_only_zone"`
	PassengerFares []PassengerFare   `json:"passenger_fares"`
}

// PassengerFare is a fare for a passenger type.
type PassengerFare struct {
	PassengerType string  `json:"passenger_type"`
	Fare2Hour     float64 `json:"fare_2_hour"`
	FareDaily     float64 `json:"fare_daily"`
	FareWeekly    float64 `json:"fare_weekly"`
	FareMonthly   float64 `json:"fare_monthly"`
	Pass7Days     float64 `json:"pass_7_days"`
	Pass28To69DaysPerDay float64 `json:"pass_28_to_69_day_per_day"`
	Pass70PlusDaysPerDay float64 `json:"pass_70_plus_day_per_day"`
	FareWeekend   float64 `json:"fare_weekend_cap"`
	HolidayCap    float64 `json:"holiday_cap"`
}

// RouteTypesResponse is the response from GET /v3/route_types.
type RouteTypesResponse struct {
	RouteTypes []RouteType `json:"route_types"`
	Status     Status      `json:"status"`
}

// RouteType describes a route type.
type RouteType struct {
	RouteTypeName string `json:"route_type_name"`
	RouteTypeID   int    `json:"route_type"`
}
