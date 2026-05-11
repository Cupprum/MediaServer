package main

import (
	"fmt"
	"log"
	"os"

	"MediaServer/configuration/jellyfin"
	"MediaServer/configuration/prowlarr"
	"MediaServer/configuration/qbittorrent"
)

func main() {
	fmt.Print(os.Getenv("MEDIASERVER_QBITTORRENT_USERNAME"))
	fmt.Print(os.Getenv("MEDIASERVER_QBITTORRENT_RAW_PASSWORD"))
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
