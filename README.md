# Rtyq

Rtyq is a command-line tool used to create spatially indexed databases of polygon data and to provide an API for spatial queries. It creates persisent database files on disk and generates in-memory R-tree spatial indexes to do (very) fast queries. Rtyq supports queries by point, tile or feature ID and serves GeoJSON data via a REST API. 

## Hybrid Data Model

Rtyq uses a hybrid data model for querying geometry data in order to resolve feature bounding box overlap issues. Database filesizes are minimized by only storing feature IDs. Fast, R-tree spatial queries (i.e. Intersect) return a best-guess set of features ID's which are used to load geometry from files on disk. Slower, last-mile spatial checks on feature geometries are used to resolve overlap issues, namely point-in-polygon checks for point queries and high-zoom level tile overlap checks for tile queries.

## Features

* Generates an R-tree spatially indexed database by reading geometry files on disk
* Fast directory traversel using [karrick/godirwalk](https://github.com/karrick/godirwalk)
* Supports Polygon and MultiPolygon GeoJSON Features 
* In-memory spatial index and database that persists on disk using [tidwall/buntdb](https://github.com/tidwall/buntdb)
* Hybrid model (Feature IDs in database / Feature geometries on disk) balances database write time and file size
* REST API allows queries by point, tile or ID using [go-chi/chi](https://github.com/go-chi/chi)
* API response as JSON array of GeoJSON Features
* CLI flags for single data layer or config.json file for multiple data layers

## Install

```
go get -u github.com/engelsjk/rtyq/cmd/rtyq
```

## Create

```
rtyq create \
--name="blocks"
--db="blocks.db" \
--index="block" \
--data="/path/to/blocks" \
--ext=".geojson" \
--id="GEOID10" \
```

## Start

```
rtyq start \
--name="blocks"
--db="blocks.db" \
--index="block" \
--data="/path/to/blocks" \
--ext=".geojson" \
--id="GEOID10" \
```

## Config File

```
{
    "port": 5500,
    "enable_logs": false,
    "throttle_limit": 1000,
    "layers": [
        {
            "name": "states",
            "data": {
                "path": "/path/to/states",
                "extension": ".geojson",
                "id": "GEOID"
            },
            "database": {
                "path": "states.db",
                "index": "state"
            },
            "service": {
                "endpoint": "states",
                "zoom_limit": 6
            }
        },
        {
            "name": "counties",
            "data": {
                "path": "/path/to/counties",
                "extension": ".geojson",
                "id": "GEOID"
            },
            "database": {
                "path": "counties.db",
                "index": "county"
            },
            "service": {
                "endpoint": "counties",
                "zoom_limit": 6
            }
        }
    ]
}
```

```
rtyq create --config="config.json"
rtyq start --config="config.json"
```

## Performance

## Dependencies

* [tidwall/buntdb](https://github.com/tidwall/buntdb)
* [paulmach/orb](https://github.com/paulmach/orb)
* [karrick/godirwalk](https://github.com/karrick/godirwalk)
* [schollz/progressbar](https://github.com/schollz/progressbar)
* [go-chi/chi](https://github.com/go-chi/chi)
* [alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)
