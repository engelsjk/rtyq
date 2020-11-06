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

// ParseLonLat ...
func ParseLonLat(s string) []float64 {

	// todo: better latlon string validation

	cleanLatLon := strings.ReplaceAll(s, " ", "")
	splitLatLon := strings.Split(cleanLatLon, ",")

	lon, err := strconv.ParseFloat(splitLatLon[0], 64)
	if err != nil {
		return nil
	}

	lat, err := strconv.ParseFloat(splitLatLon[1], 64)
	if err != nil {
		return nil
	}

	return []float64{lon, lat}
}

// LonLat2Point ...
func LonLat2Point(s string) orb.Point {
	lonlat := ParseLonLat(s)
	pt := orb.Point{lonlat[0], lonlat[1]}
	return pt
}

// Bounds ...
func Bounds(pt []float64) string {
	bounds := fmt.Sprintf("[%f %f]", pt[0], pt[1])
	return bounds
}

// ResolveResults ...
func ResolveResults(path, ext string, results []string, lonlat string) ([]*geojson.Feature, error) {

	features := []*geojson.Feature{}

	pt := LonLat2Point(lonlat)

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

// FilePath ...
func FilePath(path, id, ext string) string {
	fn := fmt.Sprintf("%s%s", id, ext)
	fp := filepath.Join(path, fn)
	return fp
}
