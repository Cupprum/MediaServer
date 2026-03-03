package qbittorrent_test

import (
	"encoding/json"
	"testing"

	"MediaServer/server/configuration/qbittorrent"
	"MediaServer/server/configuration/utils"
)

func config() (*qbittorrent.Config, error) {
	c, err := qbittorrent.GetConfig()
	if err != nil {
		return nil, err
	}

	err = c.Login()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func TestLogin(t *testing.T) {
	_, err := config()
	if err != nil {
		t.Error(err)
	}
}

func TestSeedingLimits(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	rb, err := utils.Request("POST", c.Url+"/api/v2/app/preferences", nil, nil, c.Client)
	if err != nil {
		t.Error(err)
	}

	var sl qbittorrent.SeedingLimits
	json.Unmarshal(rb, &sl)

	if !sl.RatioEnabled {
		t.Error("expected ratio enabled to be true")
	}
	if sl.RatioLimit != 1.0 {
		t.Errorf("expected ratio limit to be 1.0, got %v", sl.RatioLimit)
	}
	if !sl.TimeEnabled {
		t.Error("expected time enabled to be true")
	}
	if sl.TimeLimit != 60 {
		t.Errorf("expected time limit to be 60, got %v", sl.TimeLimit)
	}
	if !sl.InactiveEnabled {
		t.Error("expected inactive enabled to be true")
	}
	if sl.InactiveLimit != 60 {
		t.Errorf("expected inactive limit to be 60, got %v", sl.InactiveLimit)
	}
	if sl.Action != 1 {
		t.Errorf("expected action to be 1, got %v", sl.Action)
	}
}

func TestManagementMode(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	rb, err := utils.Request("POST", c.Url+"/api/v2/app/preferences", nil, nil, c.Client)
	if err != nil {
		t.Error(err)
	}

	var mm qbittorrent.ManagementMode
	json.Unmarshal(rb, &mm)

	if !mm.AutoMode {
		t.Error("expected Automatic Torrent Management Mode to be enabled")
	}
	if !mm.ChangePathOnCategory {
		t.Error("expected updating path on category change to be enabled")
	}
}

func TestCategories(t *testing.T) {
	c, err := config()
	if err != nil {
		t.Error(err)
	}

	rb, err := utils.Request("GET", c.Url+"/api/v2/sync/maindata", nil, nil, c.Client)
	if err != nil {
		t.Error(err)
	}

	type Category struct {
		Name     string `json:"name"`
		SavePath string `json:"savePath"`
	}
	type resp struct {
		Categories map[string]Category `json:"categories"`
	}

	var r resp
	json.Unmarshal(rb, &r)

	m, mok := r.Categories["Movies"]
	if !mok {
		t.Errorf("expected 'Movies' category")
	}
	tv, tok := r.Categories["TV"]
	if !tok {
		t.Errorf("expected 'TV' category")
	}

	if m.Name != "Movies" {
		t.Errorf("expected movies category to have name 'Movies'")
	}
	if m.SavePath != "/downloads/movies" {
		t.Errorf("expected movies category to have savePath '/downloads/movies'")
	}

	if tv.Name != "TV" {
		t.Errorf("expected tv category to have name 'TV'")
	}
	if tv.SavePath != "/downloads/tv" {
		t.Errorf("expected tv category to have savePath '/downloads/tv'")
	}
}
