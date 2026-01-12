package main

import (
	"MediaServer/server/configuration/jellyfin"
)

func main() {
	err := jellyfin.Configure()
	if err != nil {
		panic(err)
	}
}
