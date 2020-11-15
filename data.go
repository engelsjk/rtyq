package rtyq

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/karrick/godirwalk"
	"github.com/paulmach/orb/geojson"
	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnableToOpenDataFile      error = fmt.Errorf("unable to open data file")
	ErrUnableToCloseDataFile     error = fmt.Errorf("unable to close data file")
	ErrDoesNotMatchFileExtension error = fmt.Errorf("does match file extension")
	ErrInvalidGeoJSONFeature     error = fmt.Errorf("invalid geojson feature")
	ErrMissingFeatureID          error = fmt.Errorf("missing feature id")
)

// Data is a structure that contains information
// about the GeoJSON data for a single layer
// (DirPath, FileExtension and ID)
type Data struct {
	DirPath       string
	FileExtension string
	ID            string
}

// InitData initializes a Data structure, first checking that
// the data path exists
func InitData(path, ext, id string) (Data, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Data{}, fmt.Errorf("data dir path does not exist")
	}

	return Data{
		DirPath:       path,
		FileExtension: ext,
		ID:            id,
	}, nil
}

// CheckDirFiles iterates through a data path to check the validity of data files.
// It outputs metrics on the number of files, if they're readable, if they match the
// the specified file extension, if they're valid GeoJSON Features and if the Features
// include the specified ID property.
func (d Data) CheckDirFiles() (int, int, int, int, int, error) {

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

				if err == ErrUnableToOpenDataFile {
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

// ReadFile loads a GeoJSON Feature from the input filepath.
func (d Data) ReadFile(path string) (string, string, error) {

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

	// todo: remove this switch type and error if id type != string

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

// LoadFeature opens, reads and unmarshals a GeoJSON Feature from the input filepath.
func LoadFeature(path string) (*geojson.Feature, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			log.Printf("%s", ErrUnableToCloseDataFile.Error())
		}
	}(file)

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	b := buf.Bytes()

	f, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil, ErrInvalidGeoJSONFeature
	}

	return f, nil
}

// FilePath returns the full filepath from a directory, file ID and file extension.
func FilePath(dir, id, ext string) string {
	fn := fmt.Sprintf("%s%s", id, ext)
	fp := filepath.Join(dir, fn)
	return fp
}

// FileExists ...
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
