package main

import (
	"MediaServer/server/configuration/prowlarr"
)

func main() {
	err := prowlarr.Configure()
	if err != nil {
		panic(err)
	}
}
