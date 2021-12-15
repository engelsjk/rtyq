package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/engelsjk/rtyq/conf"
	"github.com/engelsjk/rtyq/data"
	"github.com/engelsjk/rtyq/server"
)

func main() {

	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	checkFlagConfigFilename := checkCmd.String("config", "config.json", "config file")

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	createFlagConfigFilename := createCmd.String("config", "config.json", "config file")

	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	startFlagConfigFilename := startCmd.String("config", "config.json", "config file")

	if len(os.Args) < 2 {
		fmt.Println("expected 'check', 'create' or 'start' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "check":
		checkCmd.Parse(os.Args[2:])
		conf.InitConfig(*checkFlagConfigFilename)
		check()
	case "create":
		createCmd.Parse(os.Args[2:])
		conf.InitConfig(*createFlagConfigFilename)
		create()
	case "start":
		startCmd.Parse(os.Args[2:])
		conf.InitConfig(*startFlagConfigFilename)
		start()
	default:
		fmt.Println("expected 'check', 'create' or 'start' subcommands")
		os.Exit(1)
	}
}

func check() {
	for _, confLayer := range conf.Configuration.Layers {
		layer := data.NewLayer(confLayer)
		log.Printf("checking layer: %s\n", layer.Name)
		if err := layer.CheckData(); err != nil {
			log.Println(err)
			continue
		}
	}
}

func create() {
	for _, confLayer := range conf.Configuration.Layers {
		layer := data.NewLayer(confLayer)
		log.Printf("creating layer: %s\n", layer.Name)
		if err := layer.CreateDatabase(); err != nil {
			log.Println(err)
			continue
		}
		if err := layer.OpenDatabase(); err != nil {
			log.Println(err)
			continue
		}
		if err := layer.AddDataToDatabase(); err != nil {
			log.Println(err)
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

		log.Printf("loading layer: %s\n", layer.Name)

		if err := layer.OpenDatabase(); err != nil {
			log.Println(err)
			continue
		}
		if err := layer.IndexDatabase(); err != nil {
			log.Println(err)
			continue
		}
		data.AddLayerToQueryHandler(layer)
	}
}

func serve() {

	srv := server.Create()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}
	}()

	log.Printf("listening at http://%s\n", srv.Addr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	abortTimeoutSec := conf.Configuration.Server.WriteTimeoutSec + 10
	chanCancelFatal := server.FatalAfter(abortTimeoutSec, "timeout on shutdown - aborting")

	close(chanCancelFatal)
}
