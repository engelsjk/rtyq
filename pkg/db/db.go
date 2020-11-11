package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tidwall/buntdb"
)

// Create ...
func Create(path string) error {

	fn := filepath.Base(path)

	var (
		ErrDatabaseFileAlreadyExists error = fmt.Errorf("database file (%s) already exists", filepath.Base(path))
		ErrDatabaseFileCreate        error = fmt.Errorf("unable to create database file (%s)", fn)
	)

	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		return ErrDatabaseFileAlreadyExists
	}

	_, err = buntdb.Open(path)
	if err != nil {
		return ErrDatabaseFileCreate
	}

	return nil
}

// Initialize ...
func Initialize(path, index string, skipIndex bool) (*buntdb.DB, error) {

	fn := filepath.Base(path)

	var (
		ErrDatabaseFileDoesNotExist error = fmt.Errorf("database file (%s) does not exist", fn)
		ErrDatabaseFileIsADir       error = fmt.Errorf("database file (%s) is a dir, needs to be a file", fn)
		ErrDatabaseFileOpen         error = fmt.Errorf("unable to open database file (%s)", fn)
		ErrSpatialIndexCreate       error = fmt.Errorf("unable to spatially index database file (%s)", fn)
	)

	fmt.Printf("initializing %s (index '%s')\n", filepath.Base(path), index)

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, ErrDatabaseFileDoesNotExist
	}

	if info.IsDir() {
		return nil, ErrDatabaseFileIsADir
	}

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, ErrDatabaseFileOpen
	}

	if skipIndex {
		return db, nil
	}

	name := index
	pattern := fmt.Sprintf("%s:*", index)

	fmt.Printf("running spatial index...\n")

	start := time.Now()

	err = db.CreateSpatialIndex(name, pattern, buntdb.IndexRect)
	if err != nil {
		return nil, ErrSpatialIndexCreate
	}

	dur := time.Since(start)
	fmt.Printf("time to index: %s sec\n", dur)

	return db, nil
}

// ListIndexes ...
func ListIndexes(bdb *buntdb.DB) {
	indexes, err := bdb.Indexes()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("db indexes: %v\n", indexes)
}

// Update ...
func Update(bdb *buntdb.DB, name, id, bounds string) error {
	// fmt.Println("updating db with %s:%s\n", id, bounds)

	return bdb.Update(func(tx *buntdb.Tx) error {
		k := fmt.Sprintf("%s:%s", name, id)
		v := bounds
		tx.Set(k, v, nil)
		return nil
	})
}

// Get ...
func Get(bdb *buntdb.DB, index, bounds string) ([]string, error) {
	// fmt.Println("looking for %s\n", bounds)

	results := []string{}

	err := bdb.View(func(tx *buntdb.Tx) error {
		tx.Intersects(index, bounds, func(k, v string) bool {
			results = append(results, fmt.Sprintf("%s:%s", k, v))
			return true
		})
		return nil
	})

	return results, err
}

// ParseResult ...
func ParseResult(result string) (string, string, string) {
	r := strings.Split(result, ":")

	if len(r) != 3 {
		return "", "", ""
	}

	index := r[0]
	id := r[1]
	bounds := r[2]

	return index, id, bounds
}
