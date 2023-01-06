package ui

import "time"

type SearchResult struct {
	Batchcomplete string         `json:"batchcomplete"`
	Continue      SearchContinue `json:"continue"`
	Query         SearchQuery    `json:"query"`
}
type SearchContinue struct {
	Sroffset int    `json:"sroffset"`
	Continue string `json:"continue"`
}
type Searchinfo struct {
	Totalhits int `json:"totalhits"`
}
type Search struct {
	Ns        int       `json:"ns"`
	Title     string    `json:"title"`
	Pageid    int       `json:"pageid"`
	Size      int       `json:"size"`
	Wordcount int       `json:"wordcount"`
	Snippet   string    `json:"snippet"`
	Timestamp time.Time `json:"timestamp"`
}
type SearchQuery struct {
	Searchinfo Searchinfo `json:"searchinfo"`
	Search     []Search   `json:"search"`
}
