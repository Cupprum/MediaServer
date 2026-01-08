package prowlarr_test

import (
	"encoding/json"
	"slices"
	"testing"

	"MediaServer/configuration/prowlarr"
	"MediaServer/configuration/utils"
)

var cc *prowlarr.ProwlarrConfig

func config() (*prowlarr.ProwlarrConfig, error) {
	if cc != nil {
		return cc, nil
	}

	c, err := prowlarr.Config()
	if err != nil {
		return nil, err
	}

	err = c.Login()
	if err != nil {
		return nil, err
	}

	// Stare `config` in `cache`
	cc = c

	return c, nil
}

func TestProwlarrLogin(t *testing.T) {
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

	// Expected indexer names
	eIndexers := []string{
		// "EZTV",
		"LimeTorrents",
		"The Pirate Bay",
		"YTS",
	}

	var indexers []string
	for _, indexer := range indexerDetails {
		indexers = append(indexers, indexer.Name)
	}

	for _, ei := range eIndexers {
		if !slices.Contains(indexers, ei) {
			t.Errorf("Expected indexer '%s' not found in actual indexers", ei)
		}
	}
}
