package data

import (
	"github.com/engelsjk/planeta/geo"
	"github.com/engelsjk/planeta/geo/geomfn"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

func pointInFeature(pt geom.Point, f *geojson.Feature) bool {
	fg, err := geo.MakeGeometryFromGeomT(f.Geometry)
	if err != nil {
		return false
	}

	ptg, err := geo.MakeGeometryFromPointCoords(pt.X(), pt.Y())
	if err != nil {
		return false
	}

	contains, err := geomfn.Contains(fg, ptg)
	if err != nil {
		return false
	}
	return contains
}

func geometryIntersectsFeature(g geom.T, f *geojson.Feature) bool {
	fg, err := geo.MakeGeometryFromGeomT(f.Geometry)
	if err != nil {
		return false
	}

	gg, err := geo.MakeGeometryFromGeomT(g)
	if err != nil {
		return false
	}

	intersects, err := geomfn.Intersects(fg, gg)
	if err != nil {
		return false
	}
	return intersects
}
