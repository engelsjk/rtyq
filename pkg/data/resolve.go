package data

import (
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/maptile/tilecover"
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

		if isPtInFeature {
			features = append(features, f)
		}
	}

	return features, nil
}

// ResolveTile ...
func ResolveTile(path, ext string, results []string, tile maptile.Tile) ([]*geojson.Feature, error) {

	var zoomOffset = 5
	var zoomMax = 22

	var newZoom maptile.Zoom
	features := []*geojson.Feature{}

	// uptile and get min/max tiles

	zoomInt := int(tile.Z)
	if zoomInt+zoomOffset > zoomMax {
		newZoom = maptile.Zoom(zoomMax)
	} else {
		newZoom = maptile.Zoom(zoomInt + zoomOffset)
	}
	minTile, maxTile := tile.Range(newZoom)

	// iterate  over results...

	for _, r := range results {

		_, id, _ := db.ParseResult(r)

		fp := FilePath(path, id, ext)

		f, err := LoadFeature(fp)
		if err != nil {
			continue
		}

		// check if tile center is in feature (only for polygons and multipolygons)...

		tileCenter := tile.Bound().Center()
		isTileCenterInFeature := false
		geom := f.Geometry
		switch g := geom.(type) {
		case orb.Polygon:
			isTileCenterInFeature = planar.PolygonContains(g, tileCenter)
		case orb.MultiPolygon:
			isTileCenterInFeature = planar.MultiPolygonContains(g, tileCenter)
		}
		if isTileCenterInFeature {
			features = append(features, f)
			continue
		}

		// todo: does the above work for non-polygon or non-multipolygon cases?
		// todo: are points and linestrings/multilinestrings handled by the below?

		// iterate over feature tileset at uptile zoom level, and check for any tile matchs

		tileSet := tilecover.Geometry(geom, newZoom)
		for tile := range tileSet {
			if (tile.X >= minTile.X && tile.Y >= minTile.Y) &&
				(tile.X <= maxTile.X && tile.Y <= maxTile.Y) {
				features = append(features, f)
				break
			}
		}
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

	features = append(features, f)
	return features, nil
}
