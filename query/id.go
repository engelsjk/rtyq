package query

import (
	"github.com/engelsjk/rtyq"
	"github.com/paulmach/orb/geojson"
)

// GetFeaturesFromID ...
func GetFeaturesFromID(id string, data *rtyq.Data) ([]*geojson.Feature, error) {

	id, err := ParseID(id)
	if err != nil {
		return nil, err
	}

	features := ResolveFeaturesFromID(id, data)

	return features, nil
}

// ParseID ...
func ParseID(id string) (string, error) {

	// todo: add id validation

	return id, nil
}

// ResolveFeaturesFromID ...
func ResolveFeaturesFromID(id string, data *rtyq.Data) []*geojson.Feature {

	features := []*geojson.Feature{}

	fp := rtyq.FilePath(data.DirPath, id, data.FileExtension)

	f, err := rtyq.LoadFeature(fp)
	if err != nil {
		return features
	}

	features = append(features, f)
	return features
}
