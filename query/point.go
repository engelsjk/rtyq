package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/engelsjk/rtyq"
	"github.com/paulmach/orb"
)

var (
	ErrInvalidPoint error = fmt.Errorf("invalid lon,lat point")
)

// GetFeaturesFromPoint parses a point string 'lon,lat',
// queries the database for results and returns
// the results as a slice of *geojson.Feature
func GetFeaturesFromPoint(pt string, db rtyq.DB, data rtyq.Data) ([][]byte, error) {

	point, err := ParsePoint(pt)
	if err != nil {
		return nil, err
	}

	results, err := db.GetResults(rtyq.Bounds(point))
	if err != nil {
		return nil, err
	}

	features := ResolveFeaturesFromPoint(point, results, data)

	return features, nil
}

// ParsePoint converts a point string 'lon, lat' to an orb.Point object
func ParsePoint(pt string) (orb.Point, error) {

	// todo: better latlon string validation

	cleanLatLon := strings.ReplaceAll(pt, " ", "")
	splitLatLon := strings.Split(cleanLatLon, ",")

	lon, err := strconv.ParseFloat(splitLatLon[0], 64)
	if err != nil {
		return orb.Point{}, ErrInvalidPoint
	}

	lat, err := strconv.ParseFloat(splitLatLon[1], 64)
	if err != nil {
		return orb.Point{}, ErrInvalidPoint
	}

	point := orb.Point{lon, lat}

	return point, nil
}

// ResolveFeaturesFromPoint converts the results from a database query,
// loads GeoJSON data from the data directory and returns a slice of *geojson.Feature
func ResolveFeaturesFromPoint(pt orb.Point, results rtyq.Results, data rtyq.Data) [][]byte {

	features := [][]byte{}

	// iterate  over results
	// and check if point is in feature geometry

	for k := range results {

		defer delete(results, k)

		_, id := rtyq.ParseKey(k)

		fp := rtyq.FilePath(data.DirPath, id, data.FileExtension)

		f, err := rtyq.LoadFeature(fp)
		if err != nil {
			continue
		}

		if isPointInFeature(f.Geometry, pt) {
			features, _ = appendFeature(features, f)
		}
	}

	return features
}
