package qbittorrent_test

import (
	"encoding/json"
	"testing"

	"MediaServer/configuration/qbittorrent"
	"MediaServer/configuration/utils"
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

	rb, err := utils.Request("POST", c.Url+"/api/v2/app/preferences", nil, nil, c.Client, 1)
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

	rb, err := utils.Request("POST", c.Url+"/api/v2/app/preferences", nil, nil, c.Client, 1)
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
