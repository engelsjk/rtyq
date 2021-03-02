package data

import (
	"strconv"
	"strings"

	"github.com/tidwall/buntdb"
)

func dbPointBounds(lon, lat float64) string {
	// bounds = [lon lat]

	var sb strings.Builder

	sb.WriteString("[")
	sb.WriteString(strconv.FormatFloat(lon, 'f', -1, 64))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatFloat(lat, 'f', -1, 64))
	sb.WriteString("]")

	bounds := sb.String()

	return bounds
}

func dbPolyBounds(minLon, minLat, maxLon, maxLat float64) string {
	// bounds = [minLon minLat], [maxLon maxLat]

	var sb strings.Builder

	sb.WriteString("[")
	sb.WriteString(strconv.FormatFloat(minLon, 'f', -1, 64))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatFloat(minLat, 'f', -1, 64))
	sb.WriteString("], [")
	sb.WriteString(strconv.FormatFloat(maxLon, 'f', -1, 64))
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatFloat(maxLat, 'f', -1, 64))
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
