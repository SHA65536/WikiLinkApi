package linkapi

import (
	"encoding/binary"
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

// Bucket names
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

// AddLinks creates multiple connections in the database, in the incoming and outgoing buckets
func (d *DatabaseHandler) AddLinks(src uint32, dst []uint32) error {
	var src_bytes = make([]byte, 4)
	var dst_bytes = idsToBytes(dst)
	binary.LittleEndian.PutUint32(src_bytes, src)
	// Adding outgoing links
	err := d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(OutLinksBucket)
		return b.Put(src_bytes, dst_bytes)
	})
	if err != nil {
		return err
	}
	// Adding incoming links
	for i := range dst {
		if err := d.createLink(dst_bytes[i*4:(i+1)*4], src_bytes, dst[i], InLinksBucket); err != nil {
			return err
		}
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

// GetAllArticles returns a map of all articles and their Id's
func (d *DatabaseHandler) GetAllArticles() (map[string]uint32, error) {
	var res = map[string]uint32{}
	return res, d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(NameToIdBucket)
		return b.ForEach(func(k, v []byte) error {
			res[string(k)] = binary.LittleEndian.Uint32(v)
			return nil
		})
	})
}

// GetName returns article name by given id
func (d *DatabaseHandler) GetName(id uint32) (string, error) {
	var res string
	var id_bytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(id_bytes, id)
	return res, d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(IdToNameBucket)
		v := b.Get(id_bytes)
		if v == nil {
			return fmt.Errorf("id %d not found", id)
		}
		res = string(v)
		return nil
	})
}

// GetId returns article id by given name
func (d *DatabaseHandler) GetId(name string) (uint32, error) {
	var res uint32
	return res, d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(NameToIdBucket)
		v := b.Get([]byte(name))
		if v == nil {
			return fmt.Errorf("name %s not found", name)
		}
		res = binary.LittleEndian.Uint32(v)
		return nil
	})
}

// GetOutgoing returns article ids by given source id
func (d *DatabaseHandler) GetOutgoing(id uint32) ([]uint32, error) {
	var res []uint32
	var id_bytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(id_bytes, id)
	return res, d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(OutLinksBucket)
		v := b.Get(id_bytes)
		if v == nil {
			return nil
		}
		res = bytesToIds(v)
		return nil
	})
}

// GetOutgoing returns article ids by given destination id
func (d *DatabaseHandler) GetIncoming(id uint32) ([]uint32, error) {
	var res []uint32
	var id_bytes = make([]byte, 4)
	binary.LittleEndian.PutUint32(id_bytes, id)
	return res, d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(InLinksBucket)
		v := b.Get(id_bytes)
		if v == nil {
			return nil
		}
		res = bytesToIds(v)
		return nil
	})
}

// IdsToNames returns a list of names of given ids
func (d *DatabaseHandler) IdsToNames(ids ...uint32) ([]string, error) {
	var res = make([]string, len(ids))
	for i := range ids {
		name, err := d.GetName(ids[i])
		if err != nil {
			return nil, err
		}
		res[i] = name
	}
	return res, nil
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

// idsToBytes converts a list of uint32 ids to a byteslice
func idsToBytes(nums []uint32) []byte {
	// Allocate a slice of bytes with the same length as the input slice
	bytes := make([]byte, len(nums)*4)

	// Iterate over the input slice of uint32 numbers
	for i, n := range nums {
		// Write the bytes of the current uint32 number in little endian order
		// to the correct offset in the slice of bytes
		bytes[i*4] = byte(n)
		bytes[i*4+1] = byte(n >> 8)
		bytes[i*4+2] = byte(n >> 16)
		bytes[i*4+3] = byte(n >> 24)
	}

	return bytes
}
