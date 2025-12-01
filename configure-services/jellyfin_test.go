package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"testing"
)

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
