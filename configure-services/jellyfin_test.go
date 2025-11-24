package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"
)

// Used to cache the Jellyfin authorization token
var jellyfinAuthorization string = ""

func getJellyfinHeaders() (map[string]string, error) {
	if jellyfinAuthorization != "" {
		return map[string]string{"Authorization": jellyfinAuthorization}, nil
	}

	fmt.Println("-- Getting Jellyfin authorization token...")
	c, err := getJellyfinConfig()
	if err != nil {
		return nil, err
	}

	b := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{c.Username, c.Password}

	// Initial auth header without token -> More details https://gist.github.com/nielsvanvelzen/ea047d9028f676185832e51ffaf12a6f
	defaultAuth := `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="1", Version="10.11.0"`
	h := map[string]string{"Authorization": defaultAuth}

	rb, err := Request("POST", c.Url+"/Users/AuthenticateByName", b, h, nil)
	if err != nil {
		return nil, err
	}

	var r struct {
		AccessToken string `json:"AccessToken"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %v", err)
	}

	// Add token to the initial auth header
	jellyfinAuthorization = fmt.Sprintf(`%s, Token="%s"`, defaultAuth, r.AccessToken)
	return map[string]string{"Authorization": jellyfinAuthorization}, nil
}

func getJellyfinItems(path string) ([]string, error) {
	h, err := getJellyfinHeaders()
	if err != nil {
		return nil, err
	}

	url := os.Getenv("JELLYFIN_URL")
	if url == "" {
		return nil, fmt.Errorf("JELLYFIN_URL environment variable not set")
	}

	rb, err := Request("GET", url+path, nil, h, nil)
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

func getJellyfinMediaFolders() ([]string, error) {
	return getJellyfinItems("/Library/MediaFolders")
}

func TestJellyfinShouldContainMoviesLibrary(t *testing.T) {
	items, err := getJellyfinMediaFolders()
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
	items, err := getJellyfinMediaFolders()
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
