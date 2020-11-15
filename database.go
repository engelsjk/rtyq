package rtyq

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
	"github.com/tidwall/buntdb"
)

var (
	ErrDatabaseFileAlreadyExists  error = fmt.Errorf("database file already exists")
	ErrDatabaseFileCreate         error = fmt.Errorf("unable to create database file")
	ErrDatabaseFileDoesNotExist   error = fmt.Errorf("database file does not exist")
	ErrDatabaseFileIsADir         error = fmt.Errorf("database file is a dir, needs to be a file")
	ErrDatabaseFileOpen           error = fmt.Errorf("unable to open database file")
	ErrSpatialIndexCreate         error = fmt.Errorf("unable to spatially index database file")
	ErrDatabaseFailedToGetResults error = fmt.Errorf("failed to get results")
)

// DB is a structure that contains information
// about the database file for a single layer
// (FilePath, FileName and Index)
type DB struct {
	FilePath string
	FileName string
	Index    string
	db       *buntdb.DB
}

// NewDB .,,
func NewDB(path string) DB {
	fn := filepath.Base(path)

	return DB{
		FilePath: path,
		FileName: fn,
	}
}

// InitDB initializes an DB object, then loads a database
// if a file already exists or creates a new database if it does not
func InitDB(path string) (DB, error) {

	db := NewDB(path)

	bdb, err := db.Create()
	if err == nil {
		db.db = bdb
		return db, nil
	}

	bdb, err = db.Load()
	if err != nil {
		return DB{}, err

	}

	db.db = bdb

	return db, nil
}

// Create initializes a new buntdb.DB object
// and creates a new database file
func (db DB) Create() (*buntdb.DB, error) {

	_, err := os.Stat(db.FilePath)

	if !os.IsNotExist(err) {
		return nil, ErrDatabaseFileAlreadyExists
	}

	bdb, err := buntdb.Open(db.FilePath)
	if err != nil {
		return nil, ErrDatabaseFileCreate
	}

	return bdb, nil
}

// Load initializes a new buntdb.DB object
// from an existing database file
func (db DB) Load() (*buntdb.DB, error) {

	info, err := os.Stat(db.FilePath)
	if os.IsNotExist(err) {
		return nil, ErrDatabaseFileDoesNotExist
	}

	if info.IsDir() {
		return nil, ErrDatabaseFileIsADir
	}

	bdb, err := buntdb.Open(db.FilePath)
	if err != nil {
		return nil, ErrDatabaseFileOpen
	}

	return bdb, nil
}

// CreateSpatialIndex runs an r-tree spatial index on the database object
// using the specified index name
func (db DB) CreateSpatialIndex(index string) error {

	db.Index = ""

	name := index
	pattern := fmt.Sprintf("%s:*", index)

	err := db.db.CreateSpatialIndex(name, pattern, buntdb.IndexRect)
	if err != nil {
		return ErrSpatialIndexCreate
	}

	db.Index = index

	return nil
}

// ListIndexes prints out the existing indexes in a buntdb.DB object
func (db DB) ListIndexes() {

	indexes, err := db.db.Indexes()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("db indexes: %v\n", indexes)
}

// Update adds an object to the database
// with a key of index:id and a value of a geometry bounds
func (db DB) Update(id, bounds string) error {

	return db.db.Update(func(tx *buntdb.Tx) error {
		k := fmt.Sprintf("%s:%s", db.Index, id)
		v := bounds
		tx.Set(k, v, nil)
		return nil
	})
}

// GetResults returns a slice of Result objects from the database
// that intersect the given bounds
func (db DB) GetResults(bounds string) ([]Result, error) {

	results := []Result{}

	err := db.db.View(func(tx *buntdb.Tx) error {
		tx.Intersects(db.Index, bounds, func(k, v string) bool {
			kv := fmt.Sprintf("%s:%s", k, v)
			r := ParseResult(kv)
			results = append(results, r)

			return true
		})
		return nil
	})

	if err != nil {
		return nil, ErrDatabaseFailedToGetResults
	}

	return results, nil
}

// Result is an object that contains Index, ID and Bounds
// of an object returned from a database
type Result struct {
	Index  string
	ID     string
	Bounds string
}

// ParseResult converts a string output from a database query
// into a Result object
func ParseResult(result string) Result {

	r := strings.Split(result, ":")

	if len(r) != 3 {
		return Result{}
	}

	index := r[0]
	id := r[1]
	bounds := r[2]

	return Result{Index: index, ID: id, Bounds: bounds}
}

// Bounds converts a geometry object (eg orb.Point or maptile.Tile)
// to the string format required to query the database
func Bounds(o interface{}) string {

	var bounds string

	switch v := o.(type) {
	case orb.Point:
		bounds = fmt.Sprintf("[%f %f]", v.Lon(), v.Lat())
	case maptile.Tile:
		bounds = fmt.Sprintf("[%f %f], [%f %f]",
			v.Bound().Min.Lon(),
			v.Bound().Min.Lat(),
			v.Bound().Max.Lon(),
			v.Bound().Max.Lat(),
		)
	default:
		//
	}

	return bounds
}
