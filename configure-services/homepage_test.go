package main

import (
	"os"
	"strings"
	"testing"
)

func TestHomepageSite(t *testing.T) {
	u := os.Getenv("HOMEPAGE_URL")
	if u == "" {
		t.Error("HOMEPAGE_URL environment variable not set")
	}

	rb, err := Request("GET", u, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}

	// TODO: change to check for Grafana, Qbittorrent and other services
	if strings.Contains(string(rb), "My Second Service") == false {
		t.Error("Homepage content not found")
	}
}
