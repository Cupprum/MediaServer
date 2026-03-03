package qbittorrent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"MediaServer/server/configuration/utils"
)

type Config struct {
	Url      string
	Username string
	Password string
	Client   *http.Client
}

func GetConfig() (*Config, error) {
	log.Println("-- Loading qBittorrent config...")

	url := os.Getenv("MEDIASERVER_QBITTORRENT_URL")
	if url == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_QBITTORRENT_URL`")
	}

	username := os.Getenv("MEDIASERVER_QBITTORRENT_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_QBITTORRENT_USERNAME`")
	}

	password := os.Getenv("MEDIASERVER_QBITTORRENT_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("missing env var: `MEDIASERVER_QBITTORRENT_PASSWORD`")
	}

	// Create a cookie jar for persisting cookies across requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &Config{
		Url:      url,
		Username: username,
		Password: password,
		Client:   client,
	}, nil
}

func (c *Config) Login() error {
	log.Println("-- Logging in to qBittorrent...")

	b := fmt.Sprintf("username=%s&password=%s", c.Username, c.Password)

	r, err := utils.Request("POST", c.Url+"/api/v2/auth/login", b, nil, c.Client)
	if err != nil {
		return err
	}

	// If response does not contain "Ok.", login failed
	if string(r) != "Ok." {
		return fmt.Errorf("not logged in")
	}
	return nil
}

func dockerLogs() (string, error) {
	cmd := exec.Command("docker", "ps", "-a", "--filter", "ancestor=qbittorrent", "-q")
	o, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to filter for qbittorrent container id: %v", err)
	}

	containerId := strings.TrimSpace(string(o))

	cmd = exec.Command("bash", "-c", "docker logs "+containerId+" | grep 'temporary password' | awk '{print $NF}'")
	o, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v", err)
	}

	pw := strings.TrimSpace(string(o))
	if pw == "" {
		return "", fmt.Errorf("failed to retrieve temporary password from docker logs")
	}
	return pw, nil
}

func kubectlLogs() (string, error) {
	cmd := exec.Command("bash", "-c", "kubectl logs -l app=qbittorrent -n server | grep 'temporary password' | awk '{print $NF}'")
	o, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v", err)
	}

	pw := strings.TrimSpace(string(o))
	if pw == "" {
		return "", fmt.Errorf("failed to retrieve temporary password from kubectl logs")
	}
	return pw, nil
}

func getPasswordFromLogs() (string, error) {
	log.Println("-- Getting initial password...")

	pw, err := dockerLogs()
	if err != nil {
		log.Println("-- Docker logs method failed, trying kubectl...")
		pw, err = kubectlLogs()
		if err != nil {
			return "", fmt.Errorf("failed to get temporary password from logs: %w", err)
		}
	}

	return pw, nil
}

func (c *Config) setPreferences(b any) error {
	log.Println("-- Setting preferences...")

	// Convert map to JSON string for form encoding
	jsonBytes, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}

	// Send as form-encoded data with json parameter
	formData := fmt.Sprintf("json=%s", url.QueryEscape(string(jsonBytes)))

	_, err = utils.Request("POST", c.Url+"/api/v2/app/setPreferences", formData, nil, c.Client)
	return err
}

func (c *Config) changePassword() error {
	log.Println("-- Changing password...")

	b := struct {
		Username string `json:"web_ui_username"`
		Password string `json:"web_ui_password"`
	}{
		Username: c.Username,
		Password: c.Password,
	}

	err := c.setPreferences(b)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	return nil
}

func (c *Config) configureUser() error {
	// Get temp password from logs
	tpw, err := getPasswordFromLogs()
	if err != nil {
		return err
	}

	// Backup original password
	opw := c.Password

	// Set temp password for login
	c.Password = tpw

	// Retry login with temp password
	if err = c.Login(); err != nil {
		return err
	}

	// Set password back to original value
	c.Password = opw

	// Change password to desired value
	if err = c.changePassword(); err != nil {
		return err
	}

	return nil
}

type SeedingLimits struct {
	RatioEnabled    bool    `json:"max_ratio_enabled"`
	RatioLimit      float64 `json:"max_ratio"`
	TimeEnabled     bool    `json:"max_seeding_time_enabled"`
	TimeLimit       int     `json:"max_seeding_time"`
	InactiveEnabled bool    `json:"max_inactive_seeding_time_enabled"`
	InactiveLimit   int     `json:"max_inactive_seeding_time"`
	Action          int     `json:"max_ratio_act"`
}

func (c *Config) setSeedingLimits() error {
	log.Println("-- Configuring Seeding Limits...")

	b := SeedingLimits{
		RatioEnabled:    true,
		RatioLimit:      1.0,
		TimeEnabled:     true,
		TimeLimit:       60, // 60 minutes
		InactiveEnabled: true,
		InactiveLimit:   60, // 60 minutes
		Action:          1,  // 1 = Remove Torrent
	}

	err := c.setPreferences(b)
	if err != nil {
		return fmt.Errorf("failed to set seeding limits: %w", err)
	}

	return nil
}

func (c *Config) createCategory(category, savePath string) error {
	log.Printf("-- Creating category: %v...\n", category)

	b := fmt.Sprintf("category=%s&savePath=%s", category, savePath)

	_, err := utils.Request("POST", c.Url+"/api/v2/torrents/createCategory", b, nil, c.Client)
	if err != nil {
		return fmt.Errorf("failed to create category %v: %w", category, err)
	}

	return nil
}

type ManagementMode struct {
	AutoMode             bool `json:"auto_tmm_enabled"`
	ChangePathOnCategory bool `json:"category_changed_tmm_enabled"` // Change download path on Category change
}

func (c *Config) setupCategories() error {
	log.Println("-- Configuring categories...")
	if err := c.createCategory("Movies", "/downloads/movies"); err != nil {
		return err
	}
	if err := c.createCategory("TV", "/downloads/tv"); err != nil {
		return err
	}

	log.Println("-- Configuring Torrent Management Mode...")

	b := ManagementMode{
		AutoMode:             true,
		ChangePathOnCategory: true,
	}

	err := c.setPreferences(b)
	if err != nil {
		return fmt.Errorf("failed to set management mode: %w", err)
	}

	return nil
}

func Configure() error {
	log.Println("- Starting qbittorrent configuration...")
	c, err := GetConfig()
	if err != nil {
		return err
	}

	// Try to login
	err = c.Login()
	if err == nil {
		log.Println("- already configured, skipping...")
		fmt.Println()
		return nil
	} else if err.Error() != "not logged in" {
		return fmt.Errorf("failed to login: %w", err)
	}
	// If error is "not logged in", proceed with configuration
	// as we need to configure user

	if err = c.configureUser(); err != nil {
		return fmt.Errorf("failed to configure user: %w", err)
	}

	// Set system preferences
	if err = c.setSeedingLimits(); err != nil {
		return err
	}
	if err = c.setupCategories(); err != nil {
		return err
	}

	log.Println("- qbittorrent configured successfully!")
	fmt.Println()
	return nil
}
