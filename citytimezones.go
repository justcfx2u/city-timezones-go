package citytimezones

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	
	"github.com/google/open-location-code/go"
)

// CityData represents a city with timezone and location information
type CityData struct {
	City        string      `json:"city"`
	CityAscii   string      `json:"city_ascii"`
	Lat         float64     `json:"lat"`
	Lng         float64     `json:"lng"`
	Pop         interface{} `json:"pop"` // Can be int, float64, or empty
	Country     string      `json:"country"`
	ISO2        interface{} `json:"iso2"` // Can be string or number
	ISO3        interface{} `json:"iso3"` // Can be string or number
	Province    string      `json:"province"`
	Timezone    string      `json:"timezone"`
	StateAnsi   interface{} `json:"state_ansi,omitempty"` // Can be string or null
	ExactCity   interface{} `json:"exactCity,omitempty"`  // Can be string or null
	ExactProvince interface{} `json:"exactProvince,omitempty"` // Can be string or null
}

var cityMapping []CityData

// Initialize loads the city data (embedded first, then external file fallback)
func init() {
	err := loadCityData()
	if err != nil {
		panic(fmt.Sprintf("Failed to load city data: %v", err))
	}
}

// loadCityData tries embedded data first, then falls back to external file
func loadCityData() error {
	// Try embedded data first
	cities, err := loadEmbeddedCityData()
	if err == nil {
		cityMapping = cities
		return nil
	}
	
	// Fallback to external JSON file
	data, err := os.ReadFile("data/cityMap.json")
	if err != nil {
		return fmt.Errorf("failed to load both embedded and external data: %w", err)
	}
	
	return json.Unmarshal(data, &cityMapping)
}

// LookupViaCity finds cities by exact name match (case-insensitive)
func LookupViaCity(city string) []CityData {
	var results []CityData
	cityTrimmed := strings.TrimSpace(city)
	if cityTrimmed == "" {
		return results
	}
	cityLower := strings.ToLower(cityTrimmed)
	
	for _, c := range cityMapping {
		if strings.ToLower(c.City) == cityLower {
			results = append(results, c)
		}
	}
	
	return results
}

// findPartialMatch checks if all search terms are found in the target strings
func findPartialMatch(searchFields []string, searchString string) bool {
	searchTrimmed := strings.TrimSpace(searchString)
	if searchTrimmed == "" {
		return false
	}
	searchTerms := strings.Fields(strings.ToLower(searchTrimmed))
	joinedFields := strings.ToLower(strings.Join(searchFields, " "))
	
	for _, term := range searchTerms {
		if !strings.Contains(joinedFields, term) {
			return false
		}
	}
	
	return true
}

// FindFromCityStateProvince finds cities by partial matching across city/state/province/country
func FindFromCityStateProvince(searchString string) []CityData {
	if searchString == "" {
		return []CityData{}
	}
	
	var results []CityData
	
	for _, c := range cityMapping {
		searchFields := []string{c.City}
		
		if c.StateAnsi != nil {
			if stateAnsi, ok := c.StateAnsi.(string); ok && stateAnsi != "" {
				searchFields = append(searchFields, stateAnsi)
			}
		}
		
		searchFields = append(searchFields, c.Province, c.Country)
		
		if findPartialMatch(searchFields, searchString) {
			results = append(results, c)
		}
	}
	
	return results
}

// FindFromIsoCode finds cities by ISO2 or ISO3 country code
func FindFromIsoCode(isoCode string) []CityData {
	isoTrimmed := strings.TrimSpace(isoCode)
	if isoTrimmed == "" {
		return []CityData{}
	}
	
	var results []CityData
	isoLower := strings.ToLower(isoTrimmed)
	
	for _, c := range cityMapping {
		iso2Match := false
		iso3Match := false
		
		if c.ISO2 != nil {
			if iso2Str, ok := c.ISO2.(string); ok {
				iso2Match = strings.ToLower(iso2Str) == isoLower
			}
		}
		
		if c.ISO3 != nil {
			if iso3Str, ok := c.ISO3.(string); ok {
				iso3Match = strings.ToLower(iso3Str) == isoLower
			}
		}
		
		if iso2Match || iso3Match {
			results = append(results, c)
		}
	}
	
	return results
}

// GetCityMapping returns the complete city dataset
func GetCityMapping() []CityData {
	return cityMapping
}

// CityDistance represents a city with its distance from a reference point
type CityDistance struct {
	CityData
	Distance float64 // in kilometers
}

// haversineDistance calculates the distance between two points on Earth using the Haversine formula
func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

// FindNearestCities finds all cities within a specified radius (in kilometers) of the given coordinates
func FindNearestCities(lat, lng, radiusKm float64) []CityData {
	var results []CityDistance

	for _, city := range cityMapping {
		distance := haversineDistance(lat, lng, city.Lat, city.Lng)
		if distance <= radiusKm {
			results = append(results, CityDistance{
				CityData: city,
				Distance: distance,
			})
		}
	}

	// Sort by distance (closest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	// Extract just the CityData - always return a slice, never nil
	cities := make([]CityData, 0, len(results))
	for _, result := range results {
		cities = append(cities, result.CityData)
	}

	return cities
}

// parseCoordinates attempts to parse coordinates from various input types
func parseCoordinates(coords interface{}) (lat, lng float64, err error) {
	switch v := coords.(type) {
	case string:
		// Parse "lat,lng" format
		parts := strings.Split(v, ",")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid coordinate string format, expected 'lat,lng'")
		}
		
		lat, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid latitude: %v", err)
		}
		
		lng, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid longitude: %v", err)
		}
		
		return lat, lng, nil
		
	case [2]float64:
		return v[0], v[1], nil
		
	case []float64:
		if len(v) != 2 {
			return 0, 0, fmt.Errorf("coordinate slice must have exactly 2 elements")
		}
		return v[0], v[1], nil
		
	default:
		return 0, 0, fmt.Errorf("unsupported coordinate type: %T", coords)
	}
}

// FindFromCoordinates finds the nearest cities to the given coordinates (flexible input)
// Supports string "lat,lng", [2]float64{lat, lng}, or []float64{lat, lng}
func FindFromCoordinates(coords interface{}) []CityData {
	lat, lng, err := parseCoordinates(coords)
	if err != nil {
		return []CityData{}
	}
	
	// Default radius of 50km for coordinate searches
	return FindNearestCities(lat, lng, 50.0)
}

// FindFromPlusCode finds cities near the location specified by a Plus Code (Open Location Code)
func FindFromPlusCode(plusCode string) []CityData {
	plusCodeTrimmed := strings.TrimSpace(plusCode)
	if plusCodeTrimmed == "" {
		return []CityData{}
	}
	
	// Decode the plus code to get coordinates
	area, err := olc.Decode(plusCodeTrimmed)
	if err != nil {
		return []CityData{}
	}
	
	// Use the center of the plus code area
	centerLat := area.LatLo + (area.LatHi-area.LatLo)/2
	centerLng := area.LngLo + (area.LngHi-area.LngLo)/2
	
	// Default radius of 50km for plus code searches
	return FindNearestCities(centerLat, centerLng, 50.0)
}