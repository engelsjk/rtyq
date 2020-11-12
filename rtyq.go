package rtyq

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
)

var (
	ErrNoConfigProvided error = fmt.Errorf("no config provided")
)

// CheckData ...
func CheckData(cfg *Config) error {

	if cfg == nil {
		return ErrNoConfigProvided
	}

	for _, layer := range cfg.Layers {

		d, err := InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			return err
		}

		fmt.Printf("checking data path: %s...", filepath.Base(layer.Data.Path))

		err = d.CheckDirFiles()
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDatabases ...
func CreateDatabases(cfg *Config) error {

	if cfg == nil {
		return ErrNoConfigProvided
	}

	for _, layer := range cfg.Layers {

		// todo: print errors and continue instead?

		db, err := NewDB(layer.Database.Path)
		if err != nil {
			fmt.Printf("error : layer (%s) : %s\n", layer.Name, err.Error())
			continue
		}

		data, err := InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			fmt.Printf("error : layer (%s) : %s\n", layer.Name, err.Error())
			continue
		}

		fmt.Printf("generating database: %s...\n", db.FileName)

		_, err = AddDataToDatabaseWithIndex(data, db, layer.Database.Index)
		if err != nil {
			fmt.Printf("error : layer (%s) : %s\n", layer.Name, err.Error())
			continue
		}
	}

	return nil
}

// AddDataToDatabaseIndex ...
func AddDataToDatabaseWithIndex(data *Data, db *DB, index string) (int, error) {

	// todo: update numLoadErrors and numUpdateErrors via the ErrorCallback function
	// warning below is never triggered since these vars aren't updating properly

	numLoadErrors := 0
	numUpdateErrors := 0
	numFiles := 0

	db.Index = index

	fmt.Printf("generating db data with index:%s (%s)...\n", db.Index, db.FileName)

	progress := progressbar.Default(-1)

	start := time.Now()

	err := godirwalk.Walk(data.DirPath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {
				ext := filepath.Ext(path)
				if ext != data.FileExtension {
					return nil
				}
				id, bound, err := data.ReadFile(path)
				if err != nil {
					numLoadErrors++
					return err
				}
				err = db.Update(id, bound)
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

	dur := time.Since(start)

	fmt.Println() // print new line after progress bar
	fmt.Printf("time to generate db: %s (%d files)\n", dur.String(), numFiles)

	if numLoadErrors > 0 || numUpdateErrors > 0 {
		fmt.Printf("warning: %d load errors | %d update errors", numLoadErrors, numUpdateErrors)
	}

	return numFiles, nil
}
