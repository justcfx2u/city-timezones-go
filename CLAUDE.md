# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Run all tests
go test -v

# Sync data from upstream repository
go run cmd/sync-data/main.go

# Debug data loading
go run cmd/debug-json/main.go
go run cmd/debug-embedded/main.go
```

## Code Architecture

This is a Go library that provides timezone lookup functionality for cities worldwide. The architecture consists of:

- **Entry Point**: `citytimezones.go` - Contains 6 main functions exported as package
- **Embedded Data**: `embedded.go` - Compressed gzipped JSON data (267KB, 86% compression)  
- **Data Source**: `data/cityMap.json` - JSON array with 7300+ city records
- **Sync Tool**: `cmd/sync-data/main.go` - Cross-platform data synchronization utility

### Core Functions

1. `LookupViaCity(city string) []CityData` - Exact city name matching (case-insensitive)
2. `FindFromCityStateProvince(searchString string) []CityData` - Partial matching with space-separated terms
3. `FindFromIsoCode(isoCode string) []CityData` - ISO2/ISO3 country code lookup
4. `FindNearestCities(lat, lng, radiusKm float64) []CityData` - Coordinate-based radius search
5. `FindFromCoordinates(coords interface{}) []CityData` - Flexible coordinate input (string/array/slice)
6. `FindFromPlusCode(plusCode string) []CityData` - Google Plus Codes (Open Location Code) support
7. `GetCityMapping() []CityData` - Direct access to complete dataset

All functions return slices (empty slice when no matches found) for consistent handling.

### Key Implementation Features

- **Zero Runtime Dependencies**: Uses only Go standard library + Google Plus Codes
- **Embedded Compressed Data**: gzipped JSON embedded at compile time
- **Thread-Safe**: All operations are read-only and concurrent-safe  
- **Flexible Input Types**: Multiple coordinate formats supported
- **Haversine Distance**: Accurate distance calculations for coordinate lookups
- **Case Insensitive**: All string comparisons normalized to lowercase
- **Consistent Returns**: Always return slices, never nil

### Data Structure

Each city record contains:
```go
type CityData struct {
    City        string      `json:"city"`           // Display name
    CityAscii   string      `json:"city_ascii"`     // ASCII version  
    Lat         float64     `json:"lat"`            // Latitude
    Lng         float64     `json:"lng"`            // Longitude
    Pop         interface{} `json:"pop"`            // Population (varies by type)
    Country     string      `json:"country"`        // Full country name
    ISO2        interface{} `json:"iso2"`           // ISO2 country code  
    ISO3        interface{} `json:"iso3"`           // ISO3 country code
    Province    string      `json:"province"`       // State/province
    Timezone    string      `json:"timezone"`       // IANA timezone identifier
    StateAnsi   interface{} `json:"state_ansi,omitempty"`    // US state abbreviation
    ExactCity   interface{} `json:"exactCity,omitempty"`     // Alternative city name
    ExactProvince interface{} `json:"exactProvince,omitempty"` // Alternative province
}
```

### Testing Strategy

- Uses Go's built-in testing framework (`citytimezones_test.go`)
- 18 comprehensive test cases covering all functions
- Tests include exact matching, coordinate lookups, Plus Codes, edge cases
- Validates deterministic results with specific latitude checks
- Tests coordinate input flexibility and distance calculations

### Development Notes

- **Build**: No build process required - direct Go execution
- **Dependencies**: Only `github.com/google/open-location-code/go` for Plus Codes
- **Data Management**: Embedded compressed data + sync utility
- **Performance**: Sub-millisecond lookups, ~267KB memory footprint
- **Threading**: Thread-safe operations for concurrent use
- **Data Source**: Syncs from upstream [kevinroberts/city-timezones](https://github.com/kevinroberts/city-timezones)

### Project Structure

```
city-timezones-go/
├── go.mod, go.sum              # Go module definition
├── citytimezones.go            # Main library implementation
├── embedded.go                 # Compressed data embedding
├── citytimezones_test.go       # Comprehensive test suite
├── data/                       # JSON data files
│   ├── cityMap.json           # Raw JSON data
│   └── cityMap.json.gz        # Compressed data for embedding
├── cmd/                       # Utilities
│   ├── sync-data/             # Data sync utility
│   ├── debug-json/            # JSON debugging tool
│   └── debug-embedded/        # Embedded data debugging tool
├── README.md                  # Complete API documentation
└── CLAUDE.md                  # Development guidance
```

This is a complete, production-ready Go library that provides timezone lookup functionality with embedded data and enhanced coordinate/Plus Code support.