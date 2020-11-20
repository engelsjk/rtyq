package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/engelsjk/rtyq/conf"
	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/buntdb"
)

type Layer struct {
	Name     string
	Dir      string
	Ext      string
	ID       string
	Filepath string
	Index    string
	db       *buntdb.DB
}

type layerError struct {
	Error   error
	Type    string
	Message string
}

func CreateLayer(layer conf.Layer) (*Layer, error) {

	if _, err := os.Stat(layer.Data.Dir); os.IsNotExist(err) {
		return nil, layerErrorDirNotFound(err).Error
	}

	return &Layer{
		Name:     layer.Name,
		Dir:      layer.Data.Dir,
		Ext:      layer.Data.Ext,
		ID:       layer.Data.ID,
		Filepath: layer.Database.Filepath,
		Index:    layer.Database.Index,
	}, nil
}

func (l Layer) CheckData() error {

	numFiles := 0
	numReadableFiles := 0
	numFilesWithExtension := 0
	numFilesWithID := 0

	fmt.Printf("layer (%s)\n", l.Name)
	fmt.Printf("checking data path...\n")

	progress := progressbar.Default(-1)

	err := godirwalk.Walk(l.Dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {

				progress.Add(1)
				numFiles++

				_, _, err := l.readFile(path)
				if err == nil {
					numFilesWithExtension++
					numReadableFiles++
					numFilesWithID++
					return nil
				}

				if err.Type == "err_invalid_file_extension" {
					return err.Error
				}
				numFilesWithExtension++

				if err.Type == "error_load_feature" {
					return err.Error
				}
				numReadableFiles++

				if err.Type == "error_no_id" {
					return err.Error
				}
				numFilesWithID++

				if err.Type != "" {
					return err.Error
				}

				return err.Error
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
	fmt.Printf("files w/ extension (%s): %d\n", l.Ext, numFilesWithExtension)
	fmt.Printf("files readable: %d\n", numReadableFiles)
	fmt.Printf("files w/ property id (%s): %d\n", l.ID, numFilesWithID)

	return nil
}

func (l Layer) readFile(path string) (string, string, *layerError) {

	ext := filepath.Ext(path)
	if ext != l.Ext {
		return "", "", layerErrorInvalidFileExtension(nil)
	}

	f, err := loadFeature(path)
	if err != nil {
		return "", "", layerErrorLoadFeature(err)
	}

	fid, ok := f.Properties[l.ID]

	if !ok {
		return "", "", layerErrorNoID(err)
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

	bounds := bounds(f.Geometry)

	return id, bounds, nil
}

func layerErrorDirNotFound(err error) *layerError {
	return &layerError{err, "err_dir_not_found", "data dir not found"}
}

func layerErrorInvalidFileExtension(err error) *layerError {
	return &layerError{err, "err_invalid_file_extension", "file has invalid extension"}
}

func layerErrorLoadFeature(err error) *layerError {
	return &layerError{err, "err_load_feature", "unable to load data file"}
}

func layerErrorNoID(err error) *layerError {
	return &layerError{err, "error_no_id", "no feature id"}
}
