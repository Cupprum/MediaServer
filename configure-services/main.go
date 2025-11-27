package main

import (
	"fmt"
	"os"
)

func main() {
	if os.Getenv("QBITTORRENT_FLAG") == "true" {
		err := ConfigureQBittorrent()
		if err != nil {
			fmt.Println("  * qBittorrent configuration failed:", err)
			os.Exit(1)
		}
	}

	if os.Getenv("PROWLARR_FLAG") == "true" {
		err := ConfigureProwlarr()
		if err != nil {
			fmt.Println("  * Prowlarr configuration failed:", err)
			os.Exit(1)
		}
	}

	if os.Getenv("JELLYFIN_FLAG") == "true" {
		err := ConfigureJellyfin()
		if err != nil {
			fmt.Println("  * Jellyfin configuration failed:", err)
			os.Exit(1)
		}
	}
}
