package jellyfin_test

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"

	"MediaServer/configuration/jellyfin"
	"MediaServer/configuration/utils"
)

// TODO: try to cleanup working with headers and tokens

// Used to cache the Jellyfin authorization token
var token string = ""

func headers() (map[string]string, error) {
	if token != "" {
		return map[string]string{"Authorization": token}, nil
	}

	fmt.Println("-- Getting authorization token...")
	c, err := jellyfin.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Initial auth header without token -> More details https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	defaultAuth := `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="1", Version="10.11.0"`

	accessToken, err := c.Login()
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	// Add token to the initial auth header
	token = fmt.Sprintf(`%s, Token="%s"`, defaultAuth, accessToken)
	return map[string]string{"Authorization": token}, nil
}

func getJellyfinItems(path string) ([]string, error) {
	h, err := headers()
	if err != nil {
		return nil, err
	}

	// TODO: if i kind of have config already, why not use it instead of getting env var?
	url := os.Getenv("JELLYFIN_URL")
	if url == "" {
		return nil, fmt.Errorf("JELLYFIN_URL environment variable not set")
	}

	rb, err := utils.Request("GET", url+path, nil, h, nil)
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
