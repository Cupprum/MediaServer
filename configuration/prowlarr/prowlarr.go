package prowlarr

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"MediaServer/configuration/utils"
)

//go:embed req_bodies/*.json
var reqBodies embed.FS

type ProwlarrConfig struct {
	Apikey              string // Set during login
	Url                 string
	Username            string
	Password            string
	QBittorrentHostname string
	QBittorrentUsername string
	QBittorrentPassword string
}

func Config() (*ProwlarrConfig, error) {
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

	return &ProwlarrConfig{
		Apikey:              "", // Set during login
		Url:                 url,
		Username:            username,
		Password:            password,
		QBittorrentHostname: QBittorrentHostname,
		QBittorrentUsername: qbittorrentUsername,
		QBittorrentPassword: qbittorrentPassword,
	}, nil
}

func (c *ProwlarrConfig) Login() error {
	fmt.Println("-- Logging in...")

	// Create a cookie jar for persisting cookies across login and subsequent requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create cookie jar: %w", err)
	}
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	b := fmt.Sprintf("username=%s&password=%s&rememberMe=on", c.Username, c.Password)
	_, err = utils.Request("POST", c.Url+"/login", b, nil, client)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	err = c.apikey(client)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	return nil
}

func (c *ProwlarrConfig) apikey(client *http.Client) error {
	fmt.Println("-- Retrieving apikey...")
	rb, err := utils.Request("GET", c.Url+"/initialize.json", nil, nil, client)
	if err != nil {
		return fmt.Errorf("failed to get apikey: %w", err)
	}

	var r struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		// If not logged in, prowlarr returns html
		if err.Error() == "invalid character '<' looking for beginning of value" {
			return fmt.Errorf("not logged in")
		}
		return fmt.Errorf("failed to decode apikey: %w", err)
	}

	// Set the API Key for Prowlarr
	c.Apikey = r.APIKey
	if c.Apikey == "" {
		return fmt.Errorf("apikey cannot be empty")
	}

	return nil
}

func (c *ProwlarrConfig) setHostSetting() error {
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

func (c *ProwlarrConfig) setDownloadClient() error {
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

func (c *ProwlarrConfig) setIndexer(filename, name string) error {
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

	c, err := Config()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// If not configured, the first call needs to be `initialization`
	// If prowlarr is already configured, the apikey retrieval will fail
	err = c.apikey(nil)
	if err != nil {
		if err.Error() == "not logged in" {
			fmt.Println("- prowlarr already configured, skipping...")
			fmt.Println()
			return nil
		} else {
			return fmt.Errorf("failed to retrieve apikey: %w", err)
		}
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

	fmt.Println("- Prowlarr configured successfully!")
	fmt.Println()

	return nil
}
