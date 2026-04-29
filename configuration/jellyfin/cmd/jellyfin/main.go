package main

import (
	"MediaServer/configuration/jellyfin"
	"log"
	"os"
)

func main() {
	if os.Getenv("MEDIASERVER_JELLYFIN_DEPLOY") == "true" {
		err := jellyfin.Configure()
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("- Skipping jellyfin deployment.")
	}
}
