package prowlarr_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"slices"
	"testing"
	"time"

	"MediaServer/server/configuration/prowlarr"
	"MediaServer/server/configuration/utils"
)

func login(c *prowlarr.Config) error {
	log.Println("-- Logging in...")

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

	// Login stores cookie in client jar, so now we can retrieve apikey
	err = c.LoadApikey(client)
	if err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	return nil
}

var configCache *prowlarr.Config

func config() (*prowlarr.Config, error) {
	if configCache != nil {
		return configCache, nil
	}

	c, err := prowlarr.GetConfig()
	if err != nil {
		return nil, err
	}

	err = login(c)
	if err != nil {
		return nil, err
	}

	// Store `config` in `cache`
	configCache = c

	return c, nil
}

func TestProwlarrLogin(t *testing.T) {
	// If we get config, login was successful
	_, err := config()
	if err != nil {
		t.Error(err)
	}
}

func TestProwlarrDownloadClients(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	h := map[string]string{"X-Api-Key": c.Apikey}
	respBody, err := utils.Request("GET", c.Url+"/api/v1/downloadclient", nil, h, nil)
	if err != nil {
		t.Error(err)
	}

	var clients []struct {
		Name string `json:"name"`
	}
	json.Unmarshal(respBody, &clients)

	for _, client := range clients {
		if client.Name == "qBittorrent" {
			return
		}
	}
	t.Error("qbittorrent download client not found in the response")
}

func TestProwlarrIndexerProxy(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	h := map[string]string{"X-Api-Key": c.Apikey}
	respBody, err := utils.Request("GET", c.Url+"/api/v1/indexerProxy", nil, h, nil)
	if err != nil {
		t.Error(err)
	}

	var proxies []struct {
		Name string `json:"name"`
	}
	json.Unmarshal(respBody, &proxies)

	for _, proxy := range proxies {
		if proxy.Name == "FlareSolverr" {
			return
		}
	}
	t.Error("flaresolverr indexer proxy not found in the response")
}

func TestProwlarrIndexers(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	h := map[string]string{"X-Api-Key": c.Apikey}
	respBody, err := utils.Request("GET", c.Url+"/api/v1/indexer", nil, h, nil)
	if err != nil {
		t.Error(err)
	}

	var indexerDetails []struct {
		Name string `json:"name"`
	}
	json.Unmarshal(respBody, &indexerDetails)

	var indexers []string
	for _, indexer := range indexerDetails {
		indexers = append(indexers, indexer.Name)
	}

	for _, i := range c.Indexers {
		if i.Enabled && !slices.Contains(indexers, i.Name) {
			t.Errorf("Expected indexer '%s' not found in actual indexers", i.Name)
		}
	}
}
