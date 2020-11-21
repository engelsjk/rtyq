package data

import (
	"fmt"
	"math"

	"github.com/engelsjk/rtyq/conf"
	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/buntdb"
)

var Layers map[string]*Layer

func init() {
	Layers = make(map[string]*Layer)
}

type Layer struct {
	Name       string
	DataDir    string
	DataExt    string
	DataID     string
	DBFilepath string
	DBIndex    string
	db         *buntdb.DB
}

func NewLayer(layer conf.Layer) *Layer {
	return &Layer{
		Name:       layer.Name,
		DataDir:    layer.Data.Dir,
		DataExt:    layer.Data.Ext,
		DataID:     layer.Data.ID,
		DBFilepath: layer.Database.Filepath,
		DBIndex:    layer.Database.Index,
	}
}

func (l *Layer) CheckData() error {

	if !dirExists(l.DataDir) {
		// log error
		return fmt.Errorf("data dir does not exist")
	}

	numFiles := 0
	var minFilesize int64 = math.MaxInt64
	var maxFilesize int64 = math.MinInt64

	fmt.Printf("layer (%s)\n", l.Name)
	fmt.Printf("checking data path...\n")

	progress := progressbar.Default(-1)

	err := godirwalk.Walk(l.DataDir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {
				progress.Add(1)

				if !validExtension(path, l.DataExt) {
					// log error
					return nil // or return error to skip?
				}

				nbytes, _, _, err := read(path, l.DataID)
				if err != nil {
					return err
				}

				numFiles++
				minFilesize = minBytes(minFilesize, nbytes)
				maxFilesize = maxBytes(maxFilesize, nbytes)
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
	fmt.Printf("largest: %d | smallest: %d\n", maxFilesize, minFilesize) // convert to KB (bytes/1024)

	return nil
}

func (l *Layer) CreateDatabase() error {
	if fileExists(l.DBFilepath) {
		return fmt.Errorf("database file already exists")
	}
	_, err := buntdb.Open(l.DBFilepath)
	if err != nil {
		return err
	}
	return nil
}

func (l *Layer) LoadDatabase() error {
	if !fileExists(l.DBFilepath) {
		return fmt.Errorf("database file does not exists")
	}
	l.db = nil
	bdb, err := buntdb.Open(l.DBFilepath)
	if err != nil {
		return err
	}
	l.db = bdb
	return nil
}

func (l *Layer) AddDataToDatabase() error {

	if !dirExists(l.DataDir) {
		// log error
		return fmt.Errorf("data dir does not exist")
	}
	if !fileExists(l.DBFilepath) {
		// log error
		return fmt.Errorf("database file does not exist")
	}
	if l.db == nil {
		// log error
		return fmt.Errorf("database not loaded")
	}

	fmt.Printf("layer (%s)\n", l.Name)
	fmt.Printf("uploading data to database...\n")

	numLoadErrors := 0
	numUpdateErrors := 0
	numFiles := 0

	progress := progressbar.Default(-1)

	err := godirwalk.Walk(l.DataDir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {

				progress.Add(1)

				_, id, bound, err := read(path, l.DataID)
				if err != nil {
					numLoadErrors++
					return err
				}

				err = dbUpdate(l.db, l.DBIndex, id, bound)
				if err != nil {
					numUpdateErrors++
					return err
				}

				numFiles++
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
	if numLoadErrors > 0 || numUpdateErrors > 0 {
		fmt.Printf("warning: %d load errors | %d update errors\n", numLoadErrors, numUpdateErrors)
	}
	fmt.Printf("num files loaded to db: %d\n", numFiles)
	return nil
}

func AddToLayers(layer *Layer) {
	Layers[layer.Name] = layer
}
