package linkapi

import (
	"encoding/binary"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

var IdToNameBucket = []byte("id-name")
var NameToIdBucket = []byte("name-id")
var OutLinksBucket = []byte("outlinks")
var InLinksBucket = []byte("inlinks")

type DatabaseHandler struct {
	DB     *bolt.DB
	Logger *log.Logger
}

func MakeDbHandler(path string) (*DatabaseHandler, error) {
	var res = &DatabaseHandler{Logger: log.Default()}
	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	res.DB = db
	return res, nil
}

func (d *DatabaseHandler) Close() {
	d.DB.Close()
}

// CreateBuckets creates the 4 buckets
func (d *DatabaseHandler) CreateBuckets() error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(IdToNameBucket)
		if err != nil {
			return fmt.Errorf("creating bucket %s: %v", IdToNameBucket, err)
		}
		_, err = tx.CreateBucketIfNotExists(NameToIdBucket)
		if err != nil {
			return fmt.Errorf("creating bucket %s: %v", NameToIdBucket, err)
		}
		_, err = tx.CreateBucketIfNotExists(OutLinksBucket)
		if err != nil {
			return fmt.Errorf("creating bucket %s: %v", OutLinksBucket, err)
		}
		_, err = tx.CreateBucketIfNotExists(InLinksBucket)
		if err != nil {
			return fmt.Errorf("creating bucket %s: %v", InLinksBucket, err)
		}
		return nil
	})
}

// CreateArticle creates an article in the database, in id-name and name-id buckets
func (d *DatabaseHandler) CreateArticle(name string, id uint32) error {
	var id_bytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(id_bytes, id)
	return d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(IdToNameBucket)
		if err := b.Put(id_bytes, []byte(name)); err != nil {
			return err
		}
		b = tx.Bucket(NameToIdBucket)
		if err := b.Put([]byte(name), id_bytes); err != nil {
			return err
		}
		return nil
	})
}

// AddLink creates a connection in the database, in the incoming and outgoing buckets
func (d *DatabaseHandler) AddLink(src, dst uint32) error {
	var src_bytes = make([]byte, 4)
	var dst_bytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(src_bytes, src)
	binary.LittleEndian.PutUint32(dst_bytes, dst)
	if err := d.createLink(src_bytes, dst_bytes, dst, OutLinksBucket); err != nil {
		return err
	}
	if err := d.createLink(dst_bytes, src_bytes, src, InLinksBucket); err != nil {
		return err
	}
	return nil
}

// createLink makes a connection in a given bucket
func (d *DatabaseHandler) createLink(src, dst []byte, dstint uint32, bucket []byte) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(src)
		if v == nil { // If value doesn't exist
			// Create new
			if err := b.Put(src, dst); err != nil {
				return err
			}
		} else { // If value exists
			// Search if dst is inside value
			for _, id := range bytesToIds(v) {
				if id == dstint { // If dst already inside value, abort
					return nil
				}
			}
			// If value not inside, add to the end of value
			v = append(v, dst...)
			// Put updated value
			if err := b.Put(src, v); err != nil {
				return err
			}
		}
		return nil
	})
}

// bytesToIds converts a byteslice to uint32s
func bytesToIds(data []byte) []uint32 {
	var res = make([]uint32, len(data)/4)
	// Each uint32 takes 4 bytes
	for i := range res {
		res[i] = binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
	}
	return res
}
