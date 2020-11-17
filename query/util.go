package query

import (
	"bytes"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

// FeaturesToString converts a slice of geojson.Features to a comma-separated
// array string
func FeaturesToString(features [][]byte) string {

	if features == nil {
		return "[]"
	}

	out := fmt.Sprintf("[%s]", string(bytes.Join(features, []byte(","))))
	return out
}

func appendFeature(bs [][]byte, f *geojson.Feature) ([][]byte, error) {
	b, err := f.MarshalJSON()
	if err != nil {
		return bs, err
	}
	bs = append(bs, b)
	return bs, nil
}

func isPointInFeature(geom orb.Geometry, pt orb.Point) bool {
	switch g := geom.(type) {
	case orb.Polygon:
		return planar.PolygonContains(g, pt)
	case orb.MultiPolygon:
		return planar.MultiPolygonContains(g, pt)
	default:
		return false
	}
}

func doTilesOverlapGeometry(geom orb.Geometry, tileset maptile.Set) bool {

	for t := range tileset {
		if isTileCenterInFeature(geom, t) {
			return true
		}
	}

	return false
}

func isTileCenterInFeature(geom orb.Geometry, tile maptile.Tile) bool {

	tileCenter := tile.Bound().Center()

	switch g := geom.(type) {
	case orb.Polygon:
		return planar.PolygonContains(g, tileCenter)
	case orb.MultiPolygon:
		return planar.MultiPolygonContains(g, tileCenter)
	default:
		return false
	}
}

func uptile(set maptile.Set, n int) maptile.Set {

	if n <= 0 {
		return set
	}

	childSet := maptile.Set{}
	for tile := range set {
		tiles := tile.Children()
		for _, t := range tiles {
			childSet[t] = true
		}
	}
	return uptile(childSet, n-1)
}
