package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ProwlarrConfig struct {
	Url                 string
	Username            string
	Password            string
	Apikey              string
	QBittorrentHostname string
	QBittorrentUsername string
	QBittorrentPassword string
}

func getProwlarrConfig() (*ProwlarrConfig, error) {
	fmt.Println("-- Loading Prowlarr config...")

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

	apikey := os.Getenv("PROWLARR_APIKEY")
	if apikey == "" {
		fmt.Println(" * env var: `PROWLARR_APIKEY` is not defined, this is expected during first run")
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
		Apikey:              apikey,
		QBittorrentHostname: QBittorrentHostname,
		QBittorrentUsername: qbittorrentUsername,
		QBittorrentPassword: qbittorrentPassword,
	}, nil
}

func (c *ProwlarrConfig) prowlarrLogin() error {
	fmt.Println("-- Logging in to Prowlarr...")

	b := fmt.Sprintf("username=%s&password=%s&rememberMe=on", c.Username, c.Password)

	req, err := requestBuilder("POST", c.Url+"/login", b, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Response verification
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status)
	}

	if string(resp.Request.URL.RawQuery) != url.PathEscape("returnUrl=/") {
		return fmt.Errorf("request failed: error 401: invalid returnUrl")
	}

	return nil
}

func (c *ProwlarrConfig) prowlarrInitialize() error {
	fmt.Println("-- Retrieving Prowlarr API Key...")
	rb, err := Request("GET", c.Url+"/initialize.json", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	var r struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	c.Apikey = r.APIKey

	updateDotEnv("PROWLARR_APIKEY", c.Apikey)

	return nil
}

func (c *ProwlarrConfig) configureProwlarrHostSettings() error {
	fmt.Println("-- Configuring Prowlarr with login details...")

	b, err := loadJSONFile("prowlarr", "host_config.json")
	if err != nil {
		return err
	}

	b["username"] = c.Username
	b["password"] = strings.TrimSpace(c.Password)
	b["passwordConfirmation"] = strings.TrimSpace(c.Password)
	b["apiKey"] = c.Apikey

	h := map[string]string{"X-Api-Key": c.Apikey}

	_, err = Request("PUT", c.Url+"/api/v1/config/host", b, h, nil)
	return err
}

func (c *ProwlarrConfig) configureProwlarrDownloadClient() error {
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

	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err = Request("POST", c.Url+"/api/v1/downloadclient", b, h, nil)
	return err
}

func (c *ProwlarrConfig) addProwlarrIndexer(filename, name string) error {
	fmt.Println("-- Adding indexer:", name)

	b, err := loadJSONFile("prowlarr", filename)
	if err != nil {
		return err
	}

	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err = Request("POST", c.Url+"/api/v1/indexer", b, h, nil)
	return err
}

func ConfigureProwlarr() error {
	fmt.Println("- Starting Prowlarr configuration...")

	c, err := getProwlarrConfig()
	if err != nil {
		return err
	}

	err = c.prowlarrInitialize()
	if err != nil {
		fmt.Println("- Prowlarr already configured, skipping...")
		fmt.Println()
		return nil
	}

	// // Otherwise, proceed with configuration
	if err = c.configureProwlarrHostSettings(); err != nil {
		return err
	}
	if err = c.configureProwlarrDownloadClient(); err != nil {
		return err
	}
	if err = c.addProwlarrIndexer("pirate_bay_indexer.json", "Pirate Bay"); err != nil {
		return err
	}
	// if err = c.addProwlarrIndexer("eztv_indexer.json", "EZTV"); err != nil {
	// 	return err
	// }
	if err = c.addProwlarrIndexer("limetorrents_indexer.json", "Limetorrents"); err != nil {
		return err
	}
	if err = c.addProwlarrIndexer("yts_indexer.json", "YTS"); err != nil {
		return err
	}

	fmt.Println("- Prowlarr configured successfully!")
	fmt.Println()
	return nil
}
