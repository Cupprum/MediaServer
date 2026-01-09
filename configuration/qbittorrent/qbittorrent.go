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

	"MediaServer/configuration/utils"
)

type Config struct {
	Url      string
	Username string
	Password string
	Client   *http.Client
}

func GetConfig() (*Config, error) {
	log.Println("-- Loading qBittorrent config...")

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

	fmt.Println("------------test")
	fmt.Println(string(o))

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

func (c *Config) changePassword() error {
	log.Println("-- Changing password...")

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

	_, err = utils.Request("POST", c.Url+"/api/v2/app/setPreferences", formData, nil, c.Client)
	return err
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

	log.Println("- qbittorrent configured successfully!")
	fmt.Println()
	return nil
}
