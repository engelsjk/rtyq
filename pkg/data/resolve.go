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

	zoomInt := int(tile.Z)

	if zoomInt+zoomOffset > zoomMax {
		newZoom = maptile.Zoom(zoomMax)
	} else {
		newZoom = maptile.Zoom(zoomInt + zoomOffset)
	}

	minTile, maxTile := tile.Range(newZoom)

	for _, r := range results {

		_, id, _ := db.ParseResult(r)

		fp := FilePath(path, id, ext)

		f, err := LoadFeature(fp)
		if err != nil {
			continue
		}

		// todo: is this too slow if original tile is entirely contained within the polygon?
		// todo: how to avoid slowdown caused by looping through large geometry at low zoom levels?

		tileSet := tilecover.Geometry(f.Geometry, newZoom)
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
