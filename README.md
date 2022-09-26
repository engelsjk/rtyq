<p align="center">
    <img src="/.github/images/logo.png"
    width="175" border="0" alt="rtyq"></a>
</p>

Rtyq is a command-line tool used to create spatially indexed databases of polygon data and to provide an API for spatial queries. It creates persistent database files on disk and generates in-memory R-tree spatial indexes to do (very) fast queries. Rtyq supports queries by point, tile or feature ID and serves GeoJSON data via a REST API.

## Install

```
go install github.com/engelsjk/rtyq
```

Or you can use an existing binary in ```bin```.

## Tool

Rtyq has three actions: check, create, start.

### Create

```create``` converts a directory of GeoJSON Feature files into a static database file with a spatial index. Each feature must have a unique property ID. A database file will be created fOr each layer specified in the configuration file.

### Start

```start``` loads each layer's database file into memory and start a web server that listens for spatial queries.

## Configuration

A ```config.json``` is required. It must have two fields: ```server``` and ```layers```.

```json
{
    "layers": [...],
    "server": {.},
}
```

### Layers

Rtyq can support an arbitrary number of data layers, where each layer is a distinct set of spatial data. Each layer will have its own API route to receive spatial queries.

In the configuration files, these layers are specified by an array of objects that specify the directory, extension and unique ID property name for the input data set, as well as the path of the output static database file.

```json
{
    "name": "states",
    "data": {
        "dir": ".../data/states",
        "ext": ".geojson",
        "id": "GEOID"
    },
    "database": {
        "filepath": ".../db/states.db",
        "index": "state"
    },
    "service": {
        "zoomlimit": 6
    }
}
```

### Server

Server options in the configuration file include a port number and other settings.

```bash
{
    "layers": [...],
    "server": {
        "port": 5500,
        "logs": true,
        "debug": false,
        "throttle": 1000,
        "readtimeoutsec": 21
    },
}
```

## Run

With a ```config.json``` in your working directory, run the command ```rtyq check``` to get information on your specified data directories.

Then run ```rtyq create``` to create a database file for each layer.

Finally, run ```rtyq start``` to load the database files into memory and start the web server. Loading the databases into memory can take some time depending on a few factors.

## Queries

The web server provides the following queries for each layer:

Bounding box: ```/{layer}/bbox/{bbox}``` where bbox is of the form {minX,minY,maxX,maxY}

Tile: ```/{layer}/tile/{z}/{x}/{y}```

ID: ```/{layer}/id/{id}```

## Dependencies

* [tidwall/buntdb](https://github.com/tidwall/buntdb)
* [paulmach/orb](https://github.com/paulmach/orb)
* [karrick/godirwalk](https://github.com/karrick/godirwalk)
* [schollz/progressbar](https://github.com/schollz/progressbar)
* [go-chi/chi](https://github.com/go-chi/chi)
* [alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)
