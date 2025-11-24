package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ProwlarrConfig struct {
	Url                 string
	Username            string
	Password            string
	QBittorrentHostname string
	QBittorrentUsername string
	QBittorrentPassword string
}

func getProwlarrConfig() (*ProwlarrConfig, error) {
	fmt.Println("-- Create Prowlarr config based on Environment Variables...")

	url := os.Getenv("PROWLARR_URL")
	if url == "" {
		return nil, fmt.Errorf("missing env var: `PROWLARR_URL`")
	}

	username := os.Getenv("PROWLARR_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing env var: `PROWLARR_USERNAME`")
	}

	password := os.Getenv("PROWLARR_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("missing env var: `PROWLARR_PASSWORD`")
	}

	QBittorrentHostname := os.Getenv("QBITTORRENT_HOSTNAME")
	if QBittorrentHostname == "" {
		return nil, fmt.Errorf("missing env var: `QBITTORRENT_HOSTNAME`")
	}

	qbittorrentUsername := os.Getenv("QBITTORRENT_USERNAME")
	if qbittorrentUsername == "" {
		return nil, fmt.Errorf("missing env var: `QBITTORRENT_USERNAME`")
	}

	qbittorrentPassword := os.Getenv("QBITTORRENT_PASSWORD")
	if qbittorrentPassword == "" {
		return nil, fmt.Errorf("missing env var: `QBITTORRENT_PASSWORD`")
	}

	return &ProwlarrConfig{
		Url:                 url,
		Username:            username,
		Password:            password,
		QBittorrentHostname: QBittorrentHostname,
		QBittorrentUsername: qbittorrentUsername,
		QBittorrentPassword: qbittorrentPassword,
	}, nil
}

// Used to cache the Prowlarr API Key
var prowlarrApiKey string = ""

func (c *ProwlarrConfig) getProwlarrHeaders() (map[string]string, error) {
	if prowlarrApiKey != "" {
		return map[string]string{"X-Api-Key": prowlarrApiKey}, nil
	}
	fmt.Println("-- Retrieving Prowlarr API Key...")

	rb, err := Request("GET", c.Url+"/initialize.json", nil, nil, nil)
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

func (c *ProwlarrConfig) configureHostSettings() error {
	fmt.Println("-- Configuring host with login details...")

	// Get headers first, as API Key is needed in body
	h, err := c.getProwlarrHeaders()
	if err != nil {
		return err
	}

	b, err := loadJSONFile("prowlarr", "host_config.json")
	if err != nil {
		return err
	}

	b["username"] = c.Username
	b["password"] = c.Password
	b["passwordConfirmation"] = c.Password
	b["apiKey"] = h["X-Api-Key"] // Set API key from headers

	_, err = Request("PUT", c.Url+"/api/v1/config/host", b, h, nil)
	return err
}

func (c *ProwlarrConfig) configureDownloadClient() error {
	fmt.Println("-- Configuring Download Client...")

	b, err := loadJSONFile("prowlarr", "qbittorrent_downloadclient.json")
	if err != nil {
		return err
	}

	// Update download client config
	if fields, ok := b["fields"].([]interface{}); ok {
		for _, field := range fields {
			if fieldMap, ok := field.(map[string]interface{}); ok {
				name, _ := fieldMap["name"].(string)
				switch name {
				case "host":
					fieldMap["value"] = c.QBittorrentHostname
				case "username":
					fieldMap["value"] = c.QBittorrentUsername
				case "password":
					fieldMap["value"] = c.QBittorrentPassword
				}
			}
		}
	}

	h, err := c.getProwlarrHeaders()
	if err != nil {
		return err
	}
	_, err = Request("POST", c.Url+"/api/v1/downloadclient", b, h, nil)
	return err
}

func (c *ProwlarrConfig) addIndexer(filename, name string) error {
	fmt.Println("-- Adding indexer:", name)

	b, err := loadJSONFile("prowlarr", filename)
	if err != nil {
		return err
	}

	h, err := c.getProwlarrHeaders()
	if err != nil {
		return err
	}
	_, err = Request("POST", c.Url+"/api/v1/indexer", b, h, nil)
	return err
}

func ConfigureProwlarr() error {
	fmt.Println("- Starting Prowlarr configuration...")

	c, err := getProwlarrConfig()
	if err != nil {
		return err
	}

	if err = c.configureHostSettings(); err != nil {
		return err
	}
	if err = c.configureDownloadClient(); err != nil {
		return err
	}
	if err = c.addIndexer("pirate_bay_indexer.json", "Pirate Bay"); err != nil {
		return err
	}
	if err = c.addIndexer("eztv_indexer.json", "EZTV"); err != nil {
		return err
	}
	if err = c.addIndexer("limetorrents_indexer.json", "Limetorrents"); err != nil {
		return err
	}
	if err = c.addIndexer("yts_indexer.json", "YTS"); err != nil {
		return err
	}

	fmt.Println("- Prowlarr configured successfully!")
	return nil
}
