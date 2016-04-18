package main

import (
	"io"
	//"fmt"
	"github.com/boltdb/bolt"
)

// BoltDB Wrappers


// Wrapper for the BoltDB handler
type BoltDB struct {
    db *bolt.DB
}

// Creates a new BoltDB database given the *bolt.DB handler
func NewBoltDB(path string) *BoltDB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	return &BoltDB{db}
}

// Returns the BoltDB raw database handler
func (bd BoltDB) DB() *bolt.DB {
    return bd.db
}

// Return true if the key in bucket exists, false otherwise.
func (bd BoltDB) Has(bucket, key string) bool {
    exist := false
    bd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if val := b.Get([]byte(key)); val != nil {
			exist = true
		}
        return nil
	})
	return exist
}

// Returns an array of all the keys in the bucket starting from offset.
// If offset is "" then starts from the first key available.
// Limit specifies the max length of the array, if 0, then de MAX_LIMIT default will be used.
func (bd BoltDB) All(bucket, offset string, limit int) []string {
    //Negative numbers defaults to MAX_LIMIT too.
    if limit <= 0 {
        limit = LIMIT
    }
    // string array to copy all keys
	data := make([]string, 0)
	bd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
        
        // If offset exists we Seek for the key and start there, first item otherwise
        start, _ := c.First()
        if offset != "" {
            start, _ = c.Seek([]byte(offset))
        }

		for k := start; k != nil; k, _ = c.Next() {
			data = append(data, string(k))
            if len(data) == limit { break }
		}

		return nil
	})
	return data
}

// Get the value of the desired key inside the bucket.
func (bd BoltDB) Get(bucket, key string) []byte {
    var data []byte
	bd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get([]byte(key))
		if v != nil {
			data = make([]byte, len(v))
			copy(data, v)
		}
		return nil
	})
	return data
}

// Set the value of the given key inside the bucket.
func (bd BoltDB) Set(bucket, key string, value []byte) error {
    return bd.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(key), value)
	})
}

// Delete the value of the given key inside the bucket.
func (bd BoltDB) Del(bucket, key string) error {
    return bd.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Delete([]byte(key))
	})
}

// Creates Bucket if it does not exist
func (bd BoltDB) CreateBucketIfNotExist(bucket string) error {
    return bd.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
    })
}

// Returns the size of the database in bytes, used to check for
// storage limits and Content-Length header when backing up.
func (bd BoltDB) Size() int64 {
	return bd.db.View(func(tx *bolt.Tx) error {
		return tx.Size()
	})
}

// Backups the entire database to the given io.Writer() interface.
// Useful for passing an http.ResponseWriter() or an os.File()
// for backup purposes via http or timestamp.
func (bd BoltDB) Backup(w io.Writer) error {
	return bd.db.View(func(tx *bolt.Tx) error {
        _, err := tx.WriteTo(w)
        return err
    })
}