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

// LoadFeature opens, reads and unmarshals a GeoJSON Feature from the input filepath.
func loadFeature(path string) (*geojson.Feature, error) {

	file, err := os.Open(path)

	if err != nil {
		if file != nil {
			file.Close()
		}
		return nil, err
	}

	buf := bytes.Buffer{}

	io.Copy(&buf, file)
	file.Close()

	b := buf.Bytes()

	f, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func bounds(o interface{}) string {

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

// FilePath returns the full filepath from a directory, file ID and file extension.
func filePath(dir, id, ext string) string {
	var sb strings.Builder
	sb.WriteString(id)
	sb.WriteString(ext)
	fp := filepath.Join(dir, sb.String())
	return fp
}

// FileExists ...
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
