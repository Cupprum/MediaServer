package main

import "MediaServer/server/configuration/homepage"

func main() {
	err := homepage.Configure()
	if err != nil {
		panic(err)
	}
}
