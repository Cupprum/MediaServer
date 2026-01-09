package main

import (
	"log"

	"MediaServer/configuration/jellyfin"
	"MediaServer/configuration/prowlarr"
	"MediaServer/configuration/qbittorrent"
)

func main() {
	err := qbittorrent.Configure()
	if err != nil {
		log.Fatalf("--- qBittorrent configuration failed: %v\n", err)
	}

	err = prowlarr.Configure()
	if err != nil {
		log.Fatalf("--- Prowlarr configuration failed: %v\n", err)
	}

	err = jellyfin.Configure()
	if err != nil {
		log.Fatalf("--- Jellyfin configuration failed: %v\n", err)
	}
}
