package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

type DownloadClient struct {
	Name string `json:"name"`
}

func TestQbittorrentPresent(t *testing.T) {
	apiKey, err := getAPIKey()
	if err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest("GET", "http://prowlarr.pi.local/api/v1/downloadclient", nil)
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var clients []DownloadClient
	json.NewDecoder(resp.Body).Decode(&clients)

	for _, client := range clients {
		if strings.Contains(strings.ToLower(client.Name), "qbittorrent") {
			return // Test passes
		}
	}
	t.Error("qbittorrent download client not found in the response")
}
