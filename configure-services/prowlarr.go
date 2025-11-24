package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const prowlarrBaseURL = "http://prowlarr.pi.local"

// Used to cache the Prowlarr API Key
var prowlarrApiKey string = ""

func getProwlarrHeaders() (map[string]string, error) {
	if prowlarrApiKey != "" {
		return map[string]string{"X-Api-Key": prowlarrApiKey}, nil
	}
	fmt.Println("Retrieving Prowlarr API Key...")

	rb, err := Request("GET", prowlarrBaseURL+"/initialize.json", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	var r struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	prowlarrApiKey = r.APIKey

	return map[string]string{"X-Api-Key": prowlarrApiKey}, nil
}

func configureHostSettings() error {
	fmt.Println("Configuring host with login details...")

	// Update host config with environment variables
	u := os.Getenv("PROWLARR_USERNAME")
	if u == "" {
		return fmt.Errorf("missing env var: PROWLARR_USERNAME")
	}

	pw := os.Getenv("PROWLARR_PASSWORD")
	if pw == "" {
		return fmt.Errorf("missing env var: PROWLARR_PASSWORD")
	}

	h, err := getProwlarrHeaders() // Fetch headers with API key
	if err != nil {
		return err
	}

	hostConfig, err := loadJSONFile("prowlarr", "host_config.json")
	if err != nil {
		return err
	}

	hostConfig["username"] = u
	hostConfig["password"] = pw
	hostConfig["passwordConfirmation"] = pw
	hostConfig["apiKey"] = h["X-Api-Key"] // Set API key from headers

	_, err = Request("PUT", prowlarrBaseURL+"/api/v1/config/host", hostConfig, h)
	return err
}

func configureDownloadClient() error {
	fmt.Println("Configuring Download Client...")

	downloadClient, err := loadJSONFile("prowlarr", "qbittorrent_downloadclient.json")
	if err != nil {
		return err
	}

	// Update download client config with environment variables
	if fields, ok := downloadClient["fields"].([]interface{}); ok {
		for _, field := range fields {
			if fieldMap, ok := field.(map[string]interface{}); ok {
				name, _ := fieldMap["name"].(string)
				switch name {
				case "host":
					if host := os.Getenv("QBITTORRENT_HOST"); host != "" {
						fieldMap["value"] = host
					}
				case "username":
					if username := os.Getenv("QBITTORRENT_USERNAME"); username != "" {
						fieldMap["value"] = username
					}
				case "password":
					if password := os.Getenv("QBITTORRENT_PASSWORD"); password != "" {
						fieldMap["value"] = password
					}
				}
			}
		}
	}

	h, err := getProwlarrHeaders()
	if err != nil {
		return err
	}
	_, err = Request("POST", prowlarrBaseURL+"/api/v1/downloadclient", downloadClient, h)
	return err
}

func addIndexer(filename, name string) error {
	fmt.Println("Adding indexer", "name", name)

	indexer, err := loadJSONFile("prowlarr", filename)
	if err != nil {
		return err
	}

	h, err := getProwlarrHeaders()
	if err != nil {
		return err
	}
	_, err = Request("POST", prowlarrBaseURL+"/api/v1/indexer", indexer, h)
	return err
}

func ConfigureProwlarr() error {
	fmt.Println("Starting Prowlarr configuration...")

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Configure Host Settings", configureHostSettings},
		{"Configure Download Client", configureDownloadClient},
		{"Add Pirate Bay Indexer", func() error { return addIndexer("pirate_bay_indexer.json", "Pirate Bay") }},
		{"Add EZTV Indexer", func() error { return addIndexer("eztv_indexer.json", "EZTV") }},
		{"Add Limetorrents Indexer", func() error { return addIndexer("limetorrents_indexer.json", "Limetorrents") }},
		{"Add YTS Indexer", func() error { return addIndexer("yts_indexer.json", "YTS") }},
	}

	for _, step := range steps {
		fmt.Println("Executing step", "step", step.name)
		if err := step.fn(); err != nil {
			logger.Error("Step failed", "step", step.name, "error", err)
			return err
		}
		fmt.Println("Step completed successfully", "step", step.name)
	}

	fmt.Println("Prowlarr configuration completed!")
	return nil
}
