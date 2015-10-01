package main

import (
	//"fmt"
	"github.com/boltdb/bolt"
)

// #DATABASE
func FindKey(db *bolt.DB, bucket, key []byte) (exist bool) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if val := b.Get(key); val != nil {
			exist = true
		} else {
			exist = false
		}
		return nil
	})
	return exist
}

func GetValue(db *bolt.DB, bucket, key []byte) []byte {
	var data []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(key)
		if v != nil {
			data = make([]byte, len(v))
			copy(data, v)
		}
		return nil
	})
	return data
}

func SetValue(db *bolt.DB, bucket, key, value []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Put(key, value)
		return err
	})
}

func DeleteValue(db *bolt.DB, bucket, key []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Delete(key)
		return err
	})
}
