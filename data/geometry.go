package data

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

func pointInGeometry(geom orb.Geometry, pt orb.Point) bool {
	switch g := geom.(type) {
	case orb.Polygon:
		return planar.PolygonContains(g, pt)
	case orb.MultiPolygon:
		return planar.MultiPolygonContains(g, pt)
	default:
		return false
	}
}

// boundOverlapsGeometry is a rough approximation
// it ignores intermediate overlaps not at corners or center
func boundOverlapsGeometry(geom orb.Geometry, bound orb.Bound) bool {
	return boundCenterInGeometry(geom, bound) || boundCornersInGeometry(geom, bound)
}

// tileOverlapsGeometry is a partial approximation
// it ignores intermediate overlaps not at corners or center
// however, it will uptile to a higher zoom and recheck overlaps
func tileOverlapsGeometry(geom orb.Geometry, tile maptile.Tile) bool {

	if boundCenterInGeometry(geom, tile.Bound()) && boundCornersInGeometry(geom, tile.Bound()) {
		return true
	}

	set := make(maptile.Set)
	set[tile] = true
	upset := uptile(set, 3)

	if tilesetOverlapsGeometry(geom, upset) {
		return true
	}

	return false
}

func tilesetOverlapsGeometry(geom orb.Geometry, tileset maptile.Set) bool {
	for t := range tileset {
		if boundCenterInGeometry(geom, t.Bound()) {
			return true
		}
	}
	return false
}

func boundCornersInGeometry(geom orb.Geometry, bound orb.Bound) bool {
	p := bound.ToPolygon()
	for _, pt := range p[0] {
		if pointInGeometry(geom, pt) {
			return true
		}
	}
	return false
}

func boundCenterInGeometry(geom orb.Geometry, bound orb.Bound) bool {
	return pointInGeometry(geom, bound.Center())
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
