package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/buntdb"
)

// Initialize ...
func Initialize(path, index string, skipIndex bool) (*buntdb.DB, error) {
	// fmt.Println("initializing db at %s\n", path)

	// todo: check if db file exists at path

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	if skipIndex {
		return db, nil
	}

	name := index
	pattern := fmt.Sprintf("%s:*", index)

	start := time.Now()

	err = db.CreateSpatialIndex(name, pattern, buntdb.IndexRect)
	if err != nil {
		return nil, err
	}

	dur := time.Since(start)
	fmt.Printf("time to index db: %s sec\n", dur)

	ListIndexes(db)

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
