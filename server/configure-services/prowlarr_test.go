package main

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestDownloadClients(t *testing.T) {
	h, err := getProwlarrHeaders()
	if err != nil {
		t.Error(err)
	}
	respBody, err := makeRequest("GET", prowlarrBaseURL+"/api/v1/downloadclient", nil, h)
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

func TestIndexers(t *testing.T) {
	h, err := getProwlarrHeaders()
	if err != nil {
		t.Error(err)
	}
	respBody, err := makeRequest("GET", prowlarrBaseURL+"/api/v1/indexer", nil, h)
	if err != nil {
		t.Error(err)
	}

	var indexerDetails []struct {
		Name string `json:"name"`
	}
	json.Unmarshal(respBody, &indexerDetails)

	// Expected indexer names
	eIndexers := []string{
		"EZTV",
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
