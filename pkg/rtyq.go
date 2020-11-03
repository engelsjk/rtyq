package rtyq

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/engelsjk/rtyq/pkg/data"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/engelsjk/rtyq/pkg/service"
	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/buntdb"
)

// Check ...
func Check(path string) error {
	return data.CheckDirFiles(path)
}

// Create ...
func Create(pathDB, dirData, index, fid string) error {

	bdb, err := db.Initialize(pathDB, index, true)
	if err != nil {
		return err
	}
	defer bdb.Close()

	n, err := Generate(bdb, dirData, index, fid)
	if err != nil {
		return err
	}

	fmt.Printf("database created: %d data files processed\n", n)
	return nil
}

// Generate ...
func Generate(bdb *buntdb.DB, dir, name, key string) (int, error) {

	progress := progressbar.Default(-1)
	numLoadErrors := 0
	numUpdateErrors := 0
	numFiles := 0

	start := time.Now()

	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {
				ext := filepath.Ext(path)
				if ext != ".geojson" {
					return nil
				}
				id, bound, err := data.ReadFile(path, key)
				if err != nil {
					numLoadErrors++
					return err
				}
				err = db.Update(bdb, name, id, bound)
				if err != nil {
					numUpdateErrors++
					return err
				}
				numFiles++
				progress.Add(1)
			}
			return nil
		},
		// ErrorCallback: ???
	})
	if err != nil {
		return 0, err
	}

	fmt.Println() // print new line after progress bar

	dur := time.Since(start)
	fmt.Printf("time to generate db: %s sec\n", dur)

	if numLoadErrors > 0 || numUpdateErrors > 0 {
		fmt.Printf("warning: %d load errors | %d update errors", numLoadErrors, numUpdateErrors)
	}

	return numFiles, nil
}

// QueryCold ...
func QueryCold(pathDB, dirData, index, lonlat string, geojson bool) error {

	// todo: add input validation w/ errors

	pt := data.ParseLonLat(lonlat) // pt = [lon lat]

	bdb, err := db.Initialize(pathDB, index, false)
	if err != nil {
		return err
	}
	defer bdb.Close()

	results, err := db.Get(bdb, index, data.Bounds(pt))
	if err != nil {
		return err
	}

	features, err := data.ResolveResults(dirData, results, lonlat)
	if err != nil {
		return err
	}

	for _, f := range features {
		b, err := f.MarshalJSON()
		if err != nil {
			continue
		}
		fmt.Println(string(b))
	}

	return nil
}

// Start ...
func Start(pathDB, dirData, index string, port int) error {

	bdb, err := db.Initialize(pathDB, index, false)
	if err != nil {
		return err
	}

	err = service.Start(port, bdb, dirData, index)
	if err != nil {
		return err
	}

	return nil
}
