package main

import (
	"MediaServer/configuration/qbittorrent"
	"log"
	"os"
)

func main() {
	if os.Getenv("MEDIASERVER_QBITTORRENT_DEPLOY") == "true" {
		err := qbittorrent.Configure()
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("- Skipping qbittorrent deployment.")
	}
}
