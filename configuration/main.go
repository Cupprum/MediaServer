package main

import (
	"log"
	"os"

	"MediaServer/configuration/homepage"
	"MediaServer/configuration/jellyfin"
	"MediaServer/configuration/prowlarr"
	"MediaServer/configuration/qbittorrent"
)

func main() {
	if os.Getenv("QBITTORRENT_DEPLOY") == "true" {
		if err := qbittorrent.Configure(); err != nil {
			log.Fatalf("--- qbittorrent configuration failed: %v\n", err)
		}
	}

	if os.Getenv("PROWLARR_DEPLOY") == "true" {
		if err := prowlarr.Configure(); err != nil {
			log.Fatalf("--- prowlarr configuration failed: %v\n", err)
		}
	}

	if os.Getenv("JELLYFIN_DEPLOY") != "true" {
		if err := jellyfin.Configure(); err != nil {
			log.Fatalf("--- jellyfin configuration failed: %v\n", err)
		}
	}

	if os.Getenv("HOMEPAGE_DEPLOY") == "true" {
		if err := homepage.Configure(); err != nil {
			log.Fatalf("--- homepage configuration failed: %v\n", err)
		}
	}
}
