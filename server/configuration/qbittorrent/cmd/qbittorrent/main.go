package main

import "MediaServer/server/configuration/qbittorrent"

func main() {
	err := qbittorrent.Configure()
	if err != nil {
		panic(err)
	}
}
