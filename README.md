# city-timezones-go

A fast and lightweight Go library for looking up timezones by city name, with additional coordinate and Plus Code support.

This is a Go port of the popular [city-timezones](https://github.com/kevinroberts/city-timezones) Node.js library, featuring embedded compressed data and zero runtime dependencies.

## Installation

```bash
go get github.com/justcfx2u/city-timezones-go
```

## Usage

```go
import "github.com/justcfx2u/city-timezones-go"
```

## API Reference

### LookupViaCity(city string) []CityData

Finds cities by exact name match (case-insensitive). Returns an array of matching cities, or an empty slice if nothing matches.

```go
cities := citytimezones.LookupViaCity("Chicago")
fmt.Printf("Found %d matches for Chicago\n", len(cities))
if len(cities) > 0 {
    fmt.Printf("Timezone: %s\n", cities[0].Timezone)
    fmt.Printf("Coordinates: %f, %f\n", cities[0].Lat, cities[0].Lng)
}
```

### FindFromCityStateProvince(searchString string) []CityData

Performs partial matching across city, state/province, and country fields. Supports space-separated search terms.

```go
cities := citytimezones.FindFromCityStateProvince("springfield mo")
// Returns cities matching both "springfield" and "mo"

cities = citytimezones.FindFromCityStateProvince("London")
// Returns all cities named London across different countries
```

### FindFromIsoCode(isoCode string) []CityData

Finds cities by ISO2 or ISO3 country code (case-insensitive).

```go
germanCities := citytimezones.FindFromIsoCode("DE")
// or
germanCities = citytimezones.FindFromIsoCode("DEU")
```

### FindNearestCities(lat, lng, radiusKm float64) []CityData

Finds all cities within a specified radius of given coordinates. Results are sorted by distance (closest first).

```go
// Find cities within 50km of coordinates
cities := citytimezones.FindNearestCities(41.8299, -87.7500, 50.0)
```

### FindFromCoordinates(coords interface{}) []CityData

Flexible coordinate input supporting multiple formats. Uses a default 50km search radius.

```go
// String format
cities := citytimezones.FindFromCoordinates("41.8299,-87.7500")

// Array format  
cities = citytimezones.FindFromCoordinates([2]float64{41.8299, -87.7500})

// Slice format
cities = citytimezones.FindFromCoordinates([]float64{41.8299, -87.7500})
```

### FindFromPlusCode(plusCode string) []CityData

Finds cities near a location specified by a [Plus Code](https://plus.codes/) (Open Location Code).

```go
// Plus code for Chicago area
cities := citytimezones.FindFromPlusCode("86HJP27M+XF")
```

### GetCityMapping() []CityData

Returns the complete dataset of all cities (7300+ entries).

```go
allCities := citytimezones.GetCityMapping()
fmt.Printf("Total cities in database: %d\n", len(allCities))
```

## Data Structure

Each city is represented by a `CityData` struct:

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

## Features

- **Zero Dependencies**: Uses only Go standard library (plus Google's Plus Codes library)
- **Embedded Data**: City data is compressed and embedded at compile time (~267KB gzipped)
- **Fast Lookups**: In-memory operations with efficient filtering
- **Flexible Input**: Multiple coordinate input formats supported
- **Plus Codes Support**: Integration with Google's Open Location Code system
- **Cross-Platform**: Works on all platforms supported by Go
- **Thread-Safe**: All operations are read-only and safe for concurrent use

## Performance

- **Data Size**: ~1.9MB JSON compressed to ~267KB (86% compression)
- **Cities**: 7300+ cities worldwide with timezone information
- **Memory Usage**: Data loaded once at initialization
- **Lookup Speed**: Sub-millisecond performance for most operations

## Data Synchronization

Keep data up-to-date with the upstream repository:

```bash
go run cmd/sync-data/main.go
```

This tool downloads the latest city data from the upstream repository, validates it, and updates both the JSON file and compressed version.

## Development

```bash
# Run tests
go test -v

# Sync data from upstream
go run cmd/sync-data/main.go

# Debug data loading
go run cmd/debug-json/main.go
go run cmd/debug-embedded/main.go
```

## License

MIT License - Same as the original [city-timezones](https://github.com/kevinroberts/city-timezones) project.

## Credits

This project is a Go port of the excellent [city-timezones](https://github.com/kevinroberts/city-timezones) JavaScript library by **Kevin Roberts**.

- **Original Author**: [Kevin Roberts](https://github.com/kevinroberts)
- **Original Project**: [kevinroberts/city-timezones](https://github.com/kevinroberts/city-timezones) (Node.js/JavaScript)
- **Data Source**: City timezone data compiled and maintained by Kevin Roberts
- **License**: MIT (same as original project)
- **Plus Codes Support**: [Google's Open Location Code](https://github.com/google/open-location-code)

Special thanks to Kevin Roberts for creating and maintaining the original comprehensive city timezone database that makes this Go library possible.