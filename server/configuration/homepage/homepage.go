package homepage

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

//go:embed services.yaml.tpl
var servicesTemplate string

func Configure() error {
	log.Println("- Starting homepage configuration...")

	configFolder := os.Getenv("MEDIASERVER_CONFIG_DIR")
	if configFolder == "" {
		return fmt.Errorf("MEDIASERVER_CONFIG_DIR environment variable is not set")
	}

	log.Println("-- Templating the services.yaml file...")
	services := os.ExpandEnv(string(servicesTemplate))

	// Make sure the config folder exists
	folder := filepath.Join(configFolder, "homepage", "config")
	err := os.MkdirAll(folder, 0755)
	if err != nil {
		return fmt.Errorf("failed to create homepage config folder: %w", err)
	}

	path := filepath.Join(folder, "services.yaml")
	log.Printf("-- Saving services.yaml file to: %s...", path)
	err = os.WriteFile(path, []byte(services), 0644)
	if err != nil {
		return fmt.Errorf("failed to write services.yaml file into homepage config folder: %w", err)
	}

	log.Println("- homepage configured successfully!")
	fmt.Println()
	return nil
}
