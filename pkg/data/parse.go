package data

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

// ParseLonLatPoint ...
func ParseLonLatPoint(p string) (orb.Point, error) {

	// todo: better latlon string validation

	cleanLatLon := strings.ReplaceAll(p, " ", "")
	splitLatLon := strings.Split(cleanLatLon, ",")

	lon, err := strconv.ParseFloat(splitLatLon[0], 64)
	if err != nil {
		return orb.Point{}, err
	}

	lat, err := strconv.ParseFloat(splitLatLon[1], 64)
	if err != nil {
		return orb.Point{}, err
	}

	pt := orb.Point{lon, lat}

	return pt, nil
}

// ParseTile ...
func ParseTile(t string) (maptile.Tile, error) {
	spl := strings.Split(t, "/")

	if len(spl) != 3 {
		return maptile.Tile{}, fmt.Errorf("invalid tile")
	}

	z, err := strconv.ParseInt(spl[0], 10, 32)
	if err != nil {
		return maptile.Tile{}, fmt.Errorf("invalid tile")
	}
	x, err := strconv.ParseInt(spl[1], 10, 32)
	if err != nil {
		return maptile.Tile{}, fmt.Errorf("invalid tile")
	}
	y, err := strconv.ParseInt(spl[2], 10, 32)
	if err != nil {
		return maptile.Tile{}, fmt.Errorf("invalid tile")
	}

	tile := maptile.Tile{
		X: uint32(x),
		Y: uint32(y),
		Z: maptile.Zoom(uint32(z)),
	}

	return tile, nil
}

// ParseID ...
func ParseID(i string) string {
	// placeholder func for future validation
	return i
}
