package main

import (
	"fmt"
	"os"

	"github.com/engelsjk/rtyq"
	"github.com/engelsjk/rtyq/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("rtyq", "generate and query spatial rtrees on disk").Version("0.0.1")

	check           = app.Command("check", "check data path")
	checkConfigFile = check.Flag("config", "config file").String()
	checkDataDir    = check.Flag("data", "data directory").Default(".").String()
	checkDataExt    = check.Flag("ext", "allowed file extension").Default(".geojson").String()
	checkDataID     = check.Flag("id", "object id").String()
	checkLayerName  = check.Flag("name", "name").String()

	create           = app.Command("create", "create an rtree db from data")
	createConfigFile = create.Flag("config", "config file").String()
	createDataDir    = create.Flag("data", "data directory").Default(".").String()
	createDataExt    = create.Flag("ext", "allowed file extension").Default(".geojson").String()
	createDataID     = create.Flag("id", "object id").String()
	createDBFile     = create.Flag("db", "database filepath").String()
	createIndex      = create.Flag("index", "index").String()
	createLayerName  = create.Flag("name", "name").String()

	start           = app.Command("start", "start api service")
	startConfigFile = start.Flag("config", "config file").String()
	startDataDir    = start.Flag("data", "data directory").Default(".").String()
	startDataExt    = start.Flag("ext", "allowed file extension").Default(".geojson").String()
	startDataID     = start.Flag("id", "unique identifier").String()
	startDBFile     = start.Flag("db", "database filepath").String()
	startIndex      = start.Flag("index", "index").String()
	startZoomLimit  = start.Flag("zoomlimit", "zoomlimit").Int()
	startEndpoint   = start.Flag("endpoint", "endpoint").String()
	startLayerName  = start.Flag("name", "name").String()
	startPort       = start.Flag("port", "api port").Default("5500").Int()
	startLogs       = start.Flag("logs", "enable logs").Default("false").Bool()
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case check.FullCommand():

		cfg, err := rtyq.LoadConfig(*checkConfigFile)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

		if cfg == nil {
			layer := rtyq.ConfigLayer{}
			layer.Name = *checkLayerName
			layer.Data.Path = *checkDataDir
			layer.Data.Extension = *checkDataExt
			layer.Data.ID = *checkDataID
			cfg = rtyq.NewConfig(layer)
		}

		err = rtyq.Check(cfg)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

	case create.FullCommand():

		cfg, err := rtyq.LoadConfig(*createConfigFile)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

		if cfg == nil {
			layer := rtyq.ConfigLayer{}
			layer.Name = *createLayerName
			layer.Data.Path = *createDataDir
			layer.Data.Extension = *createDataExt
			layer.Data.ID = *createDataID
			layer.Database.Path = *createDBFile
			layer.Database.Index = *createIndex
			cfg = rtyq.NewConfig(layer)
		}

		err = rtyq.Create(cfg)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

	case start.FullCommand():

		cfg, err := rtyq.LoadConfig(*startConfigFile)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}

		if cfg == nil {
			layer := rtyq.ConfigLayer{}
			layer.Data.Path = *startDataDir
			layer.Data.Extension = *startDataExt
			layer.Data.ID = *startDataID
			layer.Database.Path = *startDBFile
			layer.Database.Index = *startIndex
			layer.Service.ZoomLimit = *startZoomLimit
			layer.Service.Endpoint = *startEndpoint
			layer.Name = *startLayerName
			cfg = rtyq.NewConfig(layer)
			cfg.Port = *startPort
			cfg.EnableLogs = *startLogs
		}

		err = api.Start(cfg)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
	default:
		break
	}
}
