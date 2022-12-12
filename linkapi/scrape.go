package linkapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ScrapeHandler struct {
	DB     *DatabaseHandler
	Client *http.Client
}

func MakeScrapeHandler(db *DatabaseHandler) *ScrapeHandler {
	return &ScrapeHandler{
		DB:     db,
		Client: http.DefaultClient,
	}
}

// GetAllArticles returns the titles of all articles in a given region
func (s *ScrapeHandler) GetAllArticles(region string, callback func([]string)) error {
	var res []string
	var err error
	var cont string
	var done bool
	// While we're not done
	for !done {
		// Get the next batch
		res, cont, err = s.getAllArticlesCurrent(region, cont)
		if err != nil {
			return err
		}
		// Execute callback
		callback(res)
		// If we don't have cont, stop
		if cont == "" {
			done = true
		}
	}
	return nil
}

// getAllArticlesCurrent returns title of all articles in the current batch, and the start of next batch
func (s *ScrapeHandler) getAllArticlesCurrent(region, cont string) ([]string, string, error) {
	var res = &AllLinksStruct{}
	var titles []string
	// Making current url
	curUrl := fmt.Sprintf(AllLinksTemplate, region)
	if cont != "" {
		curUrl += "&gapcontinue=" + url.QueryEscape(cont)
	}
	// Making request
	resp, err := s.Client.Get(curUrl)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	// Unmarshaling response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, "", err
	}
	// Getting all titles
	titles = make([]string, len(res.Query.Pages))
	for i := range res.Query.Pages {
		titles[i] = res.Query.Pages[i].Title
	}
	return titles, res.Continue.Gapcontinue, nil
}
