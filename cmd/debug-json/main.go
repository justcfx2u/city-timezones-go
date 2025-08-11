package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type TestCityData struct {
	City        string      `json:"city"`
	CityAscii   string      `json:"city_ascii"`
	Lat         float64     `json:"lat"`
	Lng         float64     `json:"lng"`
	Pop         interface{} `json:"pop"`
	Country     string      `json:"country"`
	ISO2        interface{} `json:"iso2"` // Test as interface
	ISO3        interface{} `json:"iso3"` // Test as interface
	Province    string      `json:"province"`
	Timezone    string      `json:"timezone"`
	StateAnsi   interface{} `json:"state_ansi,omitempty"`
	ExactCity   interface{} `json:"exactCity,omitempty"`
	ExactProvince interface{} `json:"exactProvince,omitempty"`
}

func main() {
	data, err := os.ReadFile("data/cityMap.json")
	if err != nil {
		panic(err)
	}
	
	var cities []TestCityData
	if err := json.Unmarshal(data, &cities); err != nil {
		panic(err)
	}
	
	fmt.Printf("Successfully loaded %d cities\n", len(cities))
	
	// Check for problematic ISO codes
	for i, city := range cities[:100] { // Check first 100 cities
		fmt.Printf("City %d: %s, ISO2: %v (%T), ISO3: %v (%T)\n", 
			i, city.City, city.ISO2, city.ISO2, city.ISO3, city.ISO3)
	}
}