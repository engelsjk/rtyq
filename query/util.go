package query

import (
	"bytes"
	"fmt"

	"github.com/paulmach/orb/geojson"
)

// FeaturesToString converts a slice of geojson.Features to a comma-separated
// array string
func FeaturesToString(features [][]byte) string {

	if features == nil {
		return "[]"
	}

	out := fmt.Sprintf("[%s]", string(bytes.Join(features, []byte(","))))
	return out
}

func appendFeature(bs [][]byte, f *geojson.Feature) ([][]byte, error) {
	b, err := f.MarshalJSON()
	if err != nil {
		return bs, err
	}
	bs = append(bs, b)
	return bs, nil
}
