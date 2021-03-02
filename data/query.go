package data

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/paulmach/orb/maptile"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
)

var (
	ErrQueryMissingLayer error = fmt.Errorf("missing layer")
	ErrQueryInvalidLayer error = fmt.Errorf("invalid layer")
	ErrQueryMissingQuery error = fmt.Errorf("missing query")
	ErrQueryInvalidQuery error = fmt.Errorf("invalid query")
	ErrQueryMissingID    error = fmt.Errorf("missing id")
	ErrQueryInvalidID    error = fmt.Errorf("invalid id")
	ErrQueryMissingPoint error = fmt.Errorf("missing point")
	ErrQueryInvalidPoint error = fmt.Errorf("invalid point")
	ErrQueryMissingTile  error = fmt.Errorf("missing tile")
	ErrQueryInvalidTile  error = fmt.Errorf("invalid tile")
	ErrQueryMissingBBox  error = fmt.Errorf("missing bbox")
	ErrQueryInvalidBBox  error = fmt.Errorf("invalid bbox")
	// ErrQueryExceededTileZoomLimit error = fmt.Errorf("exceeded tile zoom limit")
	ErrQueryRequest error = fmt.Errorf("unable to make request")
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

func (q Query) Layers() []string {
	layers := []string{}
	for l := range q.layers {
		layers = append(layers, l)
	}
	return layers
}

func (q Query) Point(layer, pt string) (*[]geojson.Feature, error) {

	if layer == "" {
		return &[]geojson.Feature{}, ErrQueryMissingLayer
	}

	if !q.HasLayer(layer) {
		return &[]geojson.Feature{}, ErrQueryInvalidLayer
	}

	if pt == "" {
		return &[]geojson.Feature{}, ErrQueryMissingPoint
	}

	point := parsePoint(pt)
	if point == nil {
		return &[]geojson.Feature{}, ErrQueryInvalidPoint
	}

	features, err := q.layers[layer].intersects(*point)
	if err != nil {
		return &[]geojson.Feature{}, ErrQueryRequest
	}

	if len(features) == 0 {
		return &[]geojson.Feature{}, nil
	}

	return &features, nil
}

func (q Query) BBox(layer, bb string) (*[]geojson.Feature, error) {

	if layer == "" {
		return &[]geojson.Feature{}, ErrQueryMissingLayer
	}

	if !q.HasLayer(layer) {
		return &[]geojson.Feature{}, ErrQueryInvalidLayer
	}

	if bb == "" {
		return &[]geojson.Feature{}, ErrQueryMissingBBox
	}

	bbox := parseBBox(bb)
	if bbox == "" {
		return &[]geojson.Feature{}, ErrQueryInvalidBBox
	}

	// if len(features) == 0 {
	// 	return &[]geojson.Feature{}, nil
	// }

	return &[]geojson.Feature{}, nil
}

func (q Query) Tile(layer, x, y, z string) (*[]geojson.Feature, error) {

	if layer == "" {
		return &[]geojson.Feature{}, ErrQueryMissingLayer
	}

	if !q.HasLayer(layer) {
		return &[]geojson.Feature{}, ErrQueryInvalidLayer
	}

	if x == "" || y == "" || z == "" {
		return &[]geojson.Feature{}, ErrQueryMissingTile
	}

	tile := parseTile(x, y, z)
	if tile == nil {
		return &[]geojson.Feature{}, ErrQueryInvalidTile
	}

	// if int(tile.Z) < q.layers[layer].ZoomLimit {
	// 	return &[]geojson.Feature{}, ErrQueryExceededTileZoomLimit
	// }

	features, err := q.layers[layer].intersects(*tile)
	if err != nil {
		return &[]geojson.Feature{}, ErrQueryRequest
	}

	if len(features) == 0 {
		return &[]geojson.Feature{}, nil
	}

	return &features, nil
}

func (q Query) ID(layer, id string) (*[]geojson.Feature, error) {

	if layer == "" {
		return &[]geojson.Feature{}, ErrQueryMissingLayer
	}

	if !q.HasLayer(layer) {
		return &[]geojson.Feature{}, ErrQueryInvalidLayer
	}

	if id == "" {
		return &[]geojson.Feature{}, ErrQueryMissingID
	}

	// add id validation if needed

	fp := filePath(q.layers[layer].DataDir, id, q.layers[layer].DataExt)

	f, _, err := feature(fp)
	if err != nil {
		return &[]geojson.Feature{}, nil
	}

	return &[]geojson.Feature{*f}, nil
}

///////////////////////////////////////////////////////////////////////////////////////

func parsePoint(pt string) *geom.Point {

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

	point := geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{lon, lat})

	return point
}

func parseBBox(bb string) string {
	return ""
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
