package jellyfin

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"MediaServer/server/configuration/utils"
)

//go:embed req_bodies/*.json
var reqBodies embed.FS

type Config struct {
	AccessToken string // Set during login
	Url         string
	Username    string
	Password    string
}

func GetConfig() (*Config, error) {
	log.Println("-- Creating config based on Environment Variables...")

	url := os.Getenv("MEDIASERVER_JELLYFIN_URL")
	if url == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_JELLYFIN_URL`")
	}

	username := os.Getenv("MEDIASERVER_JELLYFIN_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_JELLYFIN_USERNAME`")
	}

	password := os.Getenv("MEDIASERVER_JELLYFIN_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_JELLYFIN_PASSWORD`")
	}

	// Initial AccessToken without token value -> More details https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	return &Config{
		AccessToken: `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="1", Version="10.11.0"`,
		Url:         url,
		Username:    username,
		Password:    password,
	}, nil
}

// Used to check if jellyfin was already configured and also in tests
func (c *Config) LoadAccessToken() error {
	b := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{c.Username, c.Password}

	h := map[string]string{"Authorization": c.AccessToken}

	rb, err := utils.Request("POST", c.Url+"/Users/AuthenticateByName", b, h, nil)
	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			// During first call we dont have a valid token yet
			// as our user was not created yet.
			return fmt.Errorf("not logged in")
		}
		return fmt.Errorf("failed to log in: %w", err)
	}

	var r struct {
		AccessToken string `json:"AccessToken"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return fmt.Errorf("failed to parse auth response: %v", err)
	}

	// Add retrieved token to `AccessToken`
	if !strings.Contains(c.AccessToken, "Token") {
		c.AccessToken = fmt.Sprintf(`%s, Token="%s"`, c.AccessToken, r.AccessToken)
	}

	return nil
}

func (c *Config) checkSystemInfo() error {
	log.Println("-- Checking system info...")
	_, err := utils.Request("GET", c.Url+"/System/Info", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}
	return nil
}

func (c *Config) configureStartup() error {
	log.Println("-- Configuring startup settings...")

	url := c.Url + "/Startup/Configuration"

	rconfig, err := utils.Request("GET", url, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get startup configuration: %w", err)
	}

	config := map[string]interface{}{}
	if err := json.Unmarshal(rconfig, &config); err != nil {
		return fmt.Errorf("failed to decode startup configuration response: %w", err)
	}

	_, err = utils.Request("POST", url, config, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to set startup configuration: %w", err)
	}
	return nil
}

func (c *Config) checkUser() error {
	log.Println("-- Checking user status...")
	_, err := utils.Request("GET", c.Url+"/Startup/User", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get startup user: %w", err)
	}
	return nil
}

func (c *Config) createUser() error {
	log.Println("-- Creating user...")

	b := struct {
		Name     string `json:"Name"`
		Password string `json:"Password"`
	}{
		Name:     c.Username,
		Password: c.Password,
	}
	_, err := utils.Request("POST", c.Url+"/Startup/User", b, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (c *Config) createMoviesLibrary() error {
	log.Println("-- Creating Movies library...")

	b, err := utils.LoadJSONFile(reqBodies, "library_movies.json")
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	_, err = utils.Request("POST", c.Url+"/Library/VirtualFolders?collectionType=movies&refreshLibrary=false&name=Movies", b, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create movies library: %w", err)
	}
	return nil
}

func (c *Config) createTVShowsLibrary() error {
	log.Println("-- Creating TV Shows library...")

	b, err := utils.LoadJSONFile(reqBodies, "library_tv.json")
	if err != nil {
		return fmt.Errorf("failed to retrieve json payload: %w", err)
	}

	_, err = utils.Request("POST", c.Url+"/Library/VirtualFolders?collectionType=tvshows&refreshLibrary=false&name=Shows", b, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create tv shows library: %w", err)
	}
	return nil
}

func (c *Config) configureRemoteAccess() error {
	log.Println("-- Configuring remote access...")

	// Too small to store this req body as a file
	b := struct {
		ERA bool `json:"EnableRemoteAccess"`
	}{false}

	_, err := utils.Request("POST", c.Url+"/Startup/RemoteAccess", b, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to configure remote access: %w", err)
	}
	return nil
}

func (c *Config) completeStartup() error {
	log.Println("-- Completing startup...")
	_, err := utils.Request("POST", c.Url+"/Startup/Complete", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to complete startup: %w", err)
	}
	return nil
}

func Configure() error {
	log.Println("- Starting jellyfin configuration...")

	c, err := GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Try to login
	err = c.LoadAccessToken()
	if err == nil {
		log.Println("- already configured, skipping...")
		fmt.Println()
		return nil
	} else if err.Error() != "not logged in" {
		return fmt.Errorf("failed to login: %w", err)
	}
	// If error is "not logged in", proceed with configuration

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

	log.Println("- jellyfin configured successfully!")
	fmt.Println()
	return nil
}
