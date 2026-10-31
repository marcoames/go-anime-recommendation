package anime

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchAndSaveAnime(repo *Repository) error {
	baseURL := "https://api.jikan.moe/v4"
	var allAnime []Anime
	httpClient := &http.Client{Timeout: 20 * time.Second}

	for page := 1; page <= 200; page++ {
		percentage := (float64(page) / 200.0) * 100
		fmt.Printf("\rFetching page %d/200 (%.1f%%)", page, percentage)
		url := fmt.Sprintf("%s/top/anime?page=%d", baseURL, page)

		for attempts := 1; attempts <= 5; attempts++ {
			resp, err := httpClient.Get(url)
			if err != nil {
				return fmt.Errorf("failed to fetch page %d: %v", page, err)
			}

			if resp.StatusCode == http.StatusTooManyRequests {
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
				time.Sleep(2 * time.Second)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				_ = resp.Body.Close()
				return fmt.Errorf("API returned status code %d", resp.StatusCode)
			}

			var result struct {
				Data []Anime `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				_ = resp.Body.Close()
				return fmt.Errorf("failed to decode response: %v", err)
			}
			_ = resp.Body.Close()

			allAnime = append(allAnime, result.Data...)
			time.Sleep(1 * time.Second) // Respect Jikan's rate limit.
			goto nextPage
		}
		return fmt.Errorf("failed to fetch page %d after multiple attempts", page)
	nextPage:
	}

	return repo.SaveAnimeData(allAnime)
}
