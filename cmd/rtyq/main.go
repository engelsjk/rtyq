package main

import (
	"fmt"
	"os"

	rtyq "github.com/engelsjk/rtyq/pkg"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("rtyq", "generate and query spatial rtrees on disk")

	check        = app.Command("check", "check data path")
	checkDataDir = check.Flag("data", "data directory").Default(".").String()
	checkDataExt = check.Flag("ext", "allowed file extension").Default(".geojson").String()

	create        = app.Command("create", "create an rtree db from data")
	createDataDir = create.Flag("data", "data directory").Default(".").String()
	createDataExt = create.Flag("ext", "allowed file extension").Default(".geojson").String()
	createDBFile  = create.Flag("db", "database filepath").Default("data.db").String()
	createIndex   = create.Flag("index", "index").Default("data").String()
	createID      = create.Flag("id", "object id").Default("id").String()

	query        = app.Command("query", "query by lon,lat")
	queryDataDir = query.Flag("data", "data directory").Default(".").String()
	queryDataExt = query.Flag("ext", "allowed file extension").Default(".geojson").String()
	queryDBFile  = query.Flag("db", "database filepath").Default("data.db").String()
	queryIndex   = query.Flag("index", "index").Default("data").String()
	queryLonLat  = query.Flag("lonlat", "lon,lat").Required().String()
	queryGeoJSON = query.Flag("geojson", "geojson output").Default("false").Bool()

	service        = app.Command("service", "start api service")
	serviceDataDir = service.Flag("data", "data directory").Default(".").String()
	serviceDataExt = service.Flag("ext", "allowed file extension").Default(".geojson").String()
	serviceDBFile  = service.Flag("db", "database filepath").Default("data.db").String()
	serviceIndex   = service.Flag("index", "index").Default("data").String()
	servicePort    = service.Flag("port", "api port").Default("5500").Int()
)

func main() {

	kingpin.Version("0.0.1")

	var err error

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case create.FullCommand():
		err = rtyq.Create(
			*createDataDir,
			*createDataExt,
			*createDBFile,
			*createIndex,
			*createID,
		)
	case check.FullCommand():
		err = rtyq.Check(
			*checkDataDir,
			*checkDataExt,
		)
	case query.FullCommand():
		err = rtyq.QueryCold(
			*queryDataDir,
			*queryDataExt,
			*queryDBFile,
			*queryIndex,
			*queryLonLat,
			*queryGeoJSON,
		)
	case service.FullCommand():
		err = rtyq.Start(
			*serviceDataDir,
			*serviceDataExt,
			*serviceDBFile,
			*serviceIndex,
			*servicePort,
		)
	}

	if err != nil {
		fmt.Printf("errors: %s\n", err.Error())
		return
	}
}
