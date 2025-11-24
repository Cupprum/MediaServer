package main

import (
	"fmt"
	"os"
)

func main() {
	err := ConfigureJellyfin()
	if err != nil {
		fmt.Println("- Jellyfin configuration failed:", err)
		os.Exit(1)
	}
}
