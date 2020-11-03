package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/engelsjk/rtyq/pkg/data"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/tidwall/buntdb"
)

// HandleData ...
func HandleData(w http.ResponseWriter, r *http.Request, bdb *buntdb.DB, dirData, index string) {

	// todo: clean up errors written to http

	lonlat := r.URL.Query().Get("pt")

	if lonlat == "" {
		http.Error(w, "please provide a lon,lat point", 400)
		return
	}

	pt := data.ParseLonLat(lonlat)

	results, err := db.Get(bdb, index, data.Bounds(pt))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	features, err := data.ResolveResults(dirData, results, lonlat)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	featuresStr := []string{}
	for _, f := range features {
		b, err := f.MarshalJSON()
		if err != nil {
			continue
		}
		featuresStr = append(featuresStr, string(b))
	}
	out := fmt.Sprintf("[%s]", strings.Join(featuresStr, ","))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(out))
}
