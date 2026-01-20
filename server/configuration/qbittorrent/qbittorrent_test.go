package qbittorrent_test

import (
	"testing"

	"MediaServer/server/configuration/qbittorrent"
)

func TestQbittorrentLogin(t *testing.T) {
	c, err := qbittorrent.GetConfig()
	if err != nil {
		t.Error(err)
	}

	err = c.Login()
	if err != nil {
		t.Error(err)
	}
}

// TODO: write test to check if seeding limits were set correctly
