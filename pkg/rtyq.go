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
func Check(path, ext string) error {
	return data.CheckDirFiles(path, ext)
}

// Create ...
func Create(dirData, ext, pathDB, index, fid string) error {

	bdb, err := db.Initialize(pathDB, index, true)
	if err != nil {
		return err
	}
	defer bdb.Close()

	n, err := Generate(bdb, dirData, ext, index, fid)
	if err != nil {
		return err
	}

	fmt.Printf("database created: %d data files processed\n", n)
	return nil
}

// Generate ...
func Generate(bdb *buntdb.DB, dir, ext, name, key string) (int, error) {

	progress := progressbar.Default(-1)
	numLoadErrors := 0
	numUpdateErrors := 0
	numFiles := 0

	start := time.Now()

	err := godirwalk.Walk(dir, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {
				fExt := filepath.Ext(path)
				if fExt != ext {
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

// Start ...
func Start(dirData, ext, pathDB, index string, port int) error {

	bdb, err := db.Initialize(pathDB, index, false)
	if err != nil {
		return err
	}

	return service.Start(port, dirData, ext, bdb, index)
}
