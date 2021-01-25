package data

import (
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
)

func pointInFeature(geom orb.Geometry, pt orb.Point) bool {
	switch g := geom.(type) {
	case orb.Polygon:
		return planar.PolygonContains(g, pt)
	case orb.MultiPolygon:
		return planar.MultiPolygonContains(g, pt)
	default:
		return false
	}
}

func tileOverlapsGeometry(geom orb.Geometry, tile maptile.Tile) bool {

	if tileCenterInFeature(geom, tile) && tileCornersInFeature(geom, tile) {
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
		if tileCenterInFeature(geom, t) {
			return true
		}
	}
	return false
}

func tileCornersInFeature(geom orb.Geometry, tile maptile.Tile) bool {
	p := tile.Bound().ToPolygon()
	for _, pt := range p[0] {
		if !pointInFeature(geom, pt) {
			return false
		}
	}
	return true
}

func tileCenterInFeature(geom orb.Geometry, tile maptile.Tile) bool {
	return pointInFeature(geom, tile.Bound().Center())
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
