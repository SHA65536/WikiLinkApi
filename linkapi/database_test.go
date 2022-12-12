package linkapi

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestDatabaseCreate(t *testing.T) {
	assert := assert.New(t)
	// Creating temp dir for testing
	tempDir := t.TempDir()

	// Creating database
	handler, err := MakeDbHandler(filepath.Join(tempDir, "dbcreate.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}
	defer handler.Close()

	// Creating bucekts
	if !assert.Nil(handler.CreateBuckets(), "creating buckets should work") {
		assert.FailNow("creating buckets didn't work")
	}

	// Expected bucket names
	var Expected = map[string]bool{
		string(IdToNameBucket): true,
		string(NameToIdBucket): true,
		string(OutLinksBucket): true,
		string(InLinksBucket):  true,
	}
	var Got = map[string]bool{}

	// Listing bucket names
	err = handler.DB.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			Got[string(name)] = true
			return nil
		})
		return nil
	})
	if !assert.Nil(err, "listing buckets should work") {
		assert.FailNow("listing buckets didn't work")
	}
	if !assert.Equal(Expected, Got, "4 buckets should be created") {
		assert.FailNow("creating 4 buckets didn't work")
	}
}

func TestDatabaseLinks(t *testing.T) {
	var outData = map[uint32][]uint32{
		1: {2, 3, 4},
		2: {1, 3, 4},
		3: {1},
	}
	var inData = map[uint32][]uint32{
		1: {2, 3},
		2: {1},
		3: {1, 2},
	}
	assert := assert.New(t)
	// Creating temp dir for testing
	tempDir := t.TempDir()

	// Creating database
	handler, err := MakeDbHandler(filepath.Join(tempDir, "dbcreate.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}
	defer handler.Close()

	// Creating bucekts
	if !assert.Nil(handler.CreateBuckets(), "creating buckets should work") {
		assert.FailNow("creating buckets didn't work")
	}

	// Adding links
	for k, v := range outData {
		if err := handler.AddLinks(k, v); err != nil {
			assert.FailNow("should not error adding links")
		}
	}

	// Checking outgoing links
	for k, v := range outData {
		got, err := handler.GetOutgoing(k)
		if err != nil {
			assert.FailNow("should not error getting links")
		}
		sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
		if !assert.Equal(v, got) {
			assert.FailNow("should equal original data")
		}
	}

	// Checking incoming links
	for k, v := range inData {
		got, err := handler.GetIncoming(k)
		if err != nil {
			assert.FailNow("should not error getting links")
		}
		sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
		if !assert.Equal(v, got) {
			assert.FailNow("should equal original data")
		}
	}
}

func TestDatabaseArticles(t *testing.T) {
	var arts = map[string]uint32{
		"Henlo": 1, "Does": 2, "This": 3, "Work?": 4, "טקסט": 5,
	}

	assert := assert.New(t)
	// Creating temp dir for testing
	tempDir := t.TempDir()

	// Creating database
	handler, err := MakeDbHandler(filepath.Join(tempDir, "dbcreate.db"))
	if !assert.Nil(err, "handler creation should work") {
		assert.FailNow("handler creation didn't work")
	}
	defer handler.Close()

	// Creating bucekts
	if !assert.Nil(handler.CreateBuckets(), "creating buckets should work") {
		assert.FailNow("creating buckets didn't work")
	}

	// Creating articles
	for k, v := range arts {
		if err := handler.CreateArticle(k, v); err != nil {
			assert.FailNow("should not error creating articles")
		}
	}

	// Getting ids and names
	for k, v := range arts {
		// Getting id
		gotId, err := handler.GetId(k)
		if err != nil {
			assert.FailNow("should not getting ids")
		}
		if !assert.Equal(v, gotId) {
			assert.FailNow("should equal original id")
		}
		// Getting name
		gotName, err := handler.GetName(v)
		if err != nil {
			assert.FailNow("should not getting names")
		}
		if !assert.Equal(k, gotName) {
			assert.FailNow("should equal original name")
		}
	}
}

func TestBytesToIds(t *testing.T) {
	assert := assert.New(t)
	var nums = []uint32{1, 2, 3, 4, 5, 1337, ^uint32(0)}

	// Encoding numbers to binary
	var bytes = idsToBytes(nums)

	// Decoding numbers from binary
	for i, num := range bytesToIds(bytes) {
		if !assert.Equal(nums[i], num, "decoding should be same as encoding") {
			assert.FailNow("decoding should be same as encoding")
		}
	}
}

// TODO: benchmark :)
func BenchmarkBytesToIds(b *testing.B) {

}
