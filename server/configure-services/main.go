package main

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func main() {
	// Initialize structured logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	logger.Info("Starting service configuration...")

	logger.Info("Configuring Jellyfin service...")
	if err := ConfigureJellyfin(); err != nil {
		logger.Error("Failed to configure Jellyfin", "error", err)
		os.Exit(1)
	}

	logger.Info("Configuring Prowlarr service...")
	if err := ConfigureProwlarr(); err != nil {
		logger.Error("Failed to configure Prowlarr", "error", err)
		os.Exit(1)
	}

	// TODO: Add other services here
	// Example:
	// logger.Info("Configuring qBittorrent service...")
	// if err := ConfigureQBittorrent(); err != nil {
	//     logger.Error("Failed to configure qBittorrent", "error", err)
	//     os.Exit(1)
	// }

	// logger.Info("Configuring Heimdall service...")
	// if err := ConfigureHeimdall(); err != nil {
	//     logger.Error("Failed to configure Heimdall", "error", err)
	//     os.Exit(1)
	// }

	logger.Info("All services configured successfully!")
}
