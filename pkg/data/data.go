package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/karrick/godirwalk"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/planar"
	"github.com/schollz/progressbar/v3"
)

// CheckDirFiles ...
func CheckDirFiles(dir, ext string) error {

	numFiles := 0
	progress := progressbar.Default(-1)

	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {
				fExt := filepath.Ext(path)
				if fExt != ext {
					return nil
				}
				numFiles++
				progress.Add(1)
				return nil
			}
			return nil
		},
		ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})
	if err != nil {
		return err
	}

	fmt.Println() // print new line after progress bar
	fmt.Printf("files found: %d\n", numFiles)
	return nil
}

// ReadFile ...
func ReadFile(path, fid string) (string, string, error) {

	// todo: check if path exists
	// !!! todo: check if fid exists in properties and type==string
	// todo: separate reading file and parsing id/bounds

	f, err := LoadFeature(path)
	if err != nil {
		return "", "", err
	}

	if _, ok := f.Properties[fid]; !ok {
		return "", "", err
	}

	id := f.Properties[fid].(string)

	bound := f.Geometry.Bound()
	boundStr := fmt.Sprintf("[%f %f],[%f %f]", bound.Min.X(), bound.Min.Y(), bound.Max.X(), bound.Max.Y())

	return id, boundStr, nil
}

// LoadFeature ...
func LoadFeature(path string) (*geojson.Feature, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	f, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil, err
	}
	return f, nil
}

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

// Bounds ...
func Bounds(o interface{}) string {

	var bounds string

	switch v := o.(type) {
	case orb.Point:
		bounds = fmt.Sprintf("[%f %f]", v.Lon(), v.Lat())
	case maptile.Tile:
		bounds = fmt.Sprintf("[%f %f], [%f %f]",
			v.Bound().Min.Lon(),
			v.Bound().Min.Lat(),
			v.Bound().Max.Lon(),
			v.Bound().Max.Lat(),
		)
	default:

	}

	return bounds
}

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

		ptInFeature := false

		geom := f.Geometry
		switch g := geom.(type) {
		case orb.Polygon:
			ptInFeature = planar.PolygonContains(g, pt)
		case orb.MultiPolygon:
			ptInFeature = planar.MultiPolygonContains(g, pt)
		default:
			continue
		}

		if ptInFeature {
			features = append(features, f)
		}
	}

	return features, nil
}

// ResolveTile ...
func ResolveTile(path, ext string, results []string, tile maptile.Tile) ([]*geojson.Feature, error) {

	features := []*geojson.Feature{}

	for _, r := range results {

		_, id, _ := db.ParseResult(r)

		fp := FilePath(path, id, ext)

		f, err := LoadFeature(fp)
		if err != nil {
			continue
		}

		features = append(features, f)
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
	fmt.Printf("%v\n", f)

	features = append(features, f)
	return features, nil
}

// FilePath ...
func FilePath(path, id, ext string) string {
	fn := fmt.Sprintf("%s%s", id, ext)
	fp := filepath.Join(path, fn)
	return fp
}
