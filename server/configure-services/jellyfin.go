package main

import (
	"fmt"
	"os"
)

const jellyfinBaseURL = "http://jellyfin.pi.local"

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

func checkSystemInfo() error {
	logger.Info("Checking system info...")
	_, err := makeRequest("GET", jellyfinBaseURL+"/System/Info", nil, nil)
	return err
}

func configureStartup() error {
	logger.Info("Configuring startup settings...")
	config := StartupConfig{
		UICulture:                 "en-US",
		MetadataCountryCode:       "US",
		PreferredMetadataLanguage: "en",
	}
	_, err := makeRequest("POST", jellyfinBaseURL+"/Startup/Configuration", config, nil)
	return err
}

func checkUser() error {
	logger.Info("Checking user status...")
	_, err := makeRequest("GET", jellyfinBaseURL+"/Startup/User", nil, nil)
	return err
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
	_, err := makeRequest("POST", jellyfinBaseURL+"/Startup/User", user, nil)
	return err
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

	_, err := makeRequest("POST", jellyfinBaseURL+"/Library/VirtualFolders?collectionType=movies&refreshLibrary=false&name=Movies", moviesLibrary, nil)
	return err
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

	_, err := makeRequest("POST", jellyfinBaseURL+"/Library/VirtualFolders?collectionType=tvshows&refreshLibrary=false&name=Shows", tvLibrary, nil)
	return err
}

func configureRemoteAccess() error {
	logger.Info("Configuring remote access...")
	remoteAccess := RemoteAccess{
		EnableRemoteAccess:         false,
		EnableAutomaticPortMapping: false,
	}
	_, err := makeRequest("POST", jellyfinBaseURL+"/Startup/RemoteAccess", remoteAccess, nil)
	return err
}

func completeStartup() error {
	logger.Info("Completing startup...")
	_, err := makeRequest("POST", jellyfinBaseURL+"/Startup/Complete", nil, nil)
	return err
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
