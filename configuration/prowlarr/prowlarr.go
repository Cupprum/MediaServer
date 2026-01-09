package prowlarr

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"MediaServer/configuration/utils"
)

//go:embed req_bodies/*.json
var reqBodies embed.FS

type Config struct {
	Apikey              string // Set during login
	Url                 string
	Username            string
	Password            string
	QBittorrentHostname string
	QBittorrentUsername string
	QBittorrentPassword string
}

func GetConfig() (*Config, error) {
	fmt.Println("-- Loading config...")

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

	return &Config{
		Apikey:              "", // Set during login
		Url:                 url,
		Username:            username,
		Password:            password,
		QBittorrentHostname: QBittorrentHostname,
		QBittorrentUsername: qbittorrentUsername,
		QBittorrentPassword: qbittorrentPassword,
	}, nil
}

func (c *Config) LoadApikey(client *http.Client) error {
	fmt.Println("-- Retrieving apikey...")
	// First call to initialize does not require authentication
	// the subsequent calls do authentication via cookie in client
	rb, err := utils.Request("GET", c.Url+"/initialize.json", nil, nil, client)
	if err != nil {
		return fmt.Errorf("failed to get apikey: %w", err)
	}

	var r struct {
		Apikey string `json:"apiKey"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		// After initial setup, if not logged in, prowlarr returns html
		if err.Error() == "invalid character '<' looking for beginning of value" {
			return fmt.Errorf("configured: authentication required")
		}
		return fmt.Errorf("failed to decode apikey: %w", err)
	}

	// Set the apikey in the config
	c.Apikey = r.Apikey
	if c.Apikey == "" {
		return fmt.Errorf("apikey cannot be empty")
	}

	return nil
}

func (c *Config) setHostSetting() error {
	fmt.Println("-- Set login details...")

	b, err := utils.LoadJSONFile(reqBodies, "host_config.json")
	if err != nil {
		return err
	}

	b["username"] = c.Username
	b["password"] = strings.TrimSpace(c.Password)
	b["passwordConfirmation"] = strings.TrimSpace(c.Password)
	b["apiKey"] = c.Apikey

	h := map[string]string{"X-Api-Key": c.Apikey}

	_, err = utils.Request("PUT", c.Url+"/api/v1/config/host", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to configure login details")
	}

	return nil
}

func (c *Config) setDownloadClient() error {
	fmt.Println("-- Configuring Download Client...")

	b, err := utils.LoadJSONFile(reqBodies, "qbittorrent_downloadclient.json")
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
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
	_, err = utils.Request("POST", c.Url+"/api/v1/downloadclient", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to set downloadclient: %w", err)
	}

	return nil
}

func (c *Config) setIndexer(filename, name string) error {
	fmt.Printf("-- Adding indexer: %v...\n", name)

	b, err := utils.LoadJSONFile(reqBodies, filename)
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err = utils.Request("POST", c.Url+"/api/v1/indexer", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to set indexer %v: %w", name, err)
	}

	return nil
}

func Configure() error {
	fmt.Println("- Starting Prowlarr configuration...")

	c, err := GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// If not configured, the first call needs to be `initialization` endpoint
	// If prowlarr is already configured, the apikey retrieval will fail
	err = c.LoadApikey(nil)
	if err != nil {
		if strings.Contains(err.Error(), "configured") {
			fmt.Println("- already configured, skipping...")
			fmt.Println()
			return nil
		}
		return fmt.Errorf("failed to load api key: %w", err)
	}

	if err = c.setHostSetting(); err != nil {
		return err
	}
	if err = c.setDownloadClient(); err != nil {
		return err
	}
	if err = c.setIndexer("pirate_bay_indexer.json", "Pirate Bay"); err != nil {
		return err
	}
	// if err = c.setIndexer("eztv_indexer.json", "EZTV"); err != nil {
	// 	return err
	// }
	if err = c.setIndexer("limetorrents_indexer.json", "Limetorrents"); err != nil {
		return err
	}
	if err = c.setIndexer("yts_indexer.json", "YTS"); err != nil {
		return err
	}

	fmt.Println("- prowlarr configured successfully!")
	fmt.Println()

	return nil
}
