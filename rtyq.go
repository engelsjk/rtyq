package rtyq

import (
	"fmt"
	"time"

	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
)

var (
	ErrNoConfigProvided error = fmt.Errorf("no config provided")
)

// Check ...
func Check(cfg *Config) error {

	err := ValidateConfigData(cfg)
	if err != nil {
		return err
	}

	for _, layer := range cfg.Layers {

		d, err := InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
		if err != nil {
			return err
		}

		fmt.Printf("layer (%s)\n", layer.Name)
		fmt.Printf("checking data path...\n")

		numFiles, numFilesWithExtension, numReadableFiles, numFilesValidGeoJSON, numFilesWithID, err := d.CheckDirFiles()
		if err != nil {
			return err
		}

		fmt.Println() // print new line after progress bar
		fmt.Printf("files found: %d\n", numFiles)
		fmt.Printf("files w/ extension (%s): %d\n", layer.Data.Extension, numFilesWithExtension)
		fmt.Printf("files readable: %d\n", numReadableFiles)
		fmt.Printf("files w/ valid geojson feature: %d\n", numFilesValidGeoJSON)
		fmt.Printf("files w/ property id (%s): %d\n", layer.Data.ID, numFilesWithID)
	}
	return nil
}

// Create ...
func Create(cfg *Config) error {

	err := ValidateConfigData(cfg)
	if err != nil {
		return err
	}

	err = ValidateConfigDatabase(cfg)
	if err != nil {
		return err
	}

	for _, layer := range cfg.Layers {

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

		fmt.Printf("generating database: %s\n", db.FileName)

		_, err = AddDataToDatabaseWithIndex(data, db, layer.Database.Index)
		if err != nil {
			fmt.Printf("error : layer (%s) : %s\n", layer.Name, err.Error())
			continue
		}
	}

	return nil
}

// AddDataToDatabaseWithIndex ...
func AddDataToDatabaseWithIndex(data *Data, db *DB, index string) (int, error) {

	numLoadErrors := 0
	numUpdateErrors := 0
	numFiles := 0

	db.Index = index

	fmt.Printf("adding data with db index:%s...\n", db.Index)

	progress := progressbar.Default(-1)

	start := time.Now()

	err := godirwalk.Walk(data.DirPath, &godirwalk.Options{
		Unsorted: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if de.ModeType().IsRegular() {

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
		ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})

	if err != nil {
		return 0, err
	}

	dur := time.Since(start)

	fmt.Println() // print new line after progress bar
	fmt.Printf("time to generate db: %s (%d files)\n", dur.String(), numFiles)

	if numLoadErrors > 0 || numUpdateErrors > 0 {
		fmt.Printf("warning: %d load errors | %d update errors\n", numLoadErrors, numUpdateErrors)
	}

	return numFiles, nil
}
