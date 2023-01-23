package wikilinkapi

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var articlesApi = map[string]uint32{
	"One": 1, "Two": 2, "Three": 3, "Four": 4, "Five": 5,
	"Six": 6, "Seven": 7, "Eight": 8, "Nine": 9, "Ten": 10,
	"SevenFive": 75, "SevenSix": 76, "SevenSeven": 77,
}

var graphDataApi = map[uint32][]uint32{
	1: {2, 3}, 2: {1, 7}, 3: {4}, 4: {1, 5}, 5: {2, 3, 8},
	6: {1, 4}, 7: {3, 6, 10}, 8: {9}, 9: {1}, 10: {8},
	75: {76}, 76: {77}, 77: {75},
}

var shortestPathsNumApi = [][]uint32{
	{1, 3, 4, 5, 8, 9},
	{1, 2, 7, 10, 8, 9},
}

var shortestPathsTextApi = [][]string{
	{"One", "Three", "Four", "Five", "Eight", "Nine"},
	{"One", "Two", "Seven", "Ten", "Eight", "Nine"},
}

func TestApi(t *testing.T) {
	assert := assert.New(t)

	// Creating temp dir for testing
	tempDir := t.TempDir()

	// Trying to create handler when db doesn't exist
	_, err := MakeApiHandler(filepath.Join(tempDir, "api_test.db"))
	if !assert.NotNil(err) {
		assert.FailNow("Creating api handler without database should error")
	}

	// Creating database
	handler, err := MakeDbHandler(filepath.Join(tempDir, "api_test.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}

	// Creating bucekts
	if !assert.Nil(handler.CreateBuckets(), "creating buckets should work") {
		assert.FailNow("creating buckets didn't work")
	}

	// Creating articles
	for k, v := range articlesApi {
		if err := handler.CreateArticle(k, v); err != nil {
			assert.FailNow("should not error creating articles")
		}
	}

	// Adding links
	for k, v := range graphDataApi {
		if err := handler.AddLinks(k, v); err != nil {
			assert.FailNow("should not error adding links")
		}
	}
	// Closing handler
	handler.Close()

	// Making Api Handler
	api, err := MakeApiHandler(filepath.Join(tempDir, "api_test.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}

	var res SearchResult

	// Getting search with invalid params
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	w := httptest.NewRecorder()
	api.SearchRoute(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if !assert.Nil(err) {
		assert.FailNow("reading body shouldn't error")
	}
	err = json.Unmarshal(body, &res)
	if !assert.Nil(err) {
		assert.FailNow("unmarshalling response shouldn't error")
	}
	if !assert.Equal("must have 'start' and 'end' parameters!", res.Error) {
		assert.FailNow("response must be error")
	}

	res = SearchResult{}

	//Getting search with valid params
	req = httptest.NewRequest(http.MethodGet, "/search?start=One&end=Nine", nil)
	w = httptest.NewRecorder()
	api.SearchRoute(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if !assert.Nil(err) {
		assert.FailNow("reading body shouldn't error")
	}
	err = json.Unmarshal(body, &res)
	if !assert.Nil(err) {
		assert.FailNow("unmarshalling response shouldn't error")
	}
	if !assert.Equal("", res.Error) {
		assert.FailNow("response must have no error")
	}

	if !checkAgainst(res.ResultIds, shortestPathsNumApi) {
		assert.FailNow("should return shortest path ids")
	}

	if !checkAgainst(res.ResultTitles, shortestPathsTextApi) {
		assert.FailNow("should return shortest path text")
	}

	//Getting search with no path
	req = httptest.NewRequest(http.MethodGet, "/search?start=One&end=SevenSeven", nil)
	w = httptest.NewRecorder()
	api.SearchRoute(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if !assert.Nil(err) {
		assert.FailNow("reading body shouldn't error")
	}
	err = json.Unmarshal(body, &res)
	if !assert.Nil(err) {
		assert.FailNow("unmarshalling response shouldn't error")
	}
	if !assert.Equal("no path found!", res.Error) {
		assert.FailNow("response must be no path error")
	}
}

func checkAgainst[T comparable](in []T, pos [][]T) bool {
	var found bool
All:
	for _, cur := range pos {
		if len(cur) != len(in) {
			continue
		}
		for i := range in {
			if in[i] != cur[i] {
				continue All
			}
		}
		found = true
		break
	}
	return found
}
