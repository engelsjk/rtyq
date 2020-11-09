package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
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

// FilePath ...
func FilePath(path, id, ext string) string {
	fn := fmt.Sprintf("%s%s", id, ext)
	fp := filepath.Join(path, fn)
	return fp
}
