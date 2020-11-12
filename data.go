package rtyq

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/paulmach/orb/geojson"
	"github.com/schollz/progressbar/v3"
)

// Data ...
type Data struct {
	DirPath       string
	FileExtension string
	ID            string
}

// InitData ...
func InitData(path, ext, id string) (*Data, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("data dir path does not exist")
	}

	return &Data{
		DirPath:       path,
		FileExtension: ext,
		ID:            id,
	}, nil
}

// CheckDirFiles ...
func (d *Data) CheckDirFiles() error {

	numFiles := 0
	progress := progressbar.Default(-1)

	err := godirwalk.Walk(d.DirPath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {
				fExt := filepath.Ext(path)
				if fExt != d.FileExtension {
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
func (d *Data) ReadFile(path string) (string, string, error) {

	// todo: check if path exists
	// !!! todo: check if fid exists in properties and type==string
	// todo: separate reading file and parsing id/bounds

	f, err := LoadFeature(path)
	if err != nil {
		return "", "", err
	}

	if _, ok := f.Properties[d.ID]; !ok {
		return "", "", err
	}

	id := f.Properties[d.ID].(string)

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

// FeaturesToString ...
func FeaturesToString(features []*geojson.Feature) string {

	if features == nil {
		return "[]"
	}

	featuresStr := []string{}

	for _, f := range features {
		b, err := f.MarshalJSON()
		if err != nil {
			continue
		}
		featuresStr = append(featuresStr, string(b))
	}

	out := fmt.Sprintf("[%s]", strings.Join(featuresStr, ","))
	return out
}

// FilePath ...
func FilePath(path, id, ext string) string {
	fn := fmt.Sprintf("%s%s", id, ext)
	fp := filepath.Join(path, fn)
	return fp
}
