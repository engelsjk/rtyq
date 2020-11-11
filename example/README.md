## Example

For this example, let's use US state boundaries. We'll use data from Census TIGER/Line.

```
https://www.census.gov/cgi-bin/geo/shapefiles/index.php
```

Select a year and then select the layer 'States (and equivalent)', which will download a zip file.

```
tl_2019_us_state.zip (8.9 MB)
```

Extracting this .zip will create a folder containing all of the associated Shapefile files. But we want this data in GeoJSON.
And specifically, we want to work with the newline-delimited GeoJSON format, since it allows us to read data line-by-line. Luckily, ogr2ogr has a [GeoJSONSeq](https://gdal.org/drivers/vector/geojsonseq.html) file format which can be used to convert the downloaded Shapefile into newline-delimited GeoJSON.

```
ogr2ogr -f GeoJSONSeq -t_srs crs:84 tl_2019_us_state.ndgeojson tl_2019_us_state/tl_2019_us_state.shp
```

This will create a newline-delimited GeoJSON, which we can label with a ```.ndgeojson``` file extension to distinguish it from a normal GeoJSON file.

```
tl_2019_us_state.ndgeojson (24.4 MB)
```

Next, we need to split each of the GeoJSON features into their own separate files, using a uniqe ID for the feature's filename. In this case, we'll be using a state's FIPS code, given in feature properties as 'GEOID'. A Bash script will do this easily enough.

```
#!/bin/bash

IN="tl_2019_us_state.ndgeojson"

while read p; do
  ID=$(echo $p | jq -r '.properties.GEOID')
  OUT=$ID.geojson
  echo $p > $OUT
done < $IN
```

We could also use a custom CLI, ```gjsplit```, in the [engelsjk/gjfunks](https://github.com/engelsjk/gjfunks) package that I've built for this exact task.

```gjsplit --outkey="GEOID" --nd tl_2019_us_state.ndgeojson```

Either way, we should end up with 56 files (24.4 MB total) containing (1) GeoJSON Feature in each file. Using ```gjsplit``` should save some disk space (only 20.9 MB total) since the JSON encoding removes a good bit of white space.
