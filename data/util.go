package data

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/tidwall/buntdb"
)

// paths

func filePath(dir, id, ext string) string {
	var sb strings.Builder
	sb.WriteString(id)
	sb.WriteString(ext)
	fp := filepath.Join(dir, sb.String())
	return fp
}

func filename(path string) string {
	return filepath.Base(path)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.Mode().IsRegular()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// files

func read(path, id string) (int64, string, string, error) {

	f, nbytes, err := feature(path)
	if err != nil {
		return 0, "", "", err
	}

	return nbytes, fid(f, id), bounds(f.Geometry), nil
}

func feature(path string) (*geojson.Feature, int64, error) {

	file, err := os.Open(path)

	if err != nil {
		if file != nil {
			file.Close()
		}
		return nil, 0, err
	}

	buf := bytes.Buffer{}

	nbytes, _ := io.Copy(&buf, file)
	file.Close()

	b := buf.Bytes()

	f, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil, 0, err
	}

	return f, nbytes, nil
}

func maxBytes(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func minBytes(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}

func validExtension(path, ext string) bool {
	fext := filepath.Ext(path)
	if fext != ext {
		return false
	}
	return true
}

func fid(f *geojson.Feature, key string) string {

	id, ok := f.Properties[key]
	if !ok {
		return ""
	}

	var fid string

	if v, ok := id.(string); ok {
		fid = v
	}
	if v, ok := id.(int); ok {
		fid = strconv.FormatInt(int64(v), 10)
	}
	if v, ok := id.(float64); ok {
		fid = strconv.FormatFloat(v, 'f', -1, 64)
	}

	return fid
}

func bounds(o interface{}) string {

	var bounds string

	switch v := o.(type) {
	case orb.Point:
		bounds = dbPointBounds(v)
	case orb.Polygon:
		bounds = dbPolyBounds(v.Bound())
	case orb.MultiPolygon:
		bounds = dbPolyBounds(v.Bound())
	case maptile.Tile:
		bounds = dbPolyBounds(v.Bound())
	default:
		// log unknown type
	}

	return bounds
}

// db

func dbPointBounds(p orb.Point) string {
	// bounds = [lon lat]

	var sb strings.Builder

	lon := strconv.FormatFloat(p.Lon(), 'f', -1, 64)
	lat := strconv.FormatFloat(p.Lat(), 'f', -1, 64)

	sb.WriteString("[")
	sb.WriteString(lon)
	sb.WriteString(" ")
	sb.WriteString(lat)
	sb.WriteString("]")

	bounds := sb.String()

	return bounds
}

func dbPolyBounds(b orb.Bound) string {
	// bounds = [minLon minLat], [maxLon maxLat]

	var sb strings.Builder

	minLon := strconv.FormatFloat(b.Min.Lon(), 'f', -1, 64)
	minLat := strconv.FormatFloat(b.Min.Lat(), 'f', -1, 64)
	maxLon := strconv.FormatFloat(b.Max.Lon(), 'f', -1, 64)
	maxLat := strconv.FormatFloat(b.Max.Lat(), 'f', -1, 64)

	sb.WriteString("[")
	sb.WriteString(minLon)
	sb.WriteString(" ")
	sb.WriteString(minLat)
	sb.WriteString("], [")
	sb.WriteString(maxLon)
	sb.WriteString(" ")
	sb.WriteString(maxLat)
	sb.WriteString("]")

	bounds := sb.String()

	return bounds
}

func dbUpdate(db *buntdb.DB, index, id, bounds string) error {
	return db.Update(func(tx *buntdb.Tx) error {
		k := dbKey(index, id)
		v := bounds
		tx.Set(k, v, nil)
		return nil
	})
}

func dbPattern(index string) string {
	var sb strings.Builder
	sb.WriteString(index)
	sb.WriteString(":*")
	pattern := sb.String()
	return pattern
}

func dbKey(index, id string) string {
	var sb strings.Builder
	sb.WriteString(index)
	sb.WriteString(":")
	sb.WriteString(id)
	key := sb.String()
	return key
}
