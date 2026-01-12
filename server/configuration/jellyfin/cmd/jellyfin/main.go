package main

import (
	"MediaServer/configuration/jellyfin"
)

func main() {
	err := jellyfin.Configure()
	if err != nil {
		panic(err)
	}
}
