package query

import (
	"fmt"
	"strconv"

	"github.com/engelsjk/rtyq"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
)

var (
	ErrInvalidTile           error = fmt.Errorf("invalid z/x/y tile")
	ErrTileZoomLimitExceeded error = fmt.Errorf("tile zoom limit exceeded")
)

// GetFeaturesFromTile parses a tile string 'z/x/y',
// queries the database for results and returns
// the results as a slice of *geojson.Feature
func GetFeaturesFromTile(z, x, y string, zoomLimit int, db rtyq.DB, data rtyq.Data) ([]*geojson.Feature, error) {

	tile, err := ParseTile(z, x, y)
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
func ParseTile(zs, xs, ys string) (maptile.Tile, error) {

	z, err := strconv.ParseInt(zs, 10, 32)
	if err != nil {
		return maptile.Tile{}, ErrInvalidTile
	}
	x, err := strconv.ParseInt(xs, 10, 32)
	if err != nil {
		return maptile.Tile{}, ErrInvalidTile
	}
	y, err := strconv.ParseInt(ys, 10, 32)
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
func ResolveFeaturesFromTile(tile maptile.Tile, results rtyq.Results, data rtyq.Data) []*geojson.Feature {

	features := []*geojson.Feature{}

	tileset := make(maptile.Set)
	tileset[tile] = true

	tilesetN := uptile(tileset, 3)

	// iterate  over results
	// and check if tiles overlap feature geometry

	for k := range results {

		defer delete(results, k)

		_, id := rtyq.ParseKey(k)

		fp := rtyq.FilePath(data.DirPath, id, data.FileExtension)

		f, err := rtyq.LoadFeature(fp)
		if err != nil {
			continue
		}

		if doTilesOverlapGeometry(f.Geometry, tilesetN) {
			features = appendFeature(features, f)
		}
	}

	return features
}
