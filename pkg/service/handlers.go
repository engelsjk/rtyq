package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/engelsjk/rtyq/pkg/data"
	"github.com/engelsjk/rtyq/pkg/db"
	"github.com/paulmach/orb/geojson"
	"github.com/tidwall/buntdb"
)

type Handler struct {
	DirData   string
	Extension string
	Database  *buntdb.DB
	Index     string
	Error     error
}

var (
	ErrBadQueryParam  error = fmt.Errorf(`{"error": "please provide a valid query parameter (pt, tile or id)"}`)
	ErrInvalidPoint   error = fmt.Errorf(`{"error": "please provide a valid lon,lat point"}`)
	ErrInvalidTile    error = fmt.Errorf(`{"error": "please provide a valid z/x/y tile"}`)
	ErrDatabaseGet    error = fmt.Errorf(`{"error": "unable to get data from db"}`)
	ErrResolveRequest error = fmt.Errorf(`{"error": "unable to resolve request"}`)
)

// HandleData ...
func (h *Handler) HandleData(w http.ResponseWriter, r *http.Request) {

	h.Error = nil

	w.Header().Set("Content-Type", "application/json")

	pt := r.URL.Query().Get("pt")
	tile := r.URL.Query().Get("tile")
	id := r.URL.Query().Get("id")

	if pt == "" && tile == "" && id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrBadQueryParam.Error()))
		return
	}

	var features []*geojson.Feature

	fmt.Printf("%v\n", id)

	if pt != "" {
		features = h.getFeaturesFromPoint(pt)
	} else if tile != "" {
		features = h.getFeaturesFromTile(tile)
	} else if id != "" {
		features = h.getFeaturesFromID(id)
	}

	out := h.output(features)

	switch h.Error {
	case ErrInvalidPoint:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrInvalidPoint.Error()))
	case ErrInvalidTile:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrInvalidTile.Error()))
	case ErrDatabaseGet:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrDatabaseGet.Error()))
	case ErrResolveRequest:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrResolveRequest.Error()))
	default:
		w.Write([]byte(out))
	}
}

func (h *Handler) output(features []*geojson.Feature) string {

	if h.Error != nil {
		return ""
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
	return out
}

func (h *Handler) getFeaturesFromPoint(pt string) []*geojson.Feature {

	if h.Error != nil {
		return nil
	}

	point, err := data.ParseLonLatPoint(pt)
	if err != nil {
		h.Error = ErrInvalidPoint
		return nil
	}

	results, err := db.Get(h.Database, h.Index, data.Bounds(point))
	if err != nil {
		h.Error = ErrDatabaseGet
		return nil
	}

	features, err := data.ResolvePoint(h.DirData, h.Extension, results, point)
	if err != nil {
		h.Error = ErrResolveRequest
		return nil
	}

	return features
}

func (h *Handler) getFeaturesFromTile(t string) []*geojson.Feature {
	if h.Error != nil {
		return nil
	}

	tile, err := data.ParseTile(t)
	if err != nil {
		h.Error = ErrInvalidTile
		return nil
	}

	results, err := db.Get(h.Database, h.Index, data.Bounds(tile))
	if err != nil {
		h.Error = ErrDatabaseGet
		return nil
	}

	features, err := data.ResolveTile(h.DirData, h.Extension, results, tile)
	if err != nil {
		h.Error = ErrResolveRequest
		return nil
	}

	return features
}

func (h *Handler) getFeaturesFromID(id string) []*geojson.Feature {
	features, err := data.ResolveID(h.DirData, h.Extension, id)
	if err != nil {
		h.Error = ErrResolveRequest
		return nil
	}

	return features
}
