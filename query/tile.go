package query

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/engelsjk/rtyq"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/maptile/tilecover"
	"github.com/paulmach/orb/planar"
)

var (
	ErrInvalidTile           error = fmt.Errorf("invalid z/x/y tile")
	ErrTileZoomLimitExceeded error = fmt.Errorf("tile zoom limit exceeded")
)

// GetFeaturesFromTile parses a tile string 'z/x/y',
// queries the database for results and returns
// the results as a slice of *geojson.Feature
func GetFeaturesFromTile(t string, zoomLimit int, db *rtyq.DB, data *rtyq.Data) ([][]byte, error) {

	tile, err := ParseTile(t)
	if err != nil {
		return nil, err
	}

	if zoomLimit != 0 {
		if int(tile.Z) < zoomLimit {
			return nil, ErrTileZoomLimitExceeded
		}
	}

	results, err := db.GetResults(rtyq.Bounds(tile))
	if err != nil {
		return nil, err
	}

	features := ResolveFeaturesFromTile(tile, results, data)

	return features, nil
}

// ParseTile converts a tile string 'z/x/y' to a maptile.Tile object
func ParseTile(t string) (maptile.Tile, error) {

	spl := strings.Split(t, "/")

	if len(spl) != 3 {
		return maptile.Tile{}, ErrInvalidTile
	}

	z, err := strconv.ParseInt(spl[0], 10, 32)
	if err != nil {
		return maptile.Tile{}, ErrInvalidTile
	}
	x, err := strconv.ParseInt(spl[1], 10, 32)
	if err != nil {
		return maptile.Tile{}, ErrInvalidTile
	}
	y, err := strconv.ParseInt(spl[2], 10, 32)
	if err != nil {
		return maptile.Tile{}, ErrInvalidTile
	}

	tile := maptile.Tile{
		X: uint32(x),
		Y: uint32(y),
		Z: maptile.Zoom(uint32(z)),
	}

	return tile, nil
}

// ResolveFeaturesFromTile converts the results from a database query,
// loads GeoJSON data from the data directory and returns a slice of *geojson.Feature
func ResolveFeaturesFromTile(tile maptile.Tile, results []rtyq.Result, data *rtyq.Data) [][]byte {

	var zoomOffset = 5
	var zoomMax = 22

	var newZoom maptile.Zoom
	features := [][]byte{}

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

		fp := rtyq.FilePath(data.DirPath, r.ID, data.FileExtension)

		f, err := rtyq.LoadFeature(fp)
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
			features, err = appendFeature(features, f)
			continue
		}

		// todo: does the above work for non-polygon or non-multipolygon cases?
		// todo: are points and linestrings/multilinestrings handled by the below?

		// iterate over feature tileset at uptile zoom level, and check for any tile matchs

		tileSet := tilecover.Geometry(geom, newZoom)
		for tile := range tileSet {
			if (tile.X >= minTile.X && tile.Y >= minTile.Y) &&
				(tile.X <= maxTile.X && tile.Y <= maxTile.Y) {
				features, _ = appendFeature(features, f)
				break
			}
		}
	}

	return features
}
