package main

import "MediaServer/configuration/homepage"

func main() {
	err := homepage.Configure()
	if err != nil {
		panic(err)
	}
}
