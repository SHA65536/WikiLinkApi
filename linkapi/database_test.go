package linkapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestDatabaseCreate(t *testing.T) {
	assert := assert.New(t)
	// Creating temp dir for testing
	tempDir, err := os.MkdirTemp("./", "testing")
	if !assert.Nil(err, "temp dir creation should work") {
		assert.FailNow("temp dir creation didn't work")
	}
	defer os.RemoveAll(tempDir)

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

// TODO: implements
func TestDatabaseLinks(t *testing.T) {

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
