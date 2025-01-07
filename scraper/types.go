package scraper

import "errors"

// Struct Definitions
type Manga struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string 		`json:"desc"`
	Link        string      `json:"link"`
	Thumbnail   string      `json:"thumb"`
	ReleaseDate string      `json:"release"`
	Meta        []Metadata  `json:"meta"`
}

type Metadata struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type LatestManga struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Link      string `json:"link"`
	Thumbnail string `json:"thumb"`
}

type Chapter struct {
	ID        string   `json:"id"`
	Number    int      `json:"number"`
	Title     string   `json:"title"`
	PanelURLs []string `json:"panels"`
}

type SearchResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumb"`
}

var ErrNotFound = errors.New("resource not found")