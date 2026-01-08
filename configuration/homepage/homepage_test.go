package homepagetest

import (
	"os"
	"strings"
	"testing"

	"MediaServer/configuration/utils"
)

func TestHomepageSite(t *testing.T) {
	u := os.Getenv("HOMEPAGE_URL")
	if u == "" {
		t.Error("HOMEPAGE_URL environment variable not set")
	}

	rb, err := utils.Request("GET", u, nil, nil, nil)
	if err != nil {
		t.Error(err)
	}

	if !strings.Contains(string(rb), "qBittorrent") {
		t.Error("qbittorrent service widget not found")
	}

	if !strings.Contains(string(rb), "Prowlarr") {
		t.Error("prowlarr service widget not found")
	}

	if !strings.Contains(string(rb), "Jellyfin") {
		t.Error("jellyfin service widget not found")
	}
}
