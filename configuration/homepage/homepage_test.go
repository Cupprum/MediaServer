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

	// TODO: change to check for Grafana, Qbittorrent and other services
	if strings.Contains(string(rb), "qBittorrent") == false {
		t.Error("qBittorrent service widget not found")
	}
}
