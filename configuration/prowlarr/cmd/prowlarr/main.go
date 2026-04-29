package main

import (
	"MediaServer/configuration/prowlarr"
	"log"
	"os"
)

func main() {
	if os.Getenv("MEDIASERVER_PROWLARR_DEPLOY") == "true" {
		err := prowlarr.Configure()
		if err != nil {
			panic(err)
		}
	} else {
		log.Println("- Skipping prowlarr deployment.")
	}
}
