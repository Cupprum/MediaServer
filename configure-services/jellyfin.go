package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type JellyfinConfig struct {
	Url      string
	Username string
	Password string
}

func getJellyfinConfig() (*JellyfinConfig, error) {
	fmt.Println("-- Create Jellyfin config based on Environment Variables...")

	url := os.Getenv("JELLYFIN_URL")
	if url == "" {
		return nil, fmt.Errorf("missing env var: `JELLYFIN_URL`")
	}

	username := os.Getenv("JELLYFIN_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing env var: `JELLYFIN_USERNAME`")
	}

	password := os.Getenv("JELLYFIN_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("missing env var: `JELLYFIN_PASSWORD`")
	}

	return &JellyfinConfig{
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}

func (c *JellyfinConfig) jellyfinLogin() (string, error) {
	b := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{c.Username, c.Password}

	// Initial auth header without token -> More details https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	defaultAuth := `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="1", Version="10.11.0"`
	h := map[string]string{"Authorization": defaultAuth}

	rb, err := Request("POST", c.Url+"/Users/AuthenticateByName", b, h, nil)
	if err != nil {
		return "", err
	}

	var r struct {
		AccessToken string `json:"AccessToken"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %v", err)
	}

	return r.AccessToken, nil
}

func (c *JellyfinConfig) checkSystemInfo() error {
	fmt.Println("-- Checking system info...")
	_, err := Request("GET", c.Url+"/System/Info", nil, nil, nil)
	return err
}

func (c *JellyfinConfig) configureStartup() error {
	fmt.Println("-- Configuring startup settings...")

	url := c.Url + "/Startup/Configuration"

	rconfig, err := Request("GET", url, nil, nil, nil)
	if err != nil {
		return err
	}

	config := map[string]interface{}{}
	if err := json.Unmarshal(rconfig, &config); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	_, err = Request("POST", url, config, nil, nil)
	return err
}

func (c *JellyfinConfig) checkUser() error {
	fmt.Println("-- Checking user status...")
	_, err := Request("GET", c.Url+"/Startup/User", nil, nil, nil)
	return err
}

func (c *JellyfinConfig) createUser() error {
	fmt.Println("-- Creating admin user...")

	b := struct {
		Name     string `json:"Name"`
		Password string `json:"Password"`
	}{
		Name:     c.Username,
		Password: c.Password,
	}
	_, err := Request("POST", c.Url+"/Startup/User", b, nil, nil)
	return err
}

func (c *JellyfinConfig) createMoviesLibrary() error {
	fmt.Println("-- Creating Movies library...")

	b, err := loadJSONFile("jellyfin", "library_movies.json")
	if err != nil {
		return err
	}

	_, err = Request("POST", c.Url+"/Library/VirtualFolders?collectionType=movies&refreshLibrary=false&name=Movies", b, nil, nil)
	return err
}

func (c *JellyfinConfig) createTVShowsLibrary() error {
	fmt.Println("-- Creating TV Shows library...")

	b, err := loadJSONFile("jellyfin", "library_tv.json")
	if err != nil {
		return err
	}

	_, err = Request("POST", c.Url+"/Library/VirtualFolders?collectionType=tvshows&refreshLibrary=false&name=Shows", b, nil, nil)
	return err
}

func (c *JellyfinConfig) configureRemoteAccess() error {
	fmt.Println("-- Configuring remote access...")

	// Too small to store this req body as a file
	b := struct {
		ERA bool `json:"EnableRemoteAccess"`
	}{false}

	_, err := Request("POST", c.Url+"/Startup/RemoteAccess", b, nil, nil)
	return err
}

func (c *JellyfinConfig) completeStartup() error {
	fmt.Println("-- Completing startup...")
	_, err := Request("POST", c.Url+"/Startup/Complete", nil, nil, nil)
	return err
}

func ConfigureJellyfin() error {
	fmt.Println("- Starting Jellyfin configuration...")

	c, err := getJellyfinConfig()
	if err != nil {
		return err
	}

	// Check if Jellyfin was already configured
	if _, err := c.jellyfinLogin(); err == nil {
		// If login is successful, assume already configured
		fmt.Println("- Jellyfin already configured, skipping...")
		return nil
	}

	// Otherwise, proceed with configuration
	if err = c.checkSystemInfo(); err != nil {
		return err
	}
	if err = c.configureStartup(); err != nil {
		return err
	}
	if err = c.checkUser(); err != nil {
		return err
	}
	if err = c.createUser(); err != nil {
		return err
	}
	if err = c.createMoviesLibrary(); err != nil {
		return err
	}
	if err = c.createTVShowsLibrary(); err != nil {
		return err
	}
	if err = c.configureRemoteAccess(); err != nil {
		return err
	}
	if err = c.completeStartup(); err != nil {
		return err
	}

	fmt.Println("- Jellyfin configured successfully!")
	return nil
}
