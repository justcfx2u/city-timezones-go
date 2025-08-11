package main

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
)

//go:embed ../../data/cityMap.json.gz
var embeddedCityData []byte

type TestCityData struct {
	City        string      `json:"city"`
	CityAscii   string      `json:"city_ascii"`
	Lat         float64     `json:"lat"`
	Lng         float64     `json:"lng"`
	Pop         interface{} `json:"pop"`
	Country     string      `json:"country"`
	ISO2        interface{} `json:"iso2"`
	ISO3        interface{} `json:"iso3"`
	Province    string      `json:"province"`
	Timezone    string      `json:"timezone"`
	StateAnsi   interface{} `json:"state_ansi,omitempty"`
	ExactCity   interface{} `json:"exactCity,omitempty"`
	ExactProvince interface{} `json:"exactProvince,omitempty"`
}

func main() {
	fmt.Printf("Embedded data size: %d bytes\n", len(embeddedCityData))
	
	// Create gzip reader
	gzReader, err := gzip.NewReader(bytes.NewReader(embeddedCityData))
	if err != nil {
		panic(fmt.Errorf("failed to create gzip reader: %w", err))
	}
	defer gzReader.Close()
	
	// Read decompressed data
	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		panic(fmt.Errorf("failed to decompress data: %w", err))
	}
	
	fmt.Printf("Decompressed data size: %d bytes\n", len(decompressed))
	
	// Parse JSON
	var cities []TestCityData
	if err := json.Unmarshal(decompressed, &cities); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	
	fmt.Printf("Successfully loaded %d cities from embedded data\n", len(cities))
}