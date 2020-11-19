package rtyq

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

				progress.Add(1)
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

	fid, ok := f.Properties[d.ID]

	if !ok {
		return "", "", ErrMissingFeatureID
	}

	var id string

	if v, ok := fid.(string); ok {
		id = v
	}
	if v, ok := fid.(int); ok {
		id = strconv.FormatInt(int64(v), 10)
	}
	if v, ok := fid.(float64); ok {
		id = strconv.FormatFloat(v, 'f', -1, 64)
	}

	bounds := Bounds(f.Geometry)

	return id, bounds, nil
}

// LoadFeature opens, reads and unmarshals a GeoJSON Feature from the input filepath.
func LoadFeature(path string) (*geojson.Feature, error) {

	file, err := os.Open(path)

	if err != nil {
		if file != nil {
			file.Close()
		}
		return nil, err
	}

	buf := bytes.Buffer{}

	io.Copy(&buf, file)
	file.Close()

	b := buf.Bytes()

	f, err := geojson.UnmarshalFeature(b)
	if err != nil {
		return nil, ErrInvalidGeoJSONFeature
	}

	return f, nil
}

// FilePath returns the full filepath from a directory, file ID and file extension.
func FilePath(dir, id, ext string) string {
	var sb strings.Builder
	sb.WriteString(id)
	sb.WriteString(ext)
	fp := filepath.Join(dir, sb.String())
	return fp
}

// FileExists ...
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
