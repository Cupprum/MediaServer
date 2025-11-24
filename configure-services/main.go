package main

import (
	"fmt"
	"os"
)

func main() {
	err := ConfigureQBittorrent()
	if err != nil {
		fmt.Println("  * qBittorrent configuration failed:", err)
		os.Exit(1)
	}

	err = ConfigureProwlarr()
	if err != nil {
		fmt.Println("  * Prowlarr configuration failed:", err)
		os.Exit(1)
	}

	err = ConfigureJellyfin()
	if err != nil {
		fmt.Println("  * Jellyfin configuration failed:", err)
		os.Exit(1)
	}
}
