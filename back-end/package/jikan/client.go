package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type JikanClient struct {
	BaseURL string
	Client  *http.Client
}

func NewJikanClient() *JikanClient {
	return &JikanClient{
		BaseURL: "https://api.jikan.moe/v4",
		Client:  &http.Client{Timeout: 20 * time.Second},
	}
}

func (client *JikanClient) FetchTopAnime(page int) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("%s/top/anime?page=%d", client.BaseURL, page)

	for attempts := 1; attempts <= 5; attempts++ {
		resp, err := client.Client.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch top anime: %v", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			_ = resp.Body.Close()
			fmt.Println("Rate limited. Retrying after 2 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
		}

		var result struct {
			Data []map[string]interface{} `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("failed to decode response: %v", err)
		}
		_ = resp.Body.Close()

		return result.Data, nil
	}

	return nil, fmt.Errorf("failed to fetch top anime after multiple attempts")
}

func SaveToFile(data []map[string]interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
