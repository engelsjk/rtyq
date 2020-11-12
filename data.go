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

var (
	ErrUnableToOpenDataFile      error = fmt.Errorf("unable to open data file")
	ErrUnableToReadDataFile      error = fmt.Errorf("unable to read data file")
	ErrDoesNotMatchFileExtension error = fmt.Errorf("does match file extension")
	ErrInvalidGeoJSONFeature     error = fmt.Errorf("invalid geojson feature")
	ErrMissingFeatureID          error = fmt.Errorf("missing feature id")
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
func (d *Data) CheckDirFiles() (int, int, int, int, int, error) {

	numFiles := 0
	numFilesWithExtension := 0
	numReadableFiles := 0
	numFilesValidGeoJSON := 0
	numFilesWithID := 0

	progress := progressbar.Default(-1)

	err := godirwalk.Walk(d.DirPath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {

				numFiles++

				_, _, err := d.ReadFile(path)

				if err == ErrUnableToOpenDataFile || err == ErrUnableToReadDataFile {
					return err
				}
				numReadableFiles++

				if err == ErrDoesNotMatchFileExtension {
					return err
				}
				numFilesWithExtension++

				if err == ErrInvalidGeoJSONFeature {
					return err
				}
				numFilesValidGeoJSON++

				if err == ErrMissingFeatureID {
					return err
				}
				numFilesWithID++

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
		return 0, 0, 0, 0, 0, err
	}

	return numFiles, numFilesWithExtension, numReadableFiles, numFilesValidGeoJSON, numFilesWithID, nil
}

// ReadFile ...
func (d *Data) ReadFile(path string) (string, string, error) {

	ext := filepath.Ext(path)
	if ext != d.FileExtension {
		return "", "", ErrDoesNotMatchFileExtension
	}

	f, err := LoadFeature(path)
	if err != nil {
		return "", "", err
	}

	id, ok := f.Properties[d.ID]

	if !ok {
		return "", "", ErrMissingFeatureID
	}

	var idStr string
	switch v := id.(type) {
	case string:
		idStr = v
	case int:
		idStr = fmt.Sprintf("%d", v)
	case float64:
		idStr = fmt.Sprintf("%f", v)
	}

	bound := f.Geometry.Bound()
	boundStr := fmt.Sprintf("[%f %f],[%f %f]", bound.Min.X(), bound.Min.Y(), bound.Max.X(), bound.Max.Y())

	return idStr, boundStr, nil
}

// LoadFeature ...
func LoadFeature(path string) (*geojson.Feature, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ErrUnableToOpenDataFile
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, ErrUnableToReadDataFile
	}

	f, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil, ErrInvalidGeoJSONFeature
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
