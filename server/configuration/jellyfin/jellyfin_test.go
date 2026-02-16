package jellyfin_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"MediaServer/server/configuration/jellyfin"
	"MediaServer/server/configuration/utils"
)

var configCache *jellyfin.Config

func config() (*jellyfin.Config, error) {
	if configCache != nil {
		return configCache, nil
	}

	c, err := jellyfin.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	err = c.LoadAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	configCache = c

	return c, nil
}

func getJellyfinItems(path string) ([]string, error) {
	c, err := config()
	if err != nil {
		return nil, err
	}

	h := map[string]string{"Authorization": c.AccessToken}

	rb, err := utils.Request("GET", c.Url+path, nil, h, nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		Items []struct {
			Name string `json:"Name"`
		} `json:"Items"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return nil, err
	}

	var items []string
	for _, item := range r.Items {
		items = append(items, item.Name)
	}

	return items, nil
}

func TestJellyfinShouldContainMoviesLibrary(t *testing.T) {
	items, err := getJellyfinItems("/Library/MediaFolders")
	if err != nil {
		t.Error(err)
	}

	if !slices.Contains(items, "Movies") {
		t.Error("Movies library not found in Jellyfin")
	}
}

func TestJellyfinLibraryShouldContainMovies(t *testing.T) {
	items, err := getJellyfinItems("/Items?IncludeItemTypes=Movie&Recursive=true")
	if err != nil {
		t.Error(err)
	}

	if len(items) == 0 {
		t.Error("No movies found in Jellyfin library")
	}
}

func TestJellyfinShouldContainSeriesLibrary(t *testing.T) {
	items, err := getJellyfinItems("/Library/MediaFolders")
	if err != nil {
		t.Error(err)
	}

	// Series are apparently called "Shows" in Jellyfin
	if !slices.Contains(items, "Shows") {
		t.Error("Shows library not found in Jellyfin")
	}
}

func TestJellyfinLibraryShouldContainSeries(t *testing.T) {
	items, err := getJellyfinItems("/Items?IncludeItemTypes=Series&Recursive=true")
	if err != nil {
		t.Error(err)
	}

	if len(items) == 0 {
		t.Error("No series found in Jellyfin library")
	}
}

func TestJellyfinOpensubtitlesShouldBeInstalled(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	status, err := c.GetAppStatus("Open Subtitles")
	if err != nil {
		t.Error(err)
	}

	if status != "Active" {
		t.Error("OpenSubtitles app is not active")
	}
}
