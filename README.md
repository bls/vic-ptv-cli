# vic-ptv-cli

A command-line interface for the [Public Transport Victoria (PTV) Timetable API v3](https://www.ptv.vic.gov.au/footer/data-and-reporting/datasets/ptv-timetable-api/).

Query stops, routes, departures, disruptions, and fares from the PTV network covering trains, trams, buses, V/Line, and more.

## Installation

### From source

```bash
go install github.com/bls/vic-ptv-cli@latest
```

### Binary releases

Download the latest release for your platform from the [Releases](https://github.com/bls/vic-ptv-cli/releases) page.

## Quick Start

### 1. Get API credentials

Email PTV at **APIKeyRequest@ptv.vic.gov.au** with your name and reason for access. You'll receive a Developer ID and API Key.

### 2. Configure credentials

Set environment variables:

```bash
export PTV_DEV_ID=your_dev_id
export PTV_API_KEY=your_api_key
```

Or create a config file at `~/.config/vic-ptv-cli/config.yaml`:

```yaml
devId: your_dev_id
apiKey: your_api_key
```

### 3. Start querying

```bash
# Search for a stop
ptv search "Flinders Street"

# Get upcoming train departures
ptv departures 1071 --route-type 0

# View route details
ptv route 1
```

## Commands

### `ptv search <term>`

Search for stops, routes, and outlets matching a term.

```bash
ptv search "Southern Cross"
ptv search "96" --route-types 1          # Search tram routes only
ptv search "Melbourne" --route-types 0,1  # Trains and trams
```

**Flags:**
- `--route-types` — Filter by route types (comma-separated: 0=train, 1=tram, 2=bus, 3=vline_train, 4=vline_coach)

### `ptv departures <stop_id>`

Show upcoming departures from a stop.

```bash
ptv departures 1071 --route-type 0           # Train departures from Flinders St
ptv departures 1071 --route-type 0 --limit 10 # Show more results
```

**Flags:**
- `--route-type` — Route type (required): 0=train, 1=tram, 2=bus, 3=vline_train, 4=vline_coach
- `--limit` — Maximum departures to show (default: 5)

### `ptv stop <stop_id>`

Show stop details including amenities and accessibility.

```bash
ptv stop 1071 --route-type 0
```

**Flags:**
- `--route-type` — Route type (required)

### `ptv routes`

List all routes, optionally filtered by type.

```bash
ptv routes
ptv routes --type 0    # Train routes only
ptv routes --type 1,2  # Tram and bus routes
```

**Flags:**
- `--type` — Filter by route type (comma-separated)

### `ptv route <route_id>`

Show details for a specific route.

```bash
ptv route 1
ptv route 725
```

### `ptv disruptions`

Show current service disruptions.

```bash
ptv disruptions
ptv disruptions --route 1    # Disruptions for a specific route
ptv disruptions --stop 1071  # Disruptions at a specific stop
```

**Flags:**
- `--route` — Filter by route ID
- `--stop` — Filter by stop ID

### `ptv fare <min_zone> <max_zone>`

Estimate fares between myki zones.

```bash
ptv fare 1 1    # Zone 1 only
ptv fare 1 2    # Zone 1+2
```

### `ptv route-types`

List all route types and their numeric IDs.

```bash
ptv route-types
```

### `ptv config`

Show current configuration status.

```bash
ptv config
```

## Global Flags

All commands support these flags:

- `--json` — Output raw JSON from the API
- `--dev-id` — PTV Developer ID (overrides env/config)
- `--api-key` — PTV API Key (overrides env/config)

## Configuration

Credentials are loaded in this priority order:

1. CLI flags (`--dev-id`, `--api-key`)
2. Environment variables (`PTV_DEV_ID`, `PTV_API_KEY`)
3. Config file (`~/.config/vic-ptv-cli/config.yaml`)

### Config file format

```yaml
devId: "1000001"
apiKey: "aaaabbbb-cccc-dddd-eeee-ffffffaaaaaa"
```

## Route Types

| ID | Type |
|----|------|
| 0  | Train |
| 1  | Tram |
| 2  | Bus |
| 3  | V/Line Train |
| 4  | V/Line Coach |

## Attribution

Data licensed from Public Transport Victoria under Creative Commons Attribution 3.0 Australia Licence.
