package jellyfin

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"MediaServer/configuration/utils"
)

//go:embed req_bodies/*.json
var reqBodies embed.FS

type Config struct {
	Url      string
	Username string
	Password string
}

func config() (*Config, error) {
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

	return &Config{
		Url:      url,
		Username: username,
		Password: password,
	}, nil
}

// Used to cache the Jellyfin authorization token
var jellyfinAuthorization string = ""

func Headers() (map[string]string, error) {
	if jellyfinAuthorization != "" {
		return map[string]string{"Authorization": jellyfinAuthorization}, nil
	}

	fmt.Println("-- Getting Jellyfin authorization token...")
	c, err := config()
	if err != nil {
		return nil, err
	}

	// Initial auth header without token -> More details https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	defaultAuth := `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="1", Version="10.11.0"`

	accessToken, err := c.jellyfinLogin()
	if err != nil {
		return nil, err
	}

	// Add token to the initial auth header
	jellyfinAuthorization = fmt.Sprintf(`%s, Token="%s"`, defaultAuth, accessToken)
	return map[string]string{"Authorization": jellyfinAuthorization}, nil
}

// TODO: change naming to get rid of Jellyfin
func (c *Config) jellyfinLogin() (string, error) {
	b := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{c.Username, c.Password}

	// Initial auth header without token -> More details https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	defaultAuth := `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="1", Version="10.11.0"`
	h := map[string]string{"Authorization": defaultAuth}

	rb, err := utils.Request("POST", c.Url+"/Users/AuthenticateByName", b, h, nil)
	if err != nil {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			return "", fmt.Errorf("not logged in")
		}
		return "", fmt.Errorf("failed to log in: %w", err)
	}

	var r struct {
		AccessToken string `json:"AccessToken"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %v", err)
	}

	return r.AccessToken, nil
}

func (c *Config) checkJellyfinSystemInfo() error {
	fmt.Println("-- Checking system info...")
	_, err := utils.Request("GET", c.Url+"/System/Info", nil, nil, nil)
	return err
}

func (c *Config) configureJellyfinStartup() error {
	fmt.Println("-- Configuring startup settings...")

	url := c.Url + "/Startup/Configuration"

	rconfig, err := utils.Request("GET", url, nil, nil, nil)
	if err != nil {
		return err
	}

	config := map[string]interface{}{}
	if err := json.Unmarshal(rconfig, &config); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	_, err = utils.Request("POST", url, config, nil, nil)
	return err
}

func (c *Config) checkJellyfinUser() error {
	fmt.Println("-- Checking user status...")
	_, err := utils.Request("GET", c.Url+"/Startup/User", nil, nil, nil)
	return err
}

func (c *Config) createJellyfinUser() error {
	fmt.Println("-- Creating admin user...")

	b := struct {
		Name     string `json:"Name"`
		Password string `json:"Password"`
	}{
		Name:     c.Username,
		Password: c.Password,
	}
	_, err := utils.Request("POST", c.Url+"/Startup/User", b, nil, nil)
	return err
}

func (c *Config) createJellyfinMoviesLibrary() error {
	fmt.Println("-- Creating Movies library...")

	b, err := utils.LoadJSONFile(reqBodies, "library_movies.json")
	if err != nil {
		return err
	}

	_, err = utils.Request("POST", c.Url+"/Library/VirtualFolders?collectionType=movies&refreshLibrary=false&name=Movies", b, nil, nil)
	return err
}

func (c *Config) createJellyfinTVShowsLibrary() error {
	fmt.Println("-- Creating TV Shows library...")

	b, err := utils.LoadJSONFile(reqBodies, "library_tv.json")
	if err != nil {
		return err
	}

	_, err = utils.Request("POST", c.Url+"/Library/VirtualFolders?collectionType=tvshows&refreshLibrary=false&name=Shows", b, nil, nil)
	return err
}

func (c *Config) configureJellyfinRemoteAccess() error {
	fmt.Println("-- Configuring remote access...")

	// Too small to store this req body as a file
	b := struct {
		ERA bool `json:"EnableRemoteAccess"`
	}{false}

	_, err := utils.Request("POST", c.Url+"/Startup/RemoteAccess", b, nil, nil)
	return err
}

func (c *Config) completeJellyfinStartup() error {
	fmt.Println("-- Completing startup...")
	_, err := utils.Request("POST", c.Url+"/Startup/Complete", nil, nil, nil)
	return err
}

func Configure() error {
	fmt.Println("- Starting jellyfin configuration...")

	c, err := config()
	if err != nil {
		return err
	}

	// Try to login
	_, err = c.jellyfinLogin()
	if err == nil {
		fmt.Println("- already configured, skipping...")
		fmt.Println()
		return nil
	} else if err.Error() != "not logged in" {
		return fmt.Errorf("failed to login: %w", err)
	}
	// If error is "not logged in", proceed with configuration

	// Otherwise, proceed with configuration
	if err = c.checkJellyfinSystemInfo(); err != nil {
		return err
	}
	if err = c.configureJellyfinStartup(); err != nil {
		return err
	}
	if err = c.checkJellyfinUser(); err != nil {
		return err
	}
	if err = c.createJellyfinUser(); err != nil {
		return err
	}
	if err = c.createJellyfinMoviesLibrary(); err != nil {
		return err
	}
	if err = c.createJellyfinTVShowsLibrary(); err != nil {
		return err
	}
	if err = c.configureJellyfinRemoteAccess(); err != nil {
		return err
	}
	if err = c.completeJellyfinStartup(); err != nil {
		return err
	}

	fmt.Println("- jellyfin configured successfully!")
	fmt.Println()
	return nil
}
