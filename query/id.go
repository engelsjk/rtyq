package query

import (
	"github.com/engelsjk/rtyq"
	"github.com/paulmach/orb/geojson"
)

// GetFeaturesFromID parses an ID string 'id',
// queries the database for results and returns
// the results as a slice of *geojson.Feature
func GetFeaturesFromID(id string, data rtyq.Data) ([]*geojson.Feature, error) {

	id, err := ParseID(id)
	if err != nil {
		return nil, err
	}

	features := ResolveFeaturesFromID(id, data)

	return features, nil
}

// ParseID converts an ID string 'id'
func ParseID(id string) (string, error) {

	// todo: add id validation

	return id, nil
}

// ResolveFeaturesFromID loads GeoJSON data from the data directory
// by the requested ID and returns a slice of *geojson.Feature
func ResolveFeaturesFromID(id string, data rtyq.Data) []*geojson.Feature {

	features := []*geojson.Feature{}

	fp := rtyq.FilePath(data.DirPath, id, data.FileExtension)

	f, err := rtyq.LoadFeature(fp)
	if err != nil {
		return features
	}

	features = appendFeature(features, f)

	return features
}
