package main

import (
	"fmt"
	"os"

	"MediaServer/configuration/jellyfin"
	"MediaServer/configuration/prowlarr"
	"MediaServer/configuration/qbittorrent"
)

func main() {
	if os.Getenv("QBITTORRENT_FLAG") == "true" {
		err := qbittorrent.Configure()
		if err != nil {
			fmt.Println("  * qBittorrent configuration failed:", err)
			os.Exit(1)
		}
	}

	if os.Getenv("PROWLARR_FLAG") == "true" {
		err := prowlarr.Configure()
		if err != nil {
			fmt.Println("  * Prowlarr configuration failed:", err)
			os.Exit(1)
		}
	}

	if os.Getenv("JELLYFIN_FLAG") == "true" {
		err := jellyfin.Configure()
		if err != nil {
			fmt.Println("  * Jellyfin configuration failed:", err)
			os.Exit(1)
		}
	}
}
