package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const prowlarrBaseURL = "http://prowlarr.pi.local"

var apiKey string
var prowlarrHeaders map[string]string

type InitializeResponse struct {
	APIKey string `json:"apiKey"`
}

func getAPIKey() error {
	logger.Info("Retrieving Prowlarr API Key...")

	respBody, err := makeRequest("GET", prowlarrBaseURL+"/initialize.json", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	var initResp InitializeResponse
	if err := json.Unmarshal(respBody, &initResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	apiKey = initResp.APIKey
	prowlarrHeaders = map[string]string{
		"X-Api-Key": apiKey,
	}

	return nil
}

func loadJSONFile(filename string) (map[string]interface{}, error) {
	filePath := filepath.Join("prowlarr_req_bodies", filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filename, err)
	}

	return jsonData, nil
}

func configureHostSettings() error {
	logger.Info("Configuring host with login details...")

	hostConfig, err := loadJSONFile("host_config.json")
	if err != nil {
		return err
	}

	// Update host config with environment variables
	if username := os.Getenv("PROWLARR_USERNAME"); username == "" {
		return fmt.Errorf("missing env var: PROWLARR_USERNAME")
	} else {
		hostConfig["username"] = username
	}
	if password := os.Getenv("PROWLARR_PASSWORD"); password == "" {
		return fmt.Errorf("missing env var: PROWLARR_PASSWORD")
	} else {
		hostConfig["password"] = password
		hostConfig["passwordConfirmation"] = password
	}
	hostConfig["apiKey"] = apiKey

	_, err = makeRequest("PUT", prowlarrBaseURL+"/api/v1/config/host", hostConfig, prowlarrHeaders)
	return err
}

func configureDownloadClient() error {
	logger.Info("Configuring Download Client...")

	downloadClient, err := loadJSONFile("qbittorrent_downloadclient.json")
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

	_, err = makeRequest("POST", prowlarrBaseURL+"/api/v1/downloadclient", downloadClient, prowlarrHeaders)
	return err
}

func addIndexer(filename, name string) error {
	logger.Info("Adding indexer", "name", name)

	indexer, err := loadJSONFile(filename)
	if err != nil {
		return err
	}

	_, err = makeRequest("POST", prowlarrBaseURL+"/api/v1/indexer", indexer, prowlarrHeaders)
	return err
}

func ConfigureProwlarr() error {
	logger.Info("Starting Prowlarr configuration...")

	if err := getAPIKey(); err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

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
		logger.Info("Executing step", "step", step.name)
		if err := step.fn(); err != nil {
			logger.Error("Step failed", "step", step.name, "error", err)
			return err
		}
		logger.Info("Step completed successfully", "step", step.name)
	}

	logger.Info("Prowlarr configuration completed!")
	return nil
}
