package main

import (
	"testing"
)

var testdb *BoltDB
var err error

// Create a new Test database
func TestNewBoltDB(t *testing.T) {
	testdb, err = NewBoltDB("test.db")
	if err != nil {
		t.Fatalf("Could not create/open database file. %v", err)
	}
}

// Bucket creation
func TestCreateBucketIfNotExist(t *testing.T) {
	err = testdb.CreateBucketIfNotExist("test")
	if err != nil {
		t.Fatalf("Could not create bucket. %v", err)
	}
	// Should not fail if we create bucket again
	err = testdb.CreateBucketIfNotExist("test")
	if err != nil {
		t.Fatalf("Could not create bucket. %v", err)
	}
}