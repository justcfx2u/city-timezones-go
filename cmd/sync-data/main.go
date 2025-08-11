package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	upstreamURL    = "https://raw.githubusercontent.com/kevinroberts/city-timezones/master/data/cityMap.json"
	localJSONPath  = "data/cityMap.json"
	localGZPath    = "data/cityMap.json.gz"
)

func main() {
	fmt.Println("Syncing city data from upstream repository...")

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Printf("ERROR: Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	// Download latest data
	fmt.Printf("Downloading latest cityMap.json from %s\n", upstreamURL)
	resp, err := http.Get(upstreamURL)
	if err != nil {
		fmt.Printf("ERROR: Failed to download data: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: HTTP %d: %s\n", resp.StatusCode, resp.Status)
		os.Exit(1)
	}

	// Read response body
	newData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: Failed to read response: %v\n", err)
		os.Exit(1)
	}

	// Validate JSON
	var testData []interface{}
	if err := json.Unmarshal(newData, &testData); err != nil {
		fmt.Printf("ERROR: Downloaded file is not valid JSON: %v\n", err)
		os.Exit(1)
	}

	// Check if the new file is different from the current one
	if existingData, err := os.ReadFile(localJSONPath); err == nil {
		if bytes.Equal(existingData, newData) {
			fmt.Println("No changes detected in upstream data")
			return
		}
		fmt.Println("Changes detected in upstream data")
	} else {
		fmt.Println("No existing local data found")
	}

	// Write the new JSON file
	if err := os.WriteFile(localJSONPath, newData, 0644); err != nil {
		fmt.Printf("ERROR: Failed to write %s: %v\n", localJSONPath, err)
		os.Exit(1)
	}
	fmt.Printf("Updated %s\n", localJSONPath)

	// Create compressed version for embedding
	fmt.Println("Creating compressed version for embedding...")
	if err := compressFile(localJSONPath, localGZPath); err != nil {
		fmt.Printf("ERROR: Failed to compress file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created %s\n", localGZPath)

	// Show file sizes and compression ratio
	showFileStats()

	fmt.Println()
	fmt.Println("Sync completed successfully!")
	fmt.Println("Note: Remember to rebuild to update embedded data.")
}

func compressFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Create destination file
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	// Write compressed data
	_, err = gzWriter.Write(data)
	return err
}

func showFileStats() {
	fmt.Println()
	fmt.Println("File sizes:")

	jsonInfo, err := os.Stat(localJSONPath)
	if err == nil {
		fmt.Printf("  %s: %s\n", localJSONPath, formatBytes(jsonInfo.Size()))
	}

	gzInfo, err := os.Stat(localGZPath)
	if err == nil {
		fmt.Printf("  %s: %s\n", localGZPath, formatBytes(gzInfo.Size()))

		if jsonInfo != nil {
			ratio := 100.0 * (1.0 - float64(gzInfo.Size())/float64(jsonInfo.Size()))
			fmt.Printf("Compression ratio: %.1f%%\n", ratio)
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}