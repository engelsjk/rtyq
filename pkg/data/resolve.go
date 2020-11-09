package data

import (
	"fmt"

	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

// ResolvePoint ...
func ResolvePoint(path, ext string, results []string, pt orb.Point) ([]*geojson.Feature, error) {

	features := []*geojson.Feature{}

	for _, r := range results {

		_, id, _ := db.ParseResult(r)

		fp := FilePath(path, id, ext)

		f, err := LoadFeature(fp)
		if err != nil {
			continue
		}

		ptInFeature := false

		geom := f.Geometry
		switch g := geom.(type) {
		case orb.Polygon:
			ptInFeature = planar.PolygonContains(g, pt)
		case orb.MultiPolygon:
			ptInFeature = planar.MultiPolygonContains(g, pt)
		default:
			continue
		}

		if ptInFeature {
			features = append(features, f)
		}
	}

	return features, nil
}

// ResolveTile ...
func ResolveTile(path, ext string, results []string, tile maptile.Tile) ([]*geojson.Feature, error) {

	features := []*geojson.Feature{}

	for _, r := range results {

		_, id, _ := db.ParseResult(r)

		fp := FilePath(path, id, ext)

		f, err := LoadFeature(fp)
		if err != nil {
			continue
		}

		features = append(features, f)
	}

	return features, nil
}

// ResolveID ...
func ResolveID(path, ext string, id string) ([]*geojson.Feature, error) {
	features := []*geojson.Feature{}

	fp := FilePath(path, id, ext)

	f, err := LoadFeature(fp)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", f)

	features = append(features, f)
	return features, nil
}
