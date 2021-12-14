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
	case orb.Bound:
		bounds = dbPolyBounds(v)
	case orb.Point:
		bounds = dbPointBounds(v)
	case orb.Polygon:
		bounds = dbPolyBounds(v.Bound())
	case orb.MultiPolygon:
		bounds = dbPolyBounds(v.Bound())
	case maptile.Tile:
		bounds = dbPolyBounds(v.Bound())
	default:
		// log unknown type?
	}
	return bounds
}
