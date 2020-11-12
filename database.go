package rtyq

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

type DB struct {
	FilePath string
	FileName string
	Index    string
	db       *buntdb.DB
}

// NewDB ...
func NewDB(path string) (*DB, error) {

	fn := filepath.Base(path)

	db := &DB{
		FilePath: path,
		FileName: fn,
	}

	bdb, err := db.Create()
	if err != nil {
		return nil, err
	}

	db.db = bdb

	return db, nil
}

// LoadDB ...
func LoadDB(path string) (*DB, error) {

	fn := filepath.Base(path)

	db := &DB{
		FilePath: path,
		FileName: fn,
	}

	bdb, err := db.Load()
	if err != nil {
		return nil, err
	}

	db.db = bdb

	return db, nil
}

// InitDB ...
func InitDB(path string) (*DB, error) {
	db, err := NewDB(path)

	if err == nil {
		return db, nil
	}

	if err == ErrDatabaseFileAlreadyExists {
		return LoadDB(path)
	}

	return nil, err
}

// Create ...
func (db *DB) Create() (*buntdb.DB, error) {

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

// Load ...
func (db *DB) Load() (*buntdb.DB, error) {

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

// CreateSpatialIndex ...
func (db *DB) CreateSpatialIndex(index string) error {

	db.Index = ""

	name := index
	pattern := fmt.Sprintf("%s:*", index)

	fmt.Printf("running spatial index:%s (%s)...\n", index, db.FileName)

	start := time.Now()

	err := db.db.CreateSpatialIndex(name, pattern, buntdb.IndexRect)
	if err != nil {
		return ErrSpatialIndexCreate
	}

	db.Index = index

	dur := time.Since(start)
	fmt.Printf("time to index: %s\n", dur)

	return nil
}

// ListIndexes ...
func (db *DB) ListIndexes() {
	indexes, err := db.db.Indexes()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("db indexes: %v\n", indexes)
}

// Update ...
func (db *DB) Update(id, bounds string) error {
	// fmt.Println("updating db with %s:%s\n", id, bounds)

	return db.db.Update(func(tx *buntdb.Tx) error {
		k := fmt.Sprintf("%s:%s", db.Index, id)
		v := bounds
		tx.Set(k, v, nil)
		return nil
	})
}

// GetResults ...
func (db *DB) GetResults(bounds string) ([]Result, error) {
	// fmt.Println("looking for %s\n", bounds)

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

// Result ...
type Result struct {
	Index  string
	ID     string
	Bounds string
}

// ParseResult ...
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

// Bounds ...
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
