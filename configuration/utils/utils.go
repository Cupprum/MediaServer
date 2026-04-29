package utils

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func requestBuilder(method, url string, body interface{}, headers map[string]string) (*http.Request, error) {
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

	return req, nil
}

func Request(method, url string, body interface{}, headers map[string]string, client *http.Client, b int) ([]byte, error) {
	req, err := requestBuilder(method, url, body, headers)
	if err != nil {
		return nil, err
	}

	backoff := func(err error) {
		b--
		log.Println("--- error: backing off, going to sleep for 30s: error message:", err)
		time.Sleep(30 * time.Second)
	}

	for {
		if client == nil {
			client = &http.Client{Timeout: 60 * time.Second}
		}
		resp, err := client.Do(req)
		if err != nil {
			if b > 0 {
				backoff(err)
				continue
			}
			return nil, fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			if b > 0 {
				backoff(err)
				continue
			}
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			if b > 0 {
				backoff(err)
				continue
			}
			return nil, fmt.Errorf("request failed: %s - %s", resp.Status, string(respBody))
		}
		if os.Getenv("MEDIASERVER_LOG") == "debug" {
			log.Println("--- HTTP request completed",
				"method", method,
				"url", url,
				"status", resp.Status,
			)
		}
		return respBody, nil
	}
}

func LoadJSONFile(folder embed.FS, filename string) (map[string]interface{}, error) {
	filePath := filepath.Join("req_bodies", filename)

	data, err := folder.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from %s: %w", filename, err)
	}

	return jsonData, nil
}

func RequireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing env var: %s", key)
	}
	return val, nil
}
