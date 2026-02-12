# CLAUDE.md — Build Instructions for vic-ptv-cli

## What This Is
A Go CLI for the Public Transport Victoria (PTV) Timetable API v3.
Single binary, no runtime deps. Targets Linux/macOS/Windows.

## API Reference
- Swagger spec: `ptv-api-spec.json` (already downloaded in this repo)
- Base URL: `https://timetableapi.ptv.vic.gov.au`
- Auth: HMAC-SHA1 signature. Every request must include `devid` and `signature` query params.
  - Build the URL path with query params + `devid=XXXX`
  - HMAC-SHA1 the full path (including /v3/...) using the API key as the secret
  - Append `&signature=<hex-encoded-hmac>` to the URL

## Auth Config
Support these in priority order:
1. `--dev-id` / `--api-key` CLI flags
2. `PTV_DEV_ID` / `PTV_API_KEY` environment variables  
3. Config file: `~/.config/vic-ptv-cli/config.yaml` (devId + apiKey fields)

On first run with no config, print a helpful message about how to get API keys (email PTV).

## CLI Framework
Use Cobra for CLI commands + Viper for config. Standard Go CLI patterns.

## Commands to Implement

### `ptv search <term>`
Search for stops, routes, outlets matching a term.
- API: `GET /v3/search/{search_term}`
- Show type (stop/route/outlet), name, ID, route type
- Support `--route-types` filter (0=train, 1=tram, 2=bus, 3=vline_train, 4=vline_coach)

### `ptv departures <stop_id> [--route-type <type>]`
Show upcoming departures from a stop.
- API: `GET /v3/departures/route_type/{route_type}/stop/{stop_id}`
- Show: scheduled time, estimated time (if available), route name, direction, platform
- Default: next 5 departures. `--limit N` to change.
- `--route-type` required (infer from stop if possible, or require it)

### `ptv stop <stop_id> [--route-type <type>]`
Show stop details/facilities.
- API: `GET /v3/stops/{stop_id}/route_type/{route_type}`

### `ptv routes [--type <route_type>]`
List routes, optionally filtered by type.
- API: `GET /v3/routes`

### `ptv route <route_id>`
Show route details.
- API: `GET /v3/routes/{route_id}`

### `ptv disruptions [--route <route_id>] [--stop <stop_id>]`
Show current disruptions.
- API: `GET /v3/disruptions` (or filtered variants)

### `ptv fare <min_zone> <max_zone>`
Estimate fare between zones.
- API: `GET /v3/fare_estimate/min_zone/{minZone}/max_zone/{maxZone}`

### `ptv config`
Show current config (dev ID set? config file location? etc.)

### `ptv route-types`
List route types and their IDs.

## Output Format
- Default: human-readable table/list output (use tabwriter or similar)
- `--json` flag on all commands for raw JSON output
- Times should be displayed in local timezone by default

## Project Structure
```
cmd/           # Cobra command files
  root.go
  search.go
  departures.go
  ...
internal/
  api/         # PTV API client (auth, HTTP, response types)
  config/      # Config loading (flags, env, file)
  display/     # Output formatting (table, JSON)
go.mod
go.sum
main.go
README.md
```

## Error Handling
- Clear error messages for missing auth
- HTTP error codes mapped to human messages
- Timeout handling (default 10s)

## README.md
Generate a proper README with:
- What it is
- Installation (go install + binary releases)
- Quick start (get API key, configure, first search)
- All commands with examples
- Config file format
- Attribution: "Data licensed from Public Transport Victoria under Creative Commons Attribution 3.0 Australia Licence"

## Don't
- Don't embed any API keys
- Don't add a proxy/relay feature
- Don't add journey planning (the API doesn't support it natively)
- Don't add real-time data (not available in this API)
