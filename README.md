<p align="center">
    <img src="/.github/images/logo.png" 
    width="175" border="0" alt="rtyq"></a>
</p>

Rtyq is a command-line tool used to create spatially indexed databases of polygon data and to provide an API for spatial queries. It creates persisent database files on disk and generates in-memory R-tree spatial indexes to do (very) fast queries. Rtyq supports queries by point, tile or feature ID and serves GeoJSON data via a REST API. 

## Install

```
go get github.com/engelsjk/rtyq
```

## Run

```
rtyq check --config="config.json"
rtyq create --config="config.json"
rtyq start --config="config.json"
```

## Dependencies

* [tidwall/buntdb](https://github.com/tidwall/buntdb)
* [paulmach/orb](https://github.com/paulmach/orb)
* [karrick/godirwalk](https://github.com/karrick/godirwalk)
* [schollz/progressbar](https://github.com/schollz/progressbar)
* [go-chi/chi](https://github.com/go-chi/chi)
* [alecthomas/kingpin.v2](https://github.com/alecthomas/kingpin)
