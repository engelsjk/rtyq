package data

import (
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/tidwall/buntdb"
)

func dbPointBounds(p orb.Point) string {
	// bounds = [lon lat]

	var sb strings.Builder

	lon := strconv.FormatFloat(p.Lon(), 'f', -1, 64)
	lat := strconv.FormatFloat(p.Lat(), 'f', -1, 64)

	sb.WriteString("[")
	sb.WriteString(lon)
	sb.WriteString(" ")
	sb.WriteString(lat)
	sb.WriteString("]")

	bounds := sb.String()

	return bounds
}

func dbPolyBounds(b orb.Bound) string {
	// bounds = [minLon minLat], [maxLon maxLat]

	var sb strings.Builder

	minLon := strconv.FormatFloat(b.Min.Lon(), 'f', -1, 64)
	minLat := strconv.FormatFloat(b.Min.Lat(), 'f', -1, 64)
	maxLon := strconv.FormatFloat(b.Max.Lon(), 'f', -1, 64)
	maxLat := strconv.FormatFloat(b.Max.Lat(), 'f', -1, 64)

	sb.WriteString("[")
	sb.WriteString(minLon)
	sb.WriteString(" ")
	sb.WriteString(minLat)
	sb.WriteString("], [")
	sb.WriteString(maxLon)
	sb.WriteString(" ")
	sb.WriteString(maxLat)
	sb.WriteString("]")

	bounds := sb.String()

	return bounds
}

func dbUpdate(db *buntdb.DB, index, id, bounds string) error {
	return db.Update(func(tx *buntdb.Tx) error {
		k := dbKey(index, id)
		v := bounds
		tx.Set(k, v, nil)
		return nil
	})
}

func dbPattern(index string) string {
	var sb strings.Builder
	sb.WriteString(index)
	sb.WriteString(":*")
	pattern := sb.String()
	return pattern
}

func dbKey(index, id string) string {
	var sb strings.Builder
	sb.WriteString(index)
	sb.WriteString(":")
	sb.WriteString(id)
	key := sb.String()
	return key
}

func dbParseKey(key string) (string, string) {
	k := strings.Split(key, ":")
	return k[0], k[1]
}
