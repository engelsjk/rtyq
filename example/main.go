package main

import (
	"fmt"

	"github.com/engelsjk/rtyq"
	"github.com/engelsjk/rtyq/query"
)

func main() {

	layer := rtyq.ConfigLayer{}
	layer.Data.Path = "path/to/data"
	layer.Data.Extension = ".geojson"
	layer.Data.ID = "FIPS"
	layer.Database.Path = "path/to/database.db"
	layer.Database.Index = "state"

	/////////////////////////////////////////////////////////////////////////

	db, err := rtyq.InitDB(layer.Database.Path)
	if err != nil {
		panic(err)
	}

	data, err := rtyq.InitData(layer.Data.Path, layer.Data.Extension, layer.Data.ID)
	if err != nil {
		panic(err)
	}

	_, err = rtyq.AddDataToDatabaseWithIndex(data, db, layer.Database.Index)
	if err != nil {
		panic(err)
	}

	err = db.CreateSpatialIndex(layer.Database.Index)
	if err != nil {
		panic(err)
	}

	/////////////////////////////////////////////////////////////////////////

	pt := "-86.46283149719237,32.470450258108315"

	point, err := query.ParsePoint(pt)
	if err != nil {
		panic(err)
	}

	results, err := db.GetResults(rtyq.Bounds(point))
	if err != nil {
		panic(err)
	}

	features := query.ResolveFeaturesFromPoint(point, results, data)

	fmt.Printf("%s\n", query.FeaturesToResponse(features))

}
