package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/engelsjk/rtyq/conf"
	"github.com/engelsjk/rtyq/data"
	"github.com/engelsjk/rtyq/server"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
)

var app *kingpin.Application
var checkCmd *kingpin.CmdClause
var createCmd *kingpin.CmdClause
var startCmd *kingpin.CmdClause

var logger *zap.Logger

var flagConfigFilename string
var flagDebugOn bool

func initCommandOptions() {

	checkCmd = app.Command("check", "check data path")
	flagConfigFilename = *checkCmd.Flag("config", "config file").Short('c').String()

	createCmd = app.Command("create", "create an rtree db from data")
	flagConfigFilename = *createCmd.Flag("config", "config file").Short('c').String()

	startCmd = app.Command("start", "start api service")
	flagConfigFilename = *startCmd.Flag("config", "config file").Short('c').String()
	flagDebugOn = *startCmd.Flag("debug", "enable debug").Short('d').Default("false").Bool()

}

func main() {

	// logger, _ = zap.NewDevelopment()
	// defer logger.Sync()

	app = kingpin.New(conf.AppConfig.Name, conf.AppConfig.Help).Version(conf.AppConfig.Version)

	initCommandOptions()

	conf.InitConfig("config.json")

	// if flagDebugOn || conf.Configuration.Server.Debug {}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case checkCmd.FullCommand():
		check()
	case createCmd.FullCommand():
		create()
	case startCmd.FullCommand():
		start()
	}
}

func check() {

	for _, confLayer := range conf.Configuration.Layers {

		layer, err := data.CreateLayer(confLayer)
		if err != nil {
			panic(err)
		}

		err = layer.CheckData()
		if err != nil {
			panic(err)
		}
	}
}

func create() {}

func start() {}

func load() {}

func serve() {

	srv := server.Create()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// log shut down

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	abortTimeoutSec := conf.Configuration.Server.WriteTimeoutSec + 10
	chanCancelFatal := server.FatalAfter(abortTimeoutSec, "timeout on shutdown - aborting")

	// log server stopped
	close(chanCancelFatal)
}

// // Create initializes a database file
// // and loads all data from a directory
// // for each specified layer
// func Create(cfg Config) error {

// 	err := ValidateConfigData(cfg)
// 	if err != nil {
// 		return err
// 	}

// 	err = ValidateConfigDatabase(cfg)
// 	if err != nil {
// 		return err
// 	}

// 	for _, layer := range cfg.Layers {

// 		dbFilename := filepath.Base(filepath.Base(layer.Database.Path))

// 		fmt.Println("%************%")
// 		fmt.Printf("creating layer: %s\n", layer.Name)

// 		if FileExists(layer.Database.Path) {
// 			fmt.Printf("warning : layer (%s) : %s (%s) : skipping layer\n", layer.Name, ErrDatabaseFileAlreadyExists.Error(), dbFilename)
// 			continue
// 		}

// 		fmt.Printf("initializing database\n")

// 		db, err := InitDB(layer.Database.Path)
// 		if err != nil {
// 			fmt.Printf("warning : layer (%s) : %s (%s) : skipping layer\n", layer.Name, err.Error(), dbFilename)
// 			continue
// 		}

// 		data, err := InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
// 		if err != nil {
// 			fmt.Printf("warning : layer (%s) : %s\n", layer.Name, err.Error())
// 			continue
// 		}

// 		fmt.Printf("adding data to %s with index:%s...\n", dbFilename, layer.Database.Index)
// 		start := time.Now()

// 		numFiles, err := AddDataToDatabaseWithIndex(data, db, layer.Database.Index)
// 		if err != nil {
// 			fmt.Printf("warning : layer (%s) : %s (%s) : skipping layer\n", layer.Name, err.Error(), dbFilename)
// 			continue
// 		}

// 		dur := time.Since(start)

// 		fmt.Printf("time to generate db: %s (added %d files)\n", dur.String(), numFiles)
// 	}

// 	return nil
// }

// // AddDataToDatabaseWithIndex adds data from a data directory
// // to a database file using the specified index
// func AddDataToDatabaseWithIndex(data Data, db DB, index string) (int, error) {

// 	numLoadErrors := 0
// 	numUpdateErrors := 0
// 	numFiles := 0

// 	db.Index = index

// 	progress := progressbar.Default(-1)

// 	err := godirwalk.Walk(data.DirPath, &godirwalk.Options{
// 		Unsorted: true,
// 		Callback: func(path string, de *godirwalk.Dirent) error {
// 			if de.ModeType().IsRegular() {

// 				progress.Add(1)

// 				id, bound, err := data.ReadFile(path)
// 				if err != nil {
// 					numLoadErrors++
// 					return err
// 				}

// 				err = db.Update(id, bound)
// 				if err != nil {
// 					numUpdateErrors++
// 					return err
// 				}

// 				numFiles++
// 			}
// 			return nil
// 		},
// 		ErrorCallback: func(path string, err error) godirwalk.ErrorAction {
// 			return godirwalk.SkipNode
// 		},
// 	})

// 	if err != nil {
// 		return 0, err
// 	}

// 	fmt.Println() // print new line after progress bar

// 	if numLoadErrors > 0 || numUpdateErrors > 0 {
// 		fmt.Printf("warning: %d load errors | %d update errors\n", numLoadErrors, numUpdateErrors)
// 	}

// 	return numFiles, nil
// }
