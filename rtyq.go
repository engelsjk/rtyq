package main

import (
	"context"
	"fmt"
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
var checkFlagConfigFilename *string

var createCmd *kingpin.CmdClause
var createFlagConfigFilename *string

var startCmd *kingpin.CmdClause
var startFlagConfigFilename *string
var startFlagDebugOn *bool

var logger *zap.Logger

func initCommandOptions() {

	checkCmd = app.Command("check", "check data path")
	checkFlagConfigFilename = checkCmd.Flag("config", "config file").Short('c').String()

	createCmd = app.Command("create", "create an rtree db from data")
	createFlagConfigFilename = createCmd.Flag("config", "config file").Short('c').String()

	startCmd = app.Command("start", "start api service")
	startFlagConfigFilename = startCmd.Flag("config", "config file").Short('c').String()
	startFlagDebugOn = startCmd.Flag("debug", "enable debug").Short('d').Default("false").Bool()
}

func main() {

	// logger, _ = zap.NewDevelopment()
	// defer logger.Sync()

	app = kingpin.New(conf.AppConfig.Name, conf.AppConfig.Help).Version(conf.AppConfig.Version)

	initCommandOptions()

	// if flagDebugOn || conf.Configuration.Server.Debug {}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case checkCmd.FullCommand():
		conf.InitConfig(*checkFlagConfigFilename)
		check()
	case createCmd.FullCommand():
		conf.InitConfig(*createFlagConfigFilename)
		create()
	case startCmd.FullCommand():
		// what is happening here? why is the config filename flag empty?
		conf.InitConfig(*startFlagConfigFilename)
		start()
	}
}

func check() {
	for _, confLayer := range conf.Configuration.Layers {

		layer := data.NewLayer(confLayer)

		fmt.Printf("checking layer: %s\n", layer.Name)

		if err := layer.CheckData(); err != nil {
			// log error
			fmt.Println(err)
			continue
		}
	}
}

func create() {
	for _, confLayer := range conf.Configuration.Layers {

		layer := data.NewLayer(confLayer)

		fmt.Printf("creating layer: %s\n", layer.Name)

		if err := layer.CreateDatabase(); err != nil {
			// log error
			fmt.Println(err)
			continue
		}

		if err := layer.OpenDatabase(); err != nil {
			// log error
			fmt.Println(err)
			continue
		}

		if err := layer.AddDataToDatabase(); err != nil {
			// log error
			fmt.Println(err)
			continue
		}
	}
}

func start() {
	load()
	serve()
}

func load() {
	for _, confLayer := range conf.Configuration.Layers {
		layer := data.NewLayer(confLayer)
		
		fmt.Printf("loading layer: %s\n", layer.Name)

		// todo: check if data dir exists
		if err := layer.OpenDatabase(); err != nil {
			// log error
			fmt.Println(err)
			continue
		}
		if err := layer.IndexDatabase(); err != nil {
			// log error
			fmt.Println(err)
			continue
		}
		data.AddLayerToQueryHandler(layer)
	}
}

func serve() {

	srv := server.Create()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// log
		}
	}()

	fmt.Printf("listening at %s\n", srv.Addr)

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
