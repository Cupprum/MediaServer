package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

type QBittorrentConfig struct {
	Url      string
	Username string
	Password string
	Client   *http.Client
}

func getQBittorrentConfig() (*QBittorrentConfig, error) {
	fmt.Println("-- Loading qBittorrent config...")

	url := os.Getenv("QBITTORRENT_URL")
	if url == "" {
		return nil, fmt.Errorf("missing env var: `QBITTORRENT_URL`")
	}

	username := os.Getenv("QBITTORRENT_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("missing env var: `QBITTORRENT_USERNAME`")
	}

	password := os.Getenv("QBITTORRENT_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("missing env var: `QBITTORRENT_PASSWORD`")
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

	return &QBittorrentConfig{
		Url:      url,
		Username: username,
		Password: password,
		Client:   client,
	}, nil
}

func (c *QBittorrentConfig) qbittorrentLogin() error {
	fmt.Println("-- Logging in to qBittorrent...")

	b := fmt.Sprintf("username=%s&password=%s", c.Username, c.Password)

	r, err := Request("POST", c.Url+"/api/v2/auth/login", b, nil, c.Client)
	if err != nil {
		return err
	}

	if string(r) != "Ok." {
		return fmt.Errorf("login failed, unexpected response: %s", string(r))
	}
	return nil
}

func getQbittorrentPasswordFromLogs() (string, error) {
	fmt.Println("-- Getting initial qBittorrent password...")

	cmd := exec.Command("docker", "ps", "-a", "--filter", "ancestor=qbittorrent", "-q")
	o, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to filter for qbittorrent container id: %v", err)
	}

	containerId := strings.TrimSpace(string(o))

	cmd = exec.Command("bash", "-c",
		"docker logs "+containerId+" | grep 'temporary password' | awk '{print $NF}'")
	o, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v", err)
	}

	pw := strings.TrimSpace(string(o))
	if pw == "" {
		return "", fmt.Errorf("failed to retrieve temporary password from logs")
	}

	return pw, nil
}

func (c *QBittorrentConfig) changePassword() error {
	fmt.Println("-- Changing qBittorrent password...")

	b := struct {
		Username string `json:"web_ui_username"`
		Password string `json:"web_ui_password"`
	}{
		Username: c.Username,
		Password: c.Password,
	}

	// Convert map to JSON string for form encoding
	jsonBytes, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to marshal preferences to JSON: %w", err)
	}

	// Send as form-encoded data with json parameter
	formData := fmt.Sprintf("json=%s", url.QueryEscape(string(jsonBytes)))

	_, err = Request("POST", c.Url+"/api/v2/app/setPreferences", formData, nil, c.Client)
	if err != nil {
		return err
	}

	return nil
}

func ConfigureQBittorrent() error {
	fmt.Println("- Starting qBittorrent configuration...")
	c, err := getQBittorrentConfig()
	if err != nil {
		return err
	}

	// Try to login
	if err = c.qbittorrentLogin(); err == nil {
		// If login is successful, assume already configured
		fmt.Println("- qBittorrent already configured, skipping...")
		fmt.Println()
		return nil
	}

	// Otherwise, proceed with configuration
	pw := c.Password

	// On failure, get temp password from logs
	tempPw, err := getQbittorrentPasswordFromLogs()
	if err != nil {
		return err
	}
	c.Password = tempPw

	// Retry login with temp password
	if err = c.qbittorrentLogin(); err != nil {
		return err
	}

	// Set password back to original value
	c.Password = pw

	// Change password to desired value
	if err = c.changePassword(); err != nil {
		return err
	}

	// Retry login with original password
	if err = c.qbittorrentLogin(); err != nil {
		return err
	}

	fmt.Println("- qBittorrent configured successfully!")
	fmt.Println()
	return nil
}
