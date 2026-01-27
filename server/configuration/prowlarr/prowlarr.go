package prowlarr

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"MediaServer/server/configuration/utils"
)

//go:embed req_bodies/*.json
var reqBodies embed.FS

type Config struct {
	Apikey                string // Set during login
	Url                   string
	Username              string
	Password              string
	QBittorrentHostname   string
	QBittorrentUsername   string
	QBittorrentPassword   string
	FlaresolverrHostUrl   string
	Deploy1337x           bool
	DeployEztv            bool
	DeployInternetArchive bool
	DeployLimetorrents    bool
	DeployPirateBay       bool
	DeployYts             bool
	DeployRutracker       bool
	RutrackerUsername     string
	RutrackerPassword     string
	DeploySkTorrent       bool
	SkTorrentUsername     string
	SkTorrentPassword     string
	DeploySkczTorrent     bool
	SkCzTorrentUsername   string
	SkCzTorrentPassword   string
}

func GetConfig() (*Config, error) {
	log.Println("-- Loading config...")

	url := os.Getenv("MEDIASERVER_PROWLARR_URL")
	if url == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_URL`")
	}

	username := os.Getenv("MEDIASERVER_PROWLARR_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_USERNAME`")
	}

	password := os.Getenv("MEDIASERVER_PROWLARR_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_PASSWORD`")
	}

	qbittorrentHostname := os.Getenv("MEDIASERVER_QBITTORRENT_HOSTNAME")
	if qbittorrentHostname == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_QBITTORRENT_HOSTNAME`")
	}

	qbittorrentUsername := os.Getenv("MEDIASERVER_QBITTORRENT_USERNAME")
	if qbittorrentUsername == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_QBITTORRENT_USERNAME`")
	}

	qbittorrentPassword := os.Getenv("MEDIASERVER_QBITTORRENT_PASSWORD")
	if qbittorrentPassword == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_QBITTORRENT_PASSWORD`")
	}

	flaresolverrHostUrl := os.Getenv("MEDIASERVER_FLARESOLVERR_HOST_URL")
	if flaresolverrHostUrl == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_FLARESOLVERR_HOST_URL`")
	}

	deploy1337x := os.Getenv("MEDIASERVER_PROWLARR_1337X_ENABLED") == "true"
	deployEztv := os.Getenv("MEDIASERVER_PROWLARR_EZTV_ENABLED") == "true"
	deployInternetArchive := os.Getenv("MEDIASERVER_PROWLARR_INTERNETARCHIVE_ENABLED") == "true"
	deployLimetorrents := os.Getenv("MEDIASERVER_PROWLARR_LIMETORRENTS_ENABLED") == "true"
	deployPirateBay := os.Getenv("MEDIASERVER_PROWLARR_PIRATEBAY_ENABLED") == "true"
	deployYts := os.Getenv("MEDIASERVER_PROWLARR_YTS_ENABLED") == "true"
	deployRutracker := os.Getenv("MEDIASERVER_PROWLARR_RUTRACKER_ENABLED") == "true"

	rutrackerUsername := os.Getenv("MEDIASERVER_PROWLARR_RUTRACKER_USERNAME")
	if rutrackerUsername == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_RUTRACKER_USERNAME`")
	}

	rutrackerPassword := os.Getenv("MEDIASERVER_PROWLARR_RUTRACKER_PASSWORD")
	if rutrackerPassword == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_RUTRACKER_PASSWORD`")
	}

	deploySkTorrent := os.Getenv("MEDIASERVER_PROWLARR_SKTORRENT_ENABLED") == "true"
	skTorrentUsername := os.Getenv("MEDIASERVER_PROWLARR_SKTORRENT_USERNAME")
	if skTorrentUsername == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_SKTORRENT_USERNAME`")
	}
	skTorrentPassword := os.Getenv("MEDIASERVER_PROWLARR_SKTORRENT_PASSWORD")
	if skTorrentPassword == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_SKTORRENT_PASSWORD`")
	}

	deploySkczTorrent := os.Getenv("MEDIASERVER_PROWLARR_SKCZTORRENT_ENABLED") == "true"
	skczTorrentUsername := os.Getenv("MEDIASERVER_PROWLARR_SKCZTORRENT_USERNAME")
	if skczTorrentUsername == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_SKCZTORRENT_USERNAME`")
	}
	skczTorrentPassword := os.Getenv("MEDIASERVER_PROWLARR_SKCZTORRENT_PASSWORD")
	if skczTorrentPassword == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_PROWLARR_SKCZTORRENT_PASSWORD`")
	}

	return &Config{
		Apikey:                "", // Set during login
		Url:                   url,
		Username:              username,
		Password:              password,
		QBittorrentHostname:   qbittorrentHostname,
		QBittorrentUsername:   qbittorrentUsername,
		QBittorrentPassword:   qbittorrentPassword,
		FlaresolverrHostUrl:   flaresolverrHostUrl,
		Deploy1337x:           deploy1337x,
		DeployEztv:            deployEztv,
		DeployInternetArchive: deployInternetArchive,
		DeployLimetorrents:    deployLimetorrents,
		DeployPirateBay:       deployPirateBay,
		DeployYts:             deployYts,
		DeployRutracker:       deployRutracker,
		RutrackerUsername:     rutrackerUsername,
		RutrackerPassword:     rutrackerPassword,
		DeploySkTorrent:       deploySkTorrent,
		SkTorrentUsername:     skTorrentUsername,
		SkTorrentPassword:     skTorrentPassword,
		DeploySkczTorrent:     deploySkczTorrent,
		SkCzTorrentUsername:   skczTorrentUsername,
		SkCzTorrentPassword:   skczTorrentPassword,
	}, nil
}

func setField(b map[string]interface{}, key string, value string) {
	if fields, ok := b["fields"].([]interface{}); ok {
		for _, field := range fields {
			if fieldMap, ok := field.(map[string]interface{}); ok {
				name, _ := fieldMap["name"].(string)
				if name == key {
					fieldMap["value"] = value
				}
			}
		}
	}
}

func (c *Config) LoadApikey(client *http.Client) error {
	log.Println("-- Retrieving apikey...")
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
	log.Println("-- Set login details...")

	b, err := utils.LoadJSONFile(reqBodies, "host_config.json")
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
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
	log.Println("-- Configuring Download Client...")

	// Configure qBittorrent download client with the following host: `qbittorrent.server.svc.cluster.local`
	b, err := utils.LoadJSONFile(reqBodies, "qbittorrent_downloadclient.json")
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	// Update download client config
	setField(b, "host", c.QBittorrentHostname)
	setField(b, "username", c.QBittorrentUsername)
	setField(b, "password", c.QBittorrentPassword)

	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err = utils.Request("POST", c.Url+"/api/v1/downloadclient", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to set downloadclient: %w", err)
	}

	return nil
}

func (c *Config) addTag(name string) error {
	log.Printf("-- Adding tag: %v...\n", name)

	b := struct {
		Label string `json:"label"`
	}{name}

	h := map[string]string{"X-Api-Key": c.Apikey}
	rb, err := utils.Request("POST", c.Url+"/api/v1/tag", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to add tag %v: %w", name, err)
	}

	var r struct {
		Id int `json:"id"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return fmt.Errorf("failed to decode tag id: %w", err)
	}

	if r.Id != 1 {
		return fmt.Errorf("tag id has to be 1, received: %v", r.Id)
	}

	return nil
}

func (c *Config) setIndexerProxy() error {
	log.Println("-- Configuring Flaresolverr Indexer Proxy...")

	b, err := utils.LoadJSONFile(reqBodies, "flaresolverr_indexer_proxy.json")
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	// Update indexer proxy config
	setField(b, "host", c.FlaresolverrHostUrl)

	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err = utils.Request("POST", c.Url+"/api/v1/indexerProxy", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to set indexer proxy: %w", err)
	}

	return nil
}

func (c *Config) setIndexer(name string, body any) error {
	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err := utils.Request("POST", c.Url+"/api/v1/indexer", body, h, nil)
	if err != nil {
		return fmt.Errorf("failed to set indexer %v: %w", name, err)
	}

	return nil
}

func (c *Config) setPublicIndexer(filename, name string) error {
	log.Printf("-- Adding public indexer: %v...\n", name)

	b, err := utils.LoadJSONFile(reqBodies, filename)
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	err = c.setIndexer(name, b)
	return err
}

func (c *Config) setPrivateIndexer(filename, name, username, password string) error {
	log.Printf("-- Adding private indexer: %v...\n", name)

	b, err := utils.LoadJSONFile(reqBodies, filename)
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	// Update indexer config
	setField(b, "username", username)
	setField(b, "password", password)

	err = c.setIndexer(name, b)
	return err
}

func Configure() error {
	log.Println("- Starting Prowlarr configuration...")

	c, err := GetConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// If not configured, the first call needs to be `initialization` endpoint
	// If prowlarr is already configured, the apikey retrieval will fail
	err = c.LoadApikey(nil)
	if err != nil {
		if strings.Contains(err.Error(), "configured") {
			log.Println("- already configured, skipping...")
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
	if err = c.addTag("flaresolverr"); err != nil {
		return err
	}
	if err = c.setIndexerProxy(); err != nil {
		return err
	}

	if c.Deploy1337x {
		if err = c.setPublicIndexer("1337x_indexer.json", "1337x"); err != nil {
			return err
		}
	}
	if c.DeployEztv {
		if err = c.setPublicIndexer("eztv_indexer.json", "EZTV"); err != nil {
			return err
		}
	}
	if c.DeployInternetArchive {
		if err = c.setPublicIndexer("internetarchive_indexer.json", "Internet Archive"); err != nil {
			return err
		}
	}
	if c.DeployLimetorrents {
		if err = c.setPublicIndexer("limetorrents_indexer.json", "LimeTorrents"); err != nil {
			return err
		}
	}
	if c.DeployPirateBay {
		if err = c.setPublicIndexer("pirate_bay_indexer.json", "The Pirate Bay"); err != nil {
			return err
		}
	}
	if c.DeployYts {
		if err = c.setPublicIndexer("yts_indexer.json", "YTS"); err != nil {
			return err
		}
	}
	if c.DeployRutracker {
		err = c.setPrivateIndexer("rutracker_indexer.json", "RuTracker.org", c.RutrackerUsername, c.RutrackerPassword)
		if err != nil {
			return err
		}
	}
	if c.DeploySkczTorrent {
		err = c.setPrivateIndexer("skcztorrent_indexer.json", "Sk-CzTorrent", c.SkCzTorrentUsername, c.SkCzTorrentPassword)
		if err != nil {
			return err
		}
	}
	if c.DeploySkTorrent {
		err = c.setPrivateIndexer("sktorrent_indexer.json", "SkTorrent.org", c.SkTorrentUsername, c.SkTorrentPassword)
		if err != nil {
			return err
		}
	}

	log.Println("- prowlarr configured successfully!")
	fmt.Println()

	return nil
}
