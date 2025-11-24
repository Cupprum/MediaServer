package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Request(method, url string, body interface{}, headers map[string]string, client *http.Client) ([]byte, error) {
	var reqBody io.Reader

	if body != nil {
		if _, ok := body.(string); ok {
			// If body is string, use it directly (Form data)
			reqBody = bytes.NewBufferString(body.(string))
		} else {
			// Otherwise marshal as JSON
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal JSON: %w", err)
			}
			reqBody = bytes.NewBuffer(jsonData)
		}
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		if _, ok := body.(string); ok {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Set additional headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed: %s - %s", resp.Status, string(respBody))
	}

	fmt.Println(" * HTTP request completed",
		"method", method,
		"url", url,
		"status", resp.Status,
	)

	return respBody, nil
}

func loadJSONFile(service string, filename string) (map[string]interface{}, error) {
	filePath := filepath.Join("req_bodies", service, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filename, err)
	}

	return jsonData, nil
}
