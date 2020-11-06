package rtyq

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/engelsjk/rtyq/pkg/config"
	"github.com/engelsjk/rtyq/pkg/data"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/engelsjk/rtyq/pkg/service"
	"github.com/karrick/godirwalk"
	"github.com/schollz/progressbar/v3"
	"github.com/tidwall/buntdb"
)

// Check ...
func Check(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("no config provided")
	}

	for _, svc := range cfg.Services {

		fmt.Printf("checking data path: %s...", filepath.Base(svc.Data.Path))

		err := data.CheckDirFiles(svc.Data.Path, svc.Data.Extension)
		if err != nil {
			return err
		}
	}
	return nil
}

// Create ...
func Create(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("no config provided")
	}

	for _, svc := range cfg.Services {

		fmt.Printf("generating database: %s...\n", filepath.Base(svc.Database.Path))

		bdb, err := db.Initialize(svc.Database.Path, svc.Database.Index, true)
		if err != nil {
			return err
		}
		defer bdb.Close()

		_, err = generate(bdb, svc.Data.Path, svc.Data.Extension, svc.Database.Index, svc.Data.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Service ...
func Service(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("no config provided")
	}
	return service.Start(cfg)
}

// generate ...
func generate(bdb *buntdb.DB, dir, ext, name, key string) (int, error) {

	//todo: update numLoadErrors and numUpdateErrors via the ErrorCallback function

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
	fmt.Printf("time to generate db: %s sec (%d files)\n", dur.String(), numFiles)

	if numLoadErrors > 0 || numUpdateErrors > 0 {
		fmt.Printf("warning: %d load errors | %d update errors", numLoadErrors, numUpdateErrors)
	}

	return numFiles, nil
}
