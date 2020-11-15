package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/engelsjk/rtyq"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

var (
	ErrInvalidPoint error = fmt.Errorf("invalid lon,lat point")
)

// GetFeaturesFromPoint parses a point string 'lon,lat',
// queries the database for results and returns
// the results as a slice of *geojson.Feature
func GetFeaturesFromPoint(pt string, db *rtyq.DB, data *rtyq.Data) ([][]byte, error) {

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
func ResolveFeaturesFromPoint(pt orb.Point, results []rtyq.Result, data *rtyq.Data) [][]byte {

	features := make([][]byte, len(results))

	for _, r := range results {

		fp := rtyq.FilePath(data.DirPath, r.ID, data.FileExtension)

		f, err := rtyq.LoadFeature(fp)
		if err != nil {
			continue
		}

		isPtInFeature := false

		geom := f.Geometry
		switch g := geom.(type) {
		case orb.Polygon:
			isPtInFeature = planar.PolygonContains(g, pt)
		case orb.MultiPolygon:
			isPtInFeature = planar.MultiPolygonContains(g, pt)
		default:
			continue
		}

		if !isPtInFeature {
			continue
		}

		features, _ = appendFeature(features, f)
	}

	return features
}
