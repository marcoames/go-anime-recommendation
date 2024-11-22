package main

import "strings"

// Genre struct represents a single genre of an anime
type Genre struct {
	MalID int    `json:"mal_id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// Broadcast struct represents broadcast information for the anime
type Broadcast struct {
	Day      string `json:"day"`
	String   string `json:"string"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
}

// Aired struct represents airing details (e.g., start and end dates)
type Aired struct {
	From   string `json:"from"`
	To     string `json:"to"`
	String string `json:"string"`
	Prop   struct {
		From struct {
			Day   int `json:"day"`
			Month int `json:"month"`
			Year  int `json:"year"`
		} `json:"from"`
		To struct {
			Day   int `json:"day"`
			Month int `json:"month"`
			Year  int `json:"year"`
		} `json:"to"`
	} `json:"prop"`
}

// Anime struct represents the main structure of an anime entry
type Anime struct {
	Aired        Aired     `json:"aired"`
	Background   string    `json:"background"`
	Broadcast    Broadcast `json:"broadcast"`
	Demographics []struct {
		MalID int    `json:"mal_id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"demographics"`
	Duration string  `json:"duration"`
	Episodes int     `json:"episodes"`
	Genres   []Genre `json:"genres"` // Genres as an array of Genre structs
	Images   struct {
		Jpg struct {
			ImageURL      string `json:"image_url"`
			LargeImageURL string `json:"large_image_url"`
			SmallImageURL string `json:"small_image_url"`
		} `json:"jpg"`
		Webp struct {
			ImageURL      string `json:"image_url"`
			LargeImageURL string `json:"large_image_url"`
			SmallImageURL string `json:"small_image_url"`
		} `json:"webp"`
	} `json:"images"`
	Licensors []struct {
		MalID int    `json:"mal_id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"licensors"`
	MalID      int `json:"mal_id"`
	Popularity int `json:"popularity"`
	Producers  []struct {
		MalID int    `json:"mal_id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"producers"`
	Rank     int     `json:"rank"`
	Rating   string  `json:"rating"`
	Score    float64 `json:"score"`
	ScoredBy int     `json:"scored_by"`
	Season   string  `json:"season"`
	Source   string  `json:"source"`
	Status   string  `json:"status"`
	Studios  []struct {
		MalID int    `json:"mal_id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		URL   string `json:"url"`
	} `json:"studios"`
	Synopsis      string   `json:"synopsis"`
	Title         string   `json:"title"`
	TitleEnglish  string   `json:"title_english"`
	TitleJapanese string   `json:"title_japanese"`
	TitleSynonyms []string `json:"title_synonyms"`
	Trailer       struct {
		EmbedURL  string `json:"embed_url"`
		ImageURL  string `json:"image_url"`
		URL       string `json:"url"`
		YoutubeID string `json:"youtube_id"`
	} `json:"trailer"`
	Type string `json:"type"`
	URL  string `json:"url"`
	Year int    `json:"year"`
}

func getAnimeIndex(animeName string, animeData []Anime) int {
	for i, anime := range animeData {
		// Compare anime names, case insensitive
		if strings.EqualFold(anime.Title, animeName) {
			return i // Return the index of the found anime
		}
	}
	// Return -1 if no anime with the specified name is found
	return -1
}
