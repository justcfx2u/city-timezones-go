package citytimezones

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
)

// Embedded compressed city data (build-time embedded)
//
//go:embed data/cityMap.json.gz
var embeddedCityData []byte

// loadEmbeddedCityData loads and decompresses the embedded city data
func loadEmbeddedCityData() ([]CityData, error) {
	// Create gzip reader
	gzReader, err := gzip.NewReader(bytes.NewReader(embeddedCityData))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()
	
	// Read decompressed data
	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}
	
	// Parse JSON
	var cities []CityData
	if err := json.Unmarshal(decompressed, &cities); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	return cities, nil
}