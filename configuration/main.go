package main

import (
	"fmt"
	"os"

	"MediaServer/configuration/jellyfin"
	"MediaServer/configuration/prowlarr"
	"MediaServer/configuration/qbittorrent"
)

func main() {
	err := qbittorrent.Configure()
	if err != nil {
		fmt.Println("  * qBittorrent configuration failed:", err)
		os.Exit(1)
	}

	err = prowlarr.Configure()
	if err != nil {
		fmt.Println("  * Prowlarr configuration failed:", err)
		os.Exit(1)
	}

	err = jellyfin.Configure()
	if err != nil {
		fmt.Println("  * Jellyfin configuration failed:", err)
		os.Exit(1)
	}
}
