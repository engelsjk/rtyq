package data

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
)

var (
	ErrQueryNoLayer               error = fmt.Errorf("layer not found")
	ErrQueryInvalidPoint          error = fmt.Errorf("invalid point")
	ErrQueryInvalidTile           error = fmt.Errorf("invalid tile")
	ErrQueryExceededTileZoomLimit error = fmt.Errorf("exceeded tile zoom limit")
	ErrQueryRequest               error = fmt.Errorf("unable to make request")
)

var QueryHandler Query

type Query struct {
	layers map[string]*Layer
}

func init() {
	QueryHandler = Query{
		layers: make(map[string]*Layer),
	}
}

func (q Query) HasLayer(layer string) bool {
	_, ok := q.layers[layer]
	return ok
}

func (q Query) Point(layer, pt string) (*[]geojson.Feature, error) {

	if !q.HasLayer(layer) {
		return nil, ErrQueryNoLayer
	}

	point := parsePoint(pt)
	if point == nil {
		return nil, ErrQueryInvalidPoint
	}

	features, err := q.layers[layer].intersects(*point)
	if err != nil {
		return nil, ErrQueryRequest
	}

	return &features, nil
}

func (q Query) BBox(layer, bbox string) (*[]geojson.Feature, error) {

	return nil, nil
}

func (q Query) Tile(layer, x, y, z string) (*[]geojson.Feature, error) {

	if !q.HasLayer(layer) {
		return nil, ErrQueryNoLayer
	}

	tile := parseTile(x, y, z)
	if tile == nil {
		return nil, ErrQueryInvalidTile
	}

	if int(tile.Z) < q.layers[layer].ZoomLimit {
		return nil, ErrQueryExceededTileZoomLimit
	}

	features, err := q.layers[layer].intersects(*tile)
	if err != nil {
		return nil, ErrQueryRequest
	}

	return &features, nil
}

func (q Query) ID(layer, id string) (*[]geojson.Feature, error) {

	if !q.HasLayer(layer) {
		return nil, ErrQueryNoLayer
	}

	// add id validation if needed

	fp := filePath(q.layers[layer].DataDir, id, q.layers[layer].DataExt)

	f, _, err := feature(fp)
	if err != nil {
		return nil, nil
	}

	return &[]geojson.Feature{*f}, nil
}

///////////////////////////////////////////////////////////////////////////////////////

func parsePoint(pt string) *orb.Point {

	// todo: better latlon string validation?

	cleanLatLon := strings.ReplaceAll(pt, " ", "")
	splitLatLon := strings.Split(cleanLatLon, ",")

	if len(splitLatLon) != 2 {
		return nil
	}

	lon, err := strconv.ParseFloat(splitLatLon[0], 64)
	if err != nil {
		return nil
	}

	lat, err := strconv.ParseFloat(splitLatLon[1], 64)
	if err != nil {
		return nil
	}

	return &orb.Point{lon, lat}
}

func parseTile(xs, ys, zs string) *maptile.Tile {

	x, err := strconv.ParseInt(xs, 10, 32)
	if err != nil {
		return nil
	}
	y, err := strconv.ParseInt(ys, 10, 32)
	if err != nil {
		return nil
	}
	z, err := strconv.ParseInt(zs, 10, 32)
	if err != nil {
		return nil
	}

	return &maptile.Tile{
		X: uint32(x),
		Y: uint32(y),
		Z: maptile.Zoom(uint32(z)),
	}
}
