package main

import (
	"encoding/json"
	"strings"
	"testing"
)

type DownloadClient struct {
	Name string `json:"name"`
}

func TestQbittorrentPresent(t *testing.T) {
	_, err := getAPIKey()
	if err != nil {
		t.Fatal(err)
	}

	respBody, err := makeRequest("GET", prowlarrBaseURL+"/api/v1/downloadclient", nil, headers)
	if err != nil {
		t.Fatal(err)
	}

	var clients []DownloadClient
	json.Unmarshal(respBody, &clients)

	for _, client := range clients {
		if strings.Contains(strings.ToLower(client.Name), "qbittorrent") {
			return // Test passes
		}
	}
	t.Error("qbittorrent download client not found in the response")
}
