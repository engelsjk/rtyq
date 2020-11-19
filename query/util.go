package query

import (
	"encoding/json"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

// EmptyResponse ...
func EmptyResponse() []byte {
	return []byte("[]")
}

// FeaturesToResponse converts a slice of geojson.Features to a byte slice
func FeaturesToResponse(fs []*geojson.Feature) []byte {

	if fs == nil {
		return EmptyResponse()
	}

	resp, err := json.Marshal(&fs)
	if err != nil {
		return EmptyResponse()
	}

	return resp
}

func appendFeature(fs []*geojson.Feature, f *geojson.Feature) {
	fs = append(fs, f)
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
