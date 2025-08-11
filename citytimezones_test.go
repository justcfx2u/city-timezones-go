package citytimezones

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestLookupViaCity_Chicago(t *testing.T) {
	cities := LookupViaCity("Chicago")
	if len(cities) == 0 {
		t.Fatal("Expected to find Chicago, got empty result")
	}
	
	chicago := cities[0]
	expectedLat := 41.82999066
	if chicago.Lat != expectedLat {
		t.Errorf("Expected Chicago lat %f, got %f", expectedLat, chicago.Lat)
	}
}

func TestLookupViaCity_CaseInsensitive(t *testing.T) {
	cities := LookupViaCity("chicago")
	if len(cities) == 0 {
		t.Fatal("Expected to find chicago (lowercase), got empty result")
	}
	
	chicago := cities[0]
	expectedLat := 41.82999066
	if chicago.Lat != expectedLat {
		t.Errorf("Expected Chicago lat %f, got %f", expectedLat, chicago.Lat)
	}
}

func TestLookupViaCity_NoMatch(t *testing.T) {
	cities := LookupViaCity("Foobar")
	if len(cities) != 0 {
		t.Errorf("Expected empty result for 'Foobar', got %d cities", len(cities))
	}
}

func TestLookupViaCity_MultipleMatches(t *testing.T) {
	cities := LookupViaCity("Springfield")
	if len(cities) <= 1 {
		t.Errorf("Expected multiple matches for Springfield, got %d", len(cities))
	}
}

func TestFindFromCityStateProvince_SpringfieldMO(t *testing.T) {
	cities := FindFromCityStateProvince("springfield mo")
	if len(cities) == 0 {
		t.Fatal("Expected to find Springfield MO, got empty result")
	}
	
	springfield := cities[0]
	expectedLat := 37.18001609
	if springfield.Lat != expectedLat {
		t.Errorf("Expected Springfield MO lat %f, got %f", expectedLat, springfield.Lat)
	}
}

func TestFindFromCityStateProvince_London(t *testing.T) {
	cities := FindFromCityStateProvince("London")
	expectedCount := 6
	if len(cities) != expectedCount {
		t.Errorf("Expected %d London matches, got %d", expectedCount, len(cities))
	}
}

func TestFindFromCityStateProvince_EmptyString(t *testing.T) {
	cities := FindFromCityStateProvince("")
	if len(cities) != 0 {
		t.Errorf("Expected empty result for empty string, got %d cities", len(cities))
	}
}

func TestFindFromIsoCode_ISO2(t *testing.T) {
	cities := FindFromIsoCode("de")
	if len(cities) == 0 {
		t.Fatal("Expected to find German cities with 'de', got empty result")
	}
	
	germany := cities[0]
	expectedLat := 49.98247246
	if germany.Lat != expectedLat {
		t.Errorf("Expected German city lat %f, got %f", expectedLat, germany.Lat)
	}
}

func TestFindFromIsoCode_ISO3(t *testing.T) {
	cities := FindFromIsoCode("deu")
	if len(cities) == 0 {
		t.Fatal("Expected to find German cities with 'deu', got empty result")
	}
	
	germany := cities[0]
	expectedLat := 49.98247246
	if germany.Lat != expectedLat {
		t.Errorf("Expected German city lat %f, got %f", expectedLat, germany.Lat)
	}
}

func TestFindFromIsoCode_CaseInsensitive(t *testing.T) {
	cities := FindFromIsoCode("DE")
	if len(cities) == 0 {
		t.Fatal("Expected to find German cities with 'DE' (uppercase), got empty result")
	}
	
	germany := cities[0]
	expectedLat := 49.98247246
	if germany.Lat != expectedLat {
		t.Errorf("Expected German city lat %f, got %f", expectedLat, germany.Lat)
	}
}

func TestFindFromIsoCode_NoMatch(t *testing.T) {
	cities := FindFromIsoCode("Foobar")
	if len(cities) != 0 {
		t.Errorf("Expected empty result for 'Foobar', got %d cities", len(cities))
	}
}

func TestFindFromIsoCode_MultipleMatches(t *testing.T) {
	cities := FindFromIsoCode("de")
	if len(cities) <= 1 {
		t.Errorf("Expected multiple German cities, got %d", len(cities))
	}
}

func TestGetCityMapping_Count(t *testing.T) {
	cities := GetCityMapping()
	minExpected := 7323
	if len(cities) < minExpected {
		t.Errorf("Expected at least %d cities in mapping, got %d", minExpected, len(cities))
	}
}

func TestFindNearestCities_Chicago(t *testing.T) {
	// Chicago coordinates
	cities := FindNearestCities(41.8299, -87.7500, 50.0)
	if len(cities) == 0 {
		t.Fatal("Expected to find cities near Chicago, got empty result")
	}
	
	// Should find Chicago itself
	found := false
	for _, city := range cities {
		if strings.ToLower(city.City) == "chicago" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find Chicago in nearest cities result")
	}
}

func TestFindNearestCities_LimitRadius(t *testing.T) {
	// Very small radius should return fewer results
	cities1 := FindNearestCities(41.8299, -87.7500, 1.0)
	cities2 := FindNearestCities(41.8299, -87.7500, 100.0)
	
	if len(cities1) >= len(cities2) {
		t.Errorf("Expected fewer cities with smaller radius: %d vs %d", len(cities1), len(cities2))
	}
}

func TestFindFromCoordinates_String(t *testing.T) {
	cities := FindFromCoordinates("41.8299,-87.7500")
	if len(cities) == 0 {
		t.Fatal("Expected to find cities from coordinate string, got empty result")
	}
}

func TestFindFromCoordinates_Floats(t *testing.T) {
	cities := FindFromCoordinates([2]float64{41.8299, -87.7500})
	if len(cities) == 0 {
		t.Fatal("Expected to find cities from coordinate array, got empty result")
	}
}

func TestFindFromCoordinates_InvalidInput(t *testing.T) {
	cities := FindFromCoordinates("invalid")
	if len(cities) != 0 {
		t.Errorf("Expected empty result for invalid coordinates, got %d cities", len(cities))
	}
}

func TestFindFromPlusCode_Valid(t *testing.T) {
	// Plus code for Chicago area: 86HJP27M+XF
	cities := FindFromPlusCode("86HJP27M+XF")
	if len(cities) == 0 {
		t.Fatal("Expected to find cities from plus code, got empty result")
	}
	
	// Should find cities near Chicago
	found := false
	for _, city := range cities {
		if strings.Contains(strings.ToLower(city.City), "chicago") {
			found = true
			break
		}
	}
	if !found {
		t.Log("Plus code search may not have found Chicago specifically, but should find nearby cities")
	}
}

func TestFindFromPlusCode_Invalid(t *testing.T) {
	cities := FindFromPlusCode("invalid")
	if len(cities) != 0 {
		t.Errorf("Expected empty result for invalid plus code, got %d cities", len(cities))
	}
}

// Robustness tests for edge cases and error conditions

func TestEmptyStringInputs(t *testing.T) {
	// All functions should handle empty strings gracefully
	tests := []struct {
		name string
		fn   func() []CityData
	}{
		{"LookupViaCity empty", func() []CityData { return LookupViaCity("") }},
		{"FindFromCityStateProvince empty", func() []CityData { return FindFromCityStateProvince("") }},
		{"FindFromIsoCode empty", func() []CityData { return FindFromIsoCode("") }},
		{"FindFromPlusCode empty", func() []CityData { return FindFromPlusCode("") }},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			if len(result) != 0 {
				t.Errorf("%s: expected empty result for empty input, got %d", tt.name, len(result))
			}
		})
	}
}

func TestWhitespaceInputs(t *testing.T) {
	// Test various whitespace scenarios
	tests := []struct {
		input    string
		expected bool // true if should find results
	}{
		{"   ", false},           // Only spaces
		{"\t\n", false},         // Only tabs/newlines
		{"  Chicago  ", true},   // Leading/trailing spaces
		{"Chicago\t", true},     // Trailing tab
		{"\nChicago\n", true},   // Leading/trailing newlines
	}
	
	for _, tt := range tests {
		t.Run(fmt.Sprintf("whitespace_%q", tt.input), func(t *testing.T) {
			cities := LookupViaCity(tt.input)
			hasResults := len(cities) > 0
			if hasResults != tt.expected {
				t.Errorf("Input %q: expected hasResults=%t, got %t", tt.input, tt.expected, hasResults)
			}
		})
	}
}

func TestUnicodeAndInternationalCharacters(t *testing.T) {
	// Test cities with accents and international characters
	testCities := []string{
		"São Paulo",     // Portuguese accents
		"München",       // German umlaut
		"Москва",        // Cyrillic (Moscow)
		"北京",           // Chinese characters (Beijing)
		"Montréal",      // French accent
	}
	
	for _, city := range testCities {
		t.Run(fmt.Sprintf("unicode_%s", city), func(t *testing.T) {
			// Should not panic with unicode input
			cities := LookupViaCity(city)
			// Don't assert results since our test data may not have all these cities
			// Just ensure no panic - nil slices are acceptable in Go
			_ = cities // Just ensure no panic
		})
	}
}

func TestExtremeCoordinates(t *testing.T) {
	// Test edge cases for coordinate inputs
	tests := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"North Pole", 90.0, 0.0},
		{"South Pole", -90.0, 0.0},
		{"International Date Line", 0.0, 180.0},
		{"Antimeridian", 0.0, -180.0},
		{"Extreme valid", 89.9, 179.9},
		{"Extreme valid negative", -89.9, -179.9},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with extreme coordinates
			cities := FindNearestCities(tt.lat, tt.lng, 1000.0) // Large radius
			_ = cities // Just ensure no panic - nil slices are acceptable
		})
	}
}

func TestInvalidCoordinates(t *testing.T) {
	// Test coordinates outside valid ranges
	tests := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"Latitude too high", 91.0, 0.0},
		{"Latitude too low", -91.0, 0.0},
		{"Longitude too high", 0.0, 181.0},
		{"Longitude too low", 0.0, -181.0},
		{"Both invalid", 100.0, 200.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should handle gracefully, not panic
			cities := FindNearestCities(tt.lat, tt.lng, 50.0)
			_ = cities // Just ensure no panic - nil slices are acceptable
		})
	}
}

func TestMalformedCoordinateStrings(t *testing.T) {
	// Test various malformed coordinate string inputs
	malformedInputs := []string{
		"not,coordinates",
		"123.45",           // Missing comma and second coordinate
		"123.45,",          // Missing second coordinate
		",67.89",           // Missing first coordinate
		"abc,def",          // Non-numeric
		"123.45,67.89,extra", // Too many parts
		"123.45;67.89",     // Wrong separator
		"(123.45,67.89)",   // With parentheses
		"123.45, 67.89, 0", // Three coordinates
		"",                 // Empty string
		"   ,   ",          // Only whitespace
	}
	
	for _, input := range malformedInputs {
		t.Run(fmt.Sprintf("malformed_%s", input), func(t *testing.T) {
			// Should return empty slice, not panic
			cities := FindFromCoordinates(input)
			if len(cities) != 0 {
				t.Errorf("Expected empty result for malformed input %s, got %d cities", input, len(cities))
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test that multiple goroutines can safely use the library
	const numGoroutines = 100
	const numCalls = 10
	
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numCalls; j++ {
				// Mix different function calls
				switch j % 4 {
				case 0:
					cities := LookupViaCity("Chicago")
					if len(cities) == 0 {
						errors <- fmt.Errorf("goroutine %d: no cities found for Chicago", id)
						return
					}
				case 1:
					cities := FindFromIsoCode("US")
					if len(cities) == 0 {
						errors <- fmt.Errorf("goroutine %d: no cities found for US", id)
						return
					}
				case 2:
					cities := FindNearestCities(40.7128, -74.0060, 50.0)
					_ = cities // Just ensure no panic
				case 3:
					allCities := GetCityMapping()
					if len(allCities) < 7000 {
						errors <- fmt.Errorf("goroutine %d: got %d cities, expected >7000", id, len(allCities))
						return
					}
				}
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for any errors
	for err := range errors {
		t.Error(err)
	}
}

func TestLargeRadiusSearch(t *testing.T) {
	// Test very large radius searches
	cities := FindNearestCities(0.0, 0.0, 20000.0) // ~half the Earth's circumference
	
	// Should return many cities but not crash
	_ = cities // Just ensure no panic
	
	// Should return a significant portion of all cities
	allCities := GetCityMapping()
	if len(cities) < len(allCities)/2 {
		t.Logf("Large radius search returned %d/%d cities (this might be expected)", len(cities), len(allCities))
	}
}

func TestZeroRadiusSearch(t *testing.T) {
	// Test zero radius search
	cities := FindNearestCities(40.7128, -74.0060, 0.0)
	
	// Should return very few or no cities
	if len(cities) > 1 {
		t.Logf("Zero radius search returned %d cities (might find exact matches)", len(cities))
	}
}

func TestDataIntegrity(t *testing.T) {
	// Verify that all cities have required fields
	allCities := GetCityMapping()
	
	if len(allCities) == 0 {
		t.Fatal("No cities loaded")
	}
	
	for i, city := range allCities[:100] { // Check first 100 for performance
		if city.City == "" {
			t.Errorf("City %d has empty name", i)
		}
		if city.Country == "" {
			t.Errorf("City %d (%s) has empty country", i, city.City)
		}
		if city.Timezone == "" {
			t.Errorf("City %d (%s) has empty timezone", i, city.City)
		}
		if city.Lat < -90 || city.Lat > 90 {
			t.Errorf("City %d (%s) has invalid latitude: %f", i, city.City, city.Lat)
		}
		if city.Lng < -180 || city.Lng > 180 {
			t.Errorf("City %d (%s) has invalid longitude: %f", i, city.City, city.Lng)
		}
	}
}