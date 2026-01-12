package main

import "MediaServer/configuration/qbittorrent"

func main() {
	err := qbittorrent.Configure()
	if err != nil {
		panic(err)
	}
}
