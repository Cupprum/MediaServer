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

type Indexer struct {
	Name     string
	Enabled  bool
	File     string
	Username string // optional, used in private trackers
	Password string // optional, used in private trackers
}

type Config struct {
	Apikey              string // Set during login
	Url                 string
	Username            string
	Password            string
	QBittorrentHostname string
	QBittorrentUsername string
	QBittorrentPassword string
	FlaresolverrHostUrl string
	Indexers            []Indexer
}

func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing env var: %s", key)
	}
	return val, nil
}

func GetConfig() (*Config, error) {
	log.Println("-- Loading config...")

	var err error
	c := &Config{}
	if c.Url, err = requireEnv("MEDIASERVER_PROWLARR_URL"); err != nil {
		return nil, err
	}
	if c.Username, err = requireEnv("MEDIASERVER_PROWLARR_USERNAME"); err != nil {
		return nil, err
	}
	if c.Password, err = requireEnv("MEDIASERVER_PROWLARR_PASSWORD"); err != nil {
		return nil, err
	}
	if c.QBittorrentHostname, err = requireEnv("MEDIASERVER_QBITTORRENT_HOSTNAME"); err != nil {
		return nil, err
	}
	if c.QBittorrentUsername, err = requireEnv("MEDIASERVER_QBITTORRENT_USERNAME"); err != nil {
		return nil, err
	}
	if c.QBittorrentPassword, err = requireEnv("MEDIASERVER_QBITTORRENT_PASSWORD"); err != nil {
		return nil, err
	}
	if c.FlaresolverrHostUrl, err = requireEnv("MEDIASERVER_FLARESOLVERR_HOST_URL"); err != nil {
		return nil, err
	}

	c.Indexers = []Indexer{
		{"1337x", os.Getenv("MEDIASERVER_PROWLARR_1337X_ENABLED") == "true", "1337x_indexer.json", "", ""},
		{"EZTV", os.Getenv("MEDIASERVER_PROWLARR_EZTV_ENABLED") == "true", "eztv_indexer.json", "", ""},
		{"Internet Archive", os.Getenv("MEDIASERVER_PROWLARR_INTERNETARCHIVE_ENABLED") == "true", "internetarchive_indexer.json", "", ""},
		{"LimeTorrents", os.Getenv("MEDIASERVER_PROWLARR_LIMETORRENTS_ENABLED") == "true", "limetorrents_indexer.json", "", ""},
		{"The Pirate Bay", os.Getenv("MEDIASERVER_PROWLARR_PIRATEBAY_ENABLED") == "true", "pirate_bay_indexer.json", "", ""},
		{"YTS", os.Getenv("MEDIASERVER_PROWLARR_YTS_ENABLED") == "true", "yts_indexer.json", "", ""},
		{"RuTracker.org", os.Getenv("MEDIASERVER_PROWLARR_RUTRACKER_ENABLED") == "true", "rutracker_indexer.json", os.Getenv("MEDIASERVER_PROWLARR_RUTRACKER_USERNAME"), os.Getenv("MEDIASERVER_PROWLARR_RUTRACKER_PASSWORD")},
		{"Sk-CzTorrent", os.Getenv("MEDIASERVER_PROWLARR_SKCZTORRENT_ENABLED") == "true", "skcztorrent_indexer.json", os.Getenv("MEDIASERVER_PROWLARR_SKCZTORRENT_USERNAME"), os.Getenv("MEDIASERVER_PROWLARR_SKCZTORRENT_PASSWORD")},
		{"SkTorrent.org", os.Getenv("MEDIASERVER_PROWLARR_SKTORRENT_ENABLED") == "true", "sktorrent_indexer.json", os.Getenv("MEDIASERVER_PROWLARR_SKTORRENT_USERNAME"), os.Getenv("MEDIASERVER_PROWLARR_SKTORRENT_PASSWORD")},
	}

	return c, nil
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

func (c *Config) setIndexer(indexer Indexer) error {
	log.Printf("-- Adding indexer: %v...\n", indexer.Name)

	b, err := utils.LoadJSONFile(reqBodies, indexer.File)
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	// Update indexer config
	if indexer.Username != "" && indexer.Password != "" {
		setField(b, "username", indexer.Username)
		setField(b, "password", indexer.Password)
	}

	h := map[string]string{"X-Api-Key": c.Apikey}
	_, err = utils.Request("POST", c.Url+"/api/v1/indexer", b, h, nil)
	if err != nil {
		return fmt.Errorf("failed to set indexer %v: %w", indexer.Name, err)
	}

	return nil
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

	for _, i := range c.Indexers {
		if i.Enabled {
			if err = c.setIndexer(i); err != nil {
				return err
			}
		}
	}

	log.Println("- prowlarr configured successfully!")
	fmt.Println()

	return nil
}
