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
	fmt.Println(os.Getenv("MEDIASERVER_QBITTORRENT_USERNAME"))
	fmt.Println(os.Getenv("MEDIASERVER_QBITTORRENT_RAW_PASSWORD"))
	fmt.Println(os.Environ())
	fmt.Println()
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
