package main

import (
	"MediaServer/configuration/prowlarr"
)

func main() {
	err := prowlarr.Configure()
	if err != nil {
		panic(err)
	}
}
