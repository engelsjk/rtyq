package main

import (
	"fmt"
	"os"

	rtyq "github.com/engelsjk/rtyq/pkg"
	"github.com/engelsjk/rtyq/pkg/config"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("rtyq", "generate and query spatial rtrees on disk")

	check           = app.Command("check", "check data path")
	checkConfigFile = check.Flag("config", "config file").String()
	checkDataDir    = check.Flag("data", "data directory").Default(".").String()
	checkDataExt    = check.Flag("ext", "allowed file extension").Default(".geojson").String()

	create           = app.Command("create", "create an rtree db from data")
	createConfigFile = create.Flag("config", "config file").String()
	createDataDir    = create.Flag("data", "data directory").Default(".").String()
	createDataExt    = create.Flag("ext", "allowed file extension").Default(".geojson").String()
	createDataID     = create.Flag("id", "object id").Default("id").String()
	createDBFile     = create.Flag("db", "database filepath").Default("data.db").String()
	createIndex      = create.Flag("index", "index").Default("data").String()

	service           = app.Command("service", "start api service")
	serviceConfigFile = service.Flag("config", "config file").String()
	serviceDataDir    = service.Flag("data", "data directory").Default(".").String()
	serviceDataExt    = service.Flag("ext", "allowed file extension").Default(".geojson").String()
	serviceDBFile     = service.Flag("db", "database filepath").Default("data.db").String()
	serviceIndex      = service.Flag("index", "index").Default("data").String()
	servicePort       = service.Flag("port", "api port").Default("5500").Int()
)

func main() {

	kingpin.Version("0.0.1")

	var err error

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case check.FullCommand():

		cfg, err := config.Load(*checkConfigFile)
		if err != nil {
			fmt.Printf("errors: %s\n", err.Error())
			break
		}

		if cfg == nil {
			svc := config.Service{}
			svc.Data.Path = *checkDataDir
			svc.Data.Extension = *checkDataExt
			cfg = config.Create(svc)
		}

		err = rtyq.Check(cfg)

	case create.FullCommand():

		cfg, err := config.Load(*createConfigFile)
		if err != nil {
			fmt.Printf("errors: %s\n", err.Error())
			break
		}

		if cfg == nil {
			svc := config.Service{}
			svc.Data.Path = *createDataDir
			svc.Data.Extension = *createDataExt
			svc.Data.ID = *createDataID
			svc.Database.Path = *createDBFile
			svc.Database.Index = *createIndex
			cfg = config.Create(svc)
		}

		err = rtyq.Create(cfg)

	case service.FullCommand():

		cfg, err := config.Load(*serviceConfigFile)
		if err != nil {
			fmt.Printf("errors: %s\n", err.Error())
			break
		}

		if cfg == nil {
			svc := config.Service{}
			svc.Data.Path = *serviceDataDir
			svc.Data.Extension = *serviceDataExt
			svc.Database.Path = *serviceDBFile
			svc.Database.Index = *serviceIndex
			cfg = config.Create(svc)
			cfg.Port = *servicePort
		}

		err = rtyq.Service(cfg)
	}

	if err != nil {
		fmt.Printf("errors: %s\n", err.Error())
		return
	}
}
