package data

import (
	"fmt"
	"math"

	"github.com/engelsjk/rtyq/conf"
	"github.com/karrick/godirwalk"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/buntdb"
)

type Layer struct {
	Name       string
	DataDir    string
	DataExt    string
	DataID     string
	DBFilepath string
	DBIndex    string
	ZoomLimit  int
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
		ZoomLimit:  layer.ZoomLimit,
	}
}

//////////////////////////////////////////////////////////

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

//////////////////////////////////////////////////////////

func (l *Layer) CreateDatabase() error {

	if fileExists(l.DBFilepath) {
		return fmt.Errorf("database file already exists")
	}

	// log
	fmt.Printf("creating db %s...\n", filename(l.DBFilepath))

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

	// log
	fmt.Printf("loading db %s...\n", filename(l.DBFilepath))

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

	// log
	fmt.Printf("uploading data to db %s...\n", filename(l.DBFilepath))

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
	fmt.Printf("%d files loaded to db: %d\n", numFiles, filename(l.DBFilepath))
	return nil
}

func (l *Layer) IndexDatabase() error {
	// log
	fmt.Printf("indexing db %s...\n", filename(l.DBFilepath))
	return l.db.CreateSpatialIndex(l.DBIndex, dbPattern(l.DBIndex), buntdb.IndexRect)
}

func (l *Layer) intersects(o interface{}) ([]geojson.Feature, error) {

	var features []geojson.Feature

	if err := l.db.View(func(tx *buntdb.Tx) error {
		tx.Intersects(l.DBIndex, bounds(o), func(k, v string) bool {
			f := resolve(l, k, o)
			if f != nil {
				features = append(features, *f)
			}
			return true
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return features, nil
}

func AddLayerToQueryHandler(layer *Layer) {
	QueryHandler.layers[layer.Name] = layer
}

func resolve(layer *Layer, k string, o interface{}) *geojson.Feature {

	index, id := dbParseKey(k)

	if index != layer.DBIndex {
		return nil
	}

	fp := filePath(layer.DataDir, id, layer.DataExt)

	f, _, err := feature(fp)
	if err != nil {
		return nil
	}

	switch v := o.(type) {
	case orb.Point:
		if pointInFeature(f.Geometry, v) {
			return f
		}
	case maptile.Tile:
		if tileOverlapsGeometry(f.Geometry, v) {
			return f
		}
	default:
		return nil
	}

	return nil
}
