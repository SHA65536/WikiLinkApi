package ui

import "time"

// Search Result from Wikipedia
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

// Search Response from the search route
type SearchResponse struct {
	Error  string          `json:"error,omitempty"`
	Result []SearchArticle `json:"result,omitempty"`
}

type SearchArticle struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Pageid  int    `json:"pageid"`
}

// Result Response from result route
type ResultResponse struct {
	Error        string   `json:"error,omitempty"`
	ResultIds    []uint32 `json:"ids,omitempty"`
	ResultTitles []string `json:"titles,omitempty"`
}
