package main

import (
	"testing"
)

func TestQbittorrentLogin(t *testing.T) {
	c, err := getQBittorrentConfig()
	if err != nil {
		t.Error(err)
	}

	err = c.qbittorrentLogin()
	if err != nil {
		t.Error(err)
	}
}
