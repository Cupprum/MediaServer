package main

import (
	"log"
	"os"

	"MediaServer/server/configuration/jellyfin"
	"MediaServer/server/configuration/prowlarr"
	"MediaServer/server/configuration/qbittorrent"
)

func main() {
	if os.Getenv("MEDIASERVER_QBITTORRENT_DEPLOY") == "true" {
		if err := qbittorrent.Configure(); err != nil {
			log.Fatalf("--- qbittorrent configuration failed: %v\n", err)
		}
	}

	if os.Getenv("MEDIASERVER_PROWLARR_DEPLOY") == "true" {
		if err := prowlarr.Configure(); err != nil {
			log.Fatalf("--- prowlarr configuration failed: %v\n", err)
		}
	}

	if os.Getenv("MEDIASERVER_JELLYFIN_DEPLOY") == "true" {
		if err := jellyfin.Configure(); err != nil {
			log.Fatalf("--- jellyfin configuration failed: %v\n", err)
		}
	}
}
