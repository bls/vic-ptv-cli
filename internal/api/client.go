package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://timetableapi.ptv.vic.gov.au"
	DefaultTimeout = 10 * time.Second
)

// Client is a PTV API client.
type Client struct {
	BaseURL    string
	DevID      string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new PTV API client.
func NewClient(devID, apiKey string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		DevID:   devID,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// get makes a signed GET request to the given API path and decodes the response.
func (c *Client) get(path string, result interface{}) error {
	signedURL, err := SignURL(c.BaseURL, path, c.DevID, c.APIKey)
	if err != nil {
		return fmt.Errorf("signing URL: %w", err)
	}

	resp, err := c.HTTPClient.Get(signedURL)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, httpErrorMessage(resp.StatusCode, string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

func httpErrorMessage(code int, body string) string {
	switch code {
	case 400:
		return "Bad request — check your parameters"
	case 403:
		return "Forbidden — check your API credentials (devid/apikey)"
	case 404:
		return "Not found — the requested resource does not exist"
	case 429:
		return "Rate limited — too many requests, please wait"
	case 500:
		return "Server error — PTV API is having issues"
	case 503:
		return "Service unavailable — PTV API is temporarily down"
	default:
		if body != "" {
			return body
		}
		return fmt.Sprintf("unexpected status code %d", code)
	}
}

// Search performs a search for stops, routes, and outlets.
func (c *Client) Search(term string, routeTypes []int) (*SearchResponse, error) {
	path := fmt.Sprintf("/v3/search/%s", url.PathEscape(term))
	if len(routeTypes) > 0 {
		parts := make([]string, len(routeTypes))
		for i, rt := range routeTypes {
			parts[i] = strconv.Itoa(rt)
		}
		path += "?route_types=" + strings.Join(parts, "&route_types=")
	}
	var resp SearchResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Departures gets upcoming departures from a stop.
func (c *Client) Departures(routeType, stopID, maxResults int) (*DeparturesResponse, error) {
	path := fmt.Sprintf("/v3/departures/route_type/%d/stop/%d?max_results=%d&expand=route&expand=direction&expand=stop",
		routeType, stopID, maxResults)
	var resp DeparturesResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Stop gets details for a specific stop.
func (c *Client) Stop(stopID, routeType int) (*StopResponse, error) {
	path := fmt.Sprintf("/v3/stops/%d/route_type/%d?stop_location=true&stop_amenities=true&stop_accessibility=true",
		stopID, routeType)
	var resp StopResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Routes lists all routes, optionally filtered by route types.
func (c *Client) Routes(routeTypes []int) (*RoutesResponse, error) {
	path := "/v3/routes"
	if len(routeTypes) > 0 {
		parts := make([]string, len(routeTypes))
		for i, rt := range routeTypes {
			parts[i] = strconv.Itoa(rt)
		}
		path += "?route_types=" + strings.Join(parts, "&route_types=")
	}
	var resp RoutesResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Route gets details for a specific route.
func (c *Client) Route(routeID int) (*RouteResponse, error) {
	path := fmt.Sprintf("/v3/routes/%d", routeID)
	var resp RouteResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Disruptions gets current disruptions.
func (c *Client) Disruptions() (*DisruptionsResponse, error) {
	path := "/v3/disruptions"
	var resp DisruptionsResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DisruptionsByRoute gets disruptions for a specific route.
func (c *Client) DisruptionsByRoute(routeID int) (*DisruptionsResponse, error) {
	path := fmt.Sprintf("/v3/disruptions/route/%d", routeID)
	var resp DisruptionsResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DisruptionsByStop gets disruptions for a specific stop.
func (c *Client) DisruptionsByStop(stopID int) (*DisruptionsResponse, error) {
	path := fmt.Sprintf("/v3/disruptions/stop/%d", stopID)
	var resp DisruptionsResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FareEstimate gets fare estimates between zones.
func (c *Client) FareEstimate(minZone, maxZone int) (*FareEstimateResponse, error) {
	path := fmt.Sprintf("/v3/fare_estimate/min_zone/%d/max_zone/%d", minZone, maxZone)
	var resp FareEstimateResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// RouteTypes lists all route types.
func (c *Client) RouteTypes() (*RouteTypesResponse, error) {
	path := "/v3/route_types"
	var resp RouteTypesResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
