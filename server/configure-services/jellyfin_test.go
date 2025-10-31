package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"
)

func getJellyfinAuthorization() (string, error) {
	logger.Info("Getting Jellyfin authorization...")

	u := os.Getenv("JELLYFIN_USERNAME")
	if u == "" {
		return "", fmt.Errorf("JELLYFIN_USERNAME environment variable not set")
	}

	pw := os.Getenv("JELLYFIN_PASSWORD")
	if pw == "" {
		return "", fmt.Errorf("JELLYFIN_PASSWORD environment variable not set")
	}

	b := struct {
		Username string `json:"Username"`
		Pw       string `json:"Pw"`
	}{u, pw}

	defaultAuth := `MediaBrowser Client="Jellyfin", Device="TestScript", DeviceId="12345", Version="10.8.0"`
	h := map[string]string{
		"Authorization": defaultAuth,
	}

	rb, err := makeRequest("POST", jellyfinBaseURL+"/Users/AuthenticateByName", b, h)
	if err != nil {
		return "", err
	}

	var r struct {
		AccessToken string `json:"AccessToken"`
	}
	if err := json.Unmarshal(rb, &r); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %v", err)
	}

	// TODO: can i just extend the previous header?
	return fmt.Sprintf(`MediaBrowser Token="%s", Client="Jellyfin", Device="TestScript", DeviceId="12345", Version="10.8.0"`, r.AccessToken), nil
}

func getJellyfinItems(path string) ([]string, error) {
	auth, err := getJellyfinAuthorization()
	if err != nil {
		return nil, err
	}
	h := map[string]string{
		"Authorization": auth,
	}

	rb, err := makeRequest("GET", jellyfinBaseURL+path, nil, h)
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
