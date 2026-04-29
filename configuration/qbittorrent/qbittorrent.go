package qbittorrent

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"MediaServer/configuration/utils"
)

type Config struct {
	Url      string
	Username string
	Password string
	Client   *http.Client
}

func GetConfig() (*Config, error) {
	log.Println("-- Loading config...")

	var err error
	c := &Config{}

	if c.Url, err = utils.RequireEnv("MEDIASERVER_QBITTORRENT_URL"); err != nil {
		return nil, err
	}
	if c.Username, err = utils.RequireEnv("MEDIASERVER_QBITTORRENT_USERNAME"); err != nil {
		return nil, err
	}
	if c.Password, err = utils.RequireEnv("MEDIASERVER_QBITTORRENT_PASSWORD"); err != nil {
		return nil, err
	}

	// Create a cookie jar for persisting cookies across requests
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	c.Client = &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return c, nil
}

func (c *Config) Login() error {
	log.Println("-- Logging in to qBittorrent...")

	b := fmt.Sprintf("username=%s&password=%s", c.Username, c.Password)

	r, err := utils.Request("POST", c.Url+"/api/v2/auth/login", b, nil, c.Client, 6)
	if err != nil {
		return err
	}

	// If response does not contain "Ok.", login failed
	if string(r) != "Ok." {
		return fmt.Errorf("not logged in")
	}
	return nil
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

	_, err = utils.Request("POST", c.Url+"/api/v2/app/setPreferences", formData, nil, c.Client, 6)
	return err
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

type ManagementMode struct {
	AutoMode             bool `json:"auto_tmm_enabled"`
	ChangePathOnCategory bool `json:"category_changed_tmm_enabled"` // Change download path on Category change
}

// Specify download path based on Category
func (c *Config) setupManagementMode() error {
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

	err = c.Login()
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	if err = c.setSeedingLimits(); err != nil {
		return err
	}
	if err = c.setupManagementMode(); err != nil {
		return err
	}

	log.Println("- qbittorrent configured successfully!")
	fmt.Println()
	return nil
}
