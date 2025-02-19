package anime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func FetchAndSaveAnime(repo *Repository) error {
	baseURL := "https://api.jikan.moe/v4"
	var allAnime []Anime

	for page := 1; page <= 200; page++ {
		percentage := (float64(page) / 200.0) * 100
		fmt.Printf("\rFetching page %d/200 (%.1f%%)", page, percentage)
		url := fmt.Sprintf("%s/top/anime?page=%d", baseURL, page)

		for attempts := 1; attempts <= 5; attempts++ {
			resp, err := http.Get(url)
			if err != nil {
				return fmt.Errorf("failed to fetch page %d: %v", page, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusTooManyRequests {
				time.Sleep(2 * time.Second)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("API returned status code %d", resp.StatusCode)
			}

			var result struct {
				Data []Anime `json:"data"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to decode response: %v", err)
			}

			allAnime = append(allAnime, result.Data...)
			time.Sleep(1 * time.Second) // Rate limiting
			break
		}
	}

	return repo.SaveAnimeData(allAnime)
}
