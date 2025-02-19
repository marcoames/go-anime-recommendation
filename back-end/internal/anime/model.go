package anime

// Date represents a day, month, and year.
type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

// AiredProp represents the "prop" object inside Aired.
type AiredProp struct {
	From Date `json:"from"`
	To   Date `json:"to"`
}

// Aired represents airing details (e.g., start and end dates) for an anime.
type Aired struct {
	From   string    `json:"from"`
	To     string    `json:"to"`
	String string    `json:"string"`
	Prop   AiredProp `json:"prop"`
}

// Info is a generic struct to represent companies or organizations related to the anime.
// This is used for demographics, licensors, producers, and studios.
type Info struct {
	MalID int    `json:"mal_id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// Genre represents a single genre of an anime.
type Genre struct {
	MalID int    `json:"mal_id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

// Broadcast represents broadcast information for the anime.
type Broadcast struct {
	Day      string `json:"day"`
	String   string `json:"string"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
}

// ImageDetail represents image URLs in different sizes.
type ImageDetail struct {
	ImageURL      string `json:"image_url"`
	LargeImageURL string `json:"large_image_url"`
	SmallImageURL string `json:"small_image_url"`
}

// Images groups together image details provided in both JPG and WebP formats.
type Images struct {
	Jpg  ImageDetail `json:"jpg"`
	Webp ImageDetail `json:"webp"`
}

// Trailer represents the trailer information for the anime.
type Trailer struct {
	EmbedURL  string `json:"embed_url"`
	ImageURL  string `json:"image_url"`
	URL       string `json:"url"`
	YoutubeID string `json:"youtube_id"`
}

// Anime represents the main structure of an anime entry.
type Anime struct {
	Aired         Aired     `json:"aired"`
	Background    string    `json:"background"`
	Broadcast     Broadcast `json:"broadcast"`
	Demographics  []Info    `json:"demographics"`
	Duration      string    `json:"duration"`
	Episodes      int       `json:"episodes"`
	Genres        []Genre   `json:"genres"`
	Images        Images    `json:"images"`
	Licensors     []Info    `json:"licensors"`
	MalID         int       `json:"mal_id"`
	Popularity    int       `json:"popularity"`
	Producers     []Info    `json:"producers"`
	Rank          int       `json:"rank"`
	Rating        string    `json:"rating"`
	Score         float64   `json:"score"`
	ScoredBy      int       `json:"scored_by"`
	Season        string    `json:"season"`
	Source        string    `json:"source"`
	Status        string    `json:"status"`
	Studios       []Info    `json:"studios"`
	Synopsis      string    `json:"synopsis"`
	Title         string    `json:"title"`
	TitleEnglish  string    `json:"title_english"`
	TitleJapanese string    `json:"title_japanese"`
	TitleSynonyms []string  `json:"title_synonyms"`
	Trailer       Trailer   `json:"trailer"`
	Type          string    `json:"type"`
	URL           string    `json:"url"`
	Year          int       `json:"year"`
}
