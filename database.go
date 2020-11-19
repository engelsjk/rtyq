package rtyq

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

// Load initializes a new buntdb.DB object
// from an existing database file
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

// CreateSpatialIndex runs an r-tree spatial index on the database object
// using the specified index name
func (db *DB) CreateSpatialIndex(index string) error {

	db.Index = index

	pattern := db.pattern()

	err := db.db.CreateSpatialIndex(index, pattern, buntdb.IndexRect)
	if err != nil {
		db.Index = ""
		return ErrSpatialIndexCreate
	}

	return nil
}

func (db *DB) pattern() string {
	var sb strings.Builder
	sb.WriteString(db.Index)
	sb.WriteString(":*")
	pattern := sb.String()
	return pattern
}

func (db *DB) key(id string) string {
	var sb strings.Builder
	sb.WriteString(db.Index)
	sb.WriteString(":")
	sb.WriteString(id)
	key := sb.String()
	return key
}

// ListIndexes prints out the existing indexes in a buntdb.DB object
func (db *DB) ListIndexes() {

	indexes, err := db.db.Indexes()
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("db indexes: %v\n", indexes)
}

// Update adds an object to the database
// with a key of index:id and a value of a geometry bounds
func (db *DB) Update(id, bounds string) error {

	return db.db.Update(func(tx *buntdb.Tx) error {
		k := db.key(id)
		v := bounds
		tx.Set(k, v, nil)
		return nil
	})
}

// GetResults returns a slice of Result objects from the database
// that intersect the given bounds
func (db *DB) GetResults(bounds string) (Results, error) {

	results := make(Results)

	err := db.db.View(func(tx *buntdb.Tx) error {
		tx.Intersects(db.Index, bounds, func(k, v string) bool {
			results[k] = v
			return true
		})
		return nil
	})

	if err != nil {
		return nil, ErrDatabaseFailedToGetResults
	}

	return results, nil
}

// Results is a map that contains Index:ID (key) and Bounds (value)
// of all objects returned from a database
type Results map[string]string

// ParseKey ...
func ParseKey(k string) (string, string) {
	r := strings.Split(k, ":")
	index := r[0]
	id := r[1]
	return index, id
}

// Bounds converts a geometry object (eg orb.Point or maptile.Tile)
// to the string format required to query the database
func Bounds(o interface{}) string {

	var sb strings.Builder

	switch v := o.(type) {
	case orb.Point:

		// bounds = [lon lat]

		lon := strconv.FormatFloat(v.Lon(), 'f', -1, 64)
		lat := strconv.FormatFloat(v.Lat(), 'f', -1, 64)

		sb.WriteString("[")
		sb.WriteString(lon)
		sb.WriteString(" ")
		sb.WriteString(lat)
		sb.WriteString("]")

	case orb.Polygon:

		// bounds = [minLon minLat], [maxLon maxLat]

		minLon := strconv.FormatFloat(v.Bound().Min.Lon(), 'f', -1, 64)
		minLat := strconv.FormatFloat(v.Bound().Min.Lat(), 'f', -1, 64)
		maxLon := strconv.FormatFloat(v.Bound().Max.Lon(), 'f', -1, 64)
		maxLat := strconv.FormatFloat(v.Bound().Max.Lat(), 'f', -1, 64)

		sb.WriteString("[")
		sb.WriteString(minLon)
		sb.WriteString(" ")
		sb.WriteString(minLat)
		sb.WriteString("], [")
		sb.WriteString(maxLon)
		sb.WriteString(" ")
		sb.WriteString(maxLat)
		sb.WriteString("]")

	case orb.MultiPolygon:

		// bounds = [minLon minLat], [maxLon maxLat]

		minLon := strconv.FormatFloat(v.Bound().Min.Lon(), 'f', -1, 64)
		minLat := strconv.FormatFloat(v.Bound().Min.Lat(), 'f', -1, 64)
		maxLon := strconv.FormatFloat(v.Bound().Max.Lon(), 'f', -1, 64)
		maxLat := strconv.FormatFloat(v.Bound().Max.Lat(), 'f', -1, 64)

		sb.WriteString("[")
		sb.WriteString(minLon)
		sb.WriteString(" ")
		sb.WriteString(minLat)
		sb.WriteString("], [")
		sb.WriteString(maxLon)
		sb.WriteString(" ")
		sb.WriteString(maxLat)
		sb.WriteString("]")

	case maptile.Tile:

		// bounds = [minLon minLat], [maxLon maxLat]

		minLon := strconv.FormatFloat(v.Bound().Min.Lon(), 'f', -1, 64)
		minLat := strconv.FormatFloat(v.Bound().Min.Lat(), 'f', -1, 64)
		maxLon := strconv.FormatFloat(v.Bound().Max.Lon(), 'f', -1, 64)
		maxLat := strconv.FormatFloat(v.Bound().Max.Lat(), 'f', -1, 64)

		sb.WriteString("[")
		sb.WriteString(minLon)
		sb.WriteString(" ")
		sb.WriteString(minLat)
		sb.WriteString("], [")
		sb.WriteString(maxLon)
		sb.WriteString(" ")
		sb.WriteString(maxLat)
		sb.WriteString("]")

	default:
		//
	}

	bounds := sb.String()

	return bounds
}
