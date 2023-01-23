package linkapibe

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var graphDataSearch = map[uint32][]uint32{
	1: {2, 3}, 2: {1, 7}, 3: {4}, 4: {1, 5}, 5: {2, 3, 8},
	6: {1, 4}, 7: {3, 6, 10}, 8: {9}, 9: {1}, 10: {8},
	75: {76}, 76: {77}, 77: {75},
}

var shortestPaths = [][]uint32{
	{1, 3, 4, 5, 8, 9},
	{1, 2, 7, 10, 8, 9},
}

func TestBfsSearch(t *testing.T) {
	assert := assert.New(t)
	// Creating temp dir for testing
	tempDir := t.TempDir()

	// Creating database
	db, err := MakeDbHandler(filepath.Join(tempDir, "dbcreate.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}
	defer db.Close()

	search := MakeSearchHandler(db)

	// Creating bucekts
	if !assert.Nil(db.CreateBuckets(), "creating buckets should work") {
		assert.FailNow("creating buckets didn't work")
	}

	// Creating links
	for k, v := range graphDataSearch {
		if err := db.AddLinks(k, v); err != nil {
			assert.FailNow("should not error adding links")
		}
	}

	// Searching shortest path
	var hops int
	res, err := search.ShortestPath(1, 9, func(u int) { hops = u })
	if err != nil {
		assert.FailNow("should not error searching")
	}
	// Comparing shortest path length
	if len(res) != hops+2 {
		assert.FailNow("path length should be same to shortest")
	}
	// Comparing to actual shortest path
	var found bool
Paths:
	for i := range shortestPaths {
		if len(shortestPaths[i]) != len(res) {
			continue
		}
		for j := range res {
			if res[j] != shortestPaths[i][j] {
				continue Paths
			}
		}
		found = true
		break
	}
	if !found {
		assert.FailNow("did not find shortest path")
	}
}

func TestBFSNotFound(t *testing.T) {
	assert := assert.New(t)
	// Creating temp dir for testing
	tempDir := t.TempDir()

	// Creating database
	db, err := MakeDbHandler(filepath.Join(tempDir, "dbcreate.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}
	defer db.Close()

	search := MakeSearchHandler(db)

	// Creating bucekts
	if !assert.Nil(db.CreateBuckets(), "creating buckets should work") {
		assert.FailNow("creating buckets didn't work")
	}

	// Creating links
	for k, v := range graphDataSearch {
		if err := db.AddLinks(k, v); err != nil {
			assert.FailNow("should not error adding links")
		}
	}

	// Searching shortest path
	res, err := search.ShortestPath(1, 77, func(u int) {})
	if err != nil {
		assert.FailNow("should not error searching")
	}
	if res != nil {
		assert.FailNow("impossible search should not give results")
	}
}
