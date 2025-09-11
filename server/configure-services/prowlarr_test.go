package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

type DownloadClient struct {
	Name string `json:"name"`
}

type Indexer struct {
	Name string `json:"name"`
}

func TestMain(m *testing.M) {
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	getAPIKey()
	m.Run()
}

func TestDownloadClients(t *testing.T) {
	respBody, err := makeRequest("GET", prowlarrBaseURL+"/api/v1/downloadclient", nil, prowlarrHeaders)
	if err != nil {
		t.Fatal(err)
	}

	var clients []DownloadClient
	json.Unmarshal(respBody, &clients)

	for _, client := range clients {
		if client.Name == "qBittorrent" {
			return
		}
	}
	t.Error("qbittorrent download client not found in the response")
}

func TestIndexers(t *testing.T) {
	respBody, err := makeRequest("GET", prowlarrBaseURL+"/api/v1/indexer", nil, prowlarrHeaders)
	if err != nil {
		t.Fatal(err)
	}

	var indexerDetails []Indexer
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
		found := false
		for _, i := range indexers {
			if ei == i {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected indexer '%s' not found in actual indexers", expectedName)
		}
	}
}
