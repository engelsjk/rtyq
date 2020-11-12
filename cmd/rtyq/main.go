package main

import (
	"fmt"
	"os"

	"github.com/engelsjk/rtyq"
	"github.com/engelsjk/rtyq/api"
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
	serviceZoomLimit  = service.Flag("zoomlimit", "zoomlimit").Int()
	servicePort       = service.Flag("port", "api port").Default("5500").Int()
)

func main() {

	kingpin.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case check.FullCommand():

		cfg, err := rtyq.LoadConfig(*checkConfigFile)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		if cfg == nil {
			layer := rtyq.ConfigLayer{}
			layer.Data.Path = *checkDataDir
			layer.Data.Extension = *checkDataExt
			cfg = rtyq.NewConfig(layer)
		}

		err = rtyq.CheckData(cfg)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

	case create.FullCommand():

		cfg, err := rtyq.LoadConfig(*createConfigFile)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		if cfg == nil {
			layer := rtyq.ConfigLayer{}
			layer.Data.Path = *checkDataDir
			layer.Data.Extension = *checkDataExt
			layer.Data.ID = *createDataID
			layer.Database.Path = *createDBFile
			layer.Database.Index = *createIndex
			cfg = rtyq.NewConfig(layer)
		}

		err = rtyq.CreateDatabases(cfg)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

	case service.FullCommand():

		cfg, err := rtyq.LoadConfig(*serviceConfigFile)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			break
		}

		if cfg == nil {
			layer := rtyq.ConfigLayer{}
			layer.Data.Path = *serviceDataDir
			layer.Data.Extension = *serviceDataExt
			layer.Database.Path = *serviceDBFile
			layer.Database.Index = *serviceIndex
			layer.Service.ZoomLimit = *serviceZoomLimit
			cfg = rtyq.NewConfig(layer)
			cfg.Port = *servicePort
		}

		err = api.StartService(cfg)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	}
}
