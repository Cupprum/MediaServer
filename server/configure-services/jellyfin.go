package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const baseURL = "http://jellyfin.pi.local"

var logger *slog.Logger

type StartupConfig struct {
	UICulture                 string `json:"UICulture"`
	MetadataCountryCode       string `json:"MetadataCountryCode"`
	PreferredMetadataLanguage string `json:"PreferredMetadataLanguage"`
}

type User struct {
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type PathInfo struct {
	Path string `json:"Path"`
}

type TypeOption struct {
	Type                 string   `json:"Type"`
	MetadataFetchers     []string `json:"MetadataFetchers"`
	MetadataFetcherOrder []string `json:"MetadataFetcherOrder"`
	ImageFetchers        []string `json:"ImageFetchers"`
	ImageFetcherOrder    []string `json:"ImageFetcherOrder"`
}

type LibraryOptions struct {
	TypeOptions []TypeOption `json:"TypeOptions"`
	PathInfos   []PathInfo   `json:"PathInfos"`
}

type VirtualFolder struct {
	LibraryOptions LibraryOptions `json:"LibraryOptions"`
}

type RemoteAccess struct {
	EnableRemoteAccess         bool `json:"EnableRemoteAccess"`
	EnableAutomaticPortMapping bool `json:"EnableAutomaticPortMapping"`
}

func makeRequest(method, url string, body interface{}) error {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s - %s", resp.Status, string(respBody))
	}

	logger.Info("HTTP request completed",
		"method", method,
		"url", url,
		"status", resp.Status,
	)

	return nil
}

func checkSystemInfo() error {
	logger.Info("Checking system info...")
	return makeRequest("GET", baseURL+"/System/Info", nil)
}

func configureStartup() error {
	logger.Info("Configuring startup settings...")
	config := StartupConfig{
		UICulture:                 "en-US",
		MetadataCountryCode:       "US",
		PreferredMetadataLanguage: "en",
	}
	return makeRequest("POST", baseURL+"/Startup/Configuration", config)
}

func checkUser() error {
	logger.Info("Checking user status...")
	return makeRequest("GET", baseURL+"/Startup/User", nil)
}

func createUser() error {
	logger.Info("Creating admin user...")

	username := os.Getenv("JELLYFIN_USERNAME")
	if username == "" {
		return fmt.Errorf("Missing env var: `JELLYFIN_USERNAME`")
	}

	password := os.Getenv("JELLYFIN_PASSWORD")
	if password == "" {
		return fmt.Errorf("Missing env var: `JELLYFIN_PASSWORD`")
	}

	user := User{
		Name:     username,
		Password: password,
	}
	return makeRequest("POST", baseURL+"/Startup/User", user)
}

func createMoviesLibrary() error {
	logger.Info("Creating Movies library...")

	moviesLibrary := VirtualFolder{
		LibraryOptions: LibraryOptions{
			TypeOptions: []TypeOption{
				{
					Type: "Movie",
					MetadataFetchers: []string{
						"TheMovieDb",
						"The Open Movie Database",
					},
					MetadataFetcherOrder: []string{
						"TheMovieDb",
						"The Open Movie Database",
					},
					ImageFetchers: []string{
						"TheMovieDb",
						"The Open Movie Database",
						"Embedded Image Extractor",
						"Screen Grabber",
					},
					ImageFetcherOrder: []string{
						"TheMovieDb",
						"The Open Movie Database",
						"Embedded Image Extractor",
						"Screen Grabber",
					},
				},
			},
			PathInfos: []PathInfo{
				{Path: "/media/movies"},
			},
		},
	}

	return makeRequest("POST", baseURL+"/Library/VirtualFolders?collectionType=movies&refreshLibrary=false&name=Movies", moviesLibrary)
}

func createTVShowsLibrary() error {
	logger.Info("Creating TV Shows library...")

	tvLibrary := VirtualFolder{
		LibraryOptions: LibraryOptions{
			TypeOptions: []TypeOption{
				{
					Type: "Series",
					MetadataFetchers: []string{
						"TheMovieDb",
						"The Open Movie Database",
					},
					MetadataFetcherOrder: []string{
						"TheMovieDb",
						"The Open Movie Database",
					},
					ImageFetchers: []string{
						"TheMovieDb",
					},
					ImageFetcherOrder: []string{
						"TheMovieDb",
					},
				},
				{
					Type: "Season",
					MetadataFetchers: []string{
						"TheMovieDb",
					},
					MetadataFetcherOrder: []string{
						"TheMovieDb",
					},
					ImageFetchers: []string{
						"TheMovieDb",
					},
					ImageFetcherOrder: []string{
						"TheMovieDb",
					},
				},
				{
					Type: "Episode",
					MetadataFetchers: []string{
						"TheMovieDb",
						"The Open Movie Database",
					},
					MetadataFetcherOrder: []string{
						"TheMovieDb",
						"The Open Movie Database",
					},
					ImageFetchers: []string{
						"TheMovieDb",
						"The Open Movie Database",
						"Embedded Image Extractor",
						"Screen Grabber",
					},
					ImageFetcherOrder: []string{
						"TheMovieDb",
						"The Open Movie Database",
						"Embedded Image Extractor",
						"Screen Grabber",
					},
				},
			},
			PathInfos: []PathInfo{
				{Path: "/media/tv"},
			},
		},
	}

	return makeRequest("POST", baseURL+"/Library/VirtualFolders?collectionType=tvshows&refreshLibrary=false&name=Shows", tvLibrary)
}

func configureRemoteAccess() error {
	logger.Info("Configuring remote access...")
	remoteAccess := RemoteAccess{
		EnableRemoteAccess:         false,
		EnableAutomaticPortMapping: false,
	}
	return makeRequest("POST", baseURL+"/Startup/RemoteAccess", remoteAccess)
}

func completeStartup() error {
	logger.Info("Completing startup...")
	return makeRequest("POST", baseURL+"/Startup/Complete", nil)
}

func ConfigureJellyfin() error {
	logger.Info("Starting Jellyfin configuration...")

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Check System Info", checkSystemInfo},
		{"Configure Startup", configureStartup},
		{"Check User", checkUser},
		{"Create User", createUser},
		{"Create Movies Library", createMoviesLibrary},
		{"Create TV Shows Library", createTVShowsLibrary},
		{"Configure Remote Access", configureRemoteAccess},
		{"Complete Startup", completeStartup},
	}

	for _, step := range steps {
		logger.Info("Executing step", "step", step.name)
		if err := step.fn(); err != nil {
			logger.Error("Step failed", "step", step.name, "error", err)
			return err
		}
		logger.Info("Step completed successfully", "step", step.name)
	}

	logger.Info("Jellyfin configuration completed!")
	return nil
}

func setJellyfinLogger(l *slog.Logger) {
	logger = l
}
