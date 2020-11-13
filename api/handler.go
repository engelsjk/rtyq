package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/engelsjk/rtyq"
	"github.com/engelsjk/rtyq/query"
	"github.com/go-chi/chi"
)

// Handler is a structure that contains pointers
// to Data and Database structures, and includes a
// zoom limit integer
type Handler struct {
	Data      *rtyq.Data
	Database  *rtyq.DB
	ZoomLimit int
}

var (
	ErrUnknown               error = fmt.Errorf(`{"error": "unknown"}`)
	ErrBadQueryParam         error = fmt.Errorf(`{"error": "please provide a valid query parameter (pt, tile or id)"}`)
	ErrUnknownQueryType      error = fmt.Errorf(`{"error": "unknown query type"}`)
	ErrInvalidPoint          error = fmt.Errorf(`{"error": "please provide a valid lon,lat point"}`)
	ErrInvalidTile           error = fmt.Errorf(`{"error": "please provide a valid z/x/y tile"}`)
	ErrTileZoomLimitExceeded error = fmt.Errorf(`{"error": "tile request zoom limit exceeded"}`)
	ErrUnableToGetDataFromDB error = fmt.Errorf(`{"error": "unable to get data from db"}`)
)

// HandleLayer parses an API query by type and runs a response function
// to write the query response
func (h *Handler) HandleLayer(w http.ResponseWriter, r *http.Request, queryType string, enableLogs bool) {

	if enableLogs {
		url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
		log.Printf(url)
	}

	switch queryType {
	case "point":
		responsePoint(w, r, h)
	case "tile":
		responseTile(w, r, h)
	case "id":
		responseID(w, r, h)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrUnknownQueryType.Error()))
		return
	}
}

func responsePoint(w http.ResponseWriter, r *http.Request, h *Handler) {

	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var response string

	point := chi.URLParam(r, "point")

	features, err := query.GetFeaturesFromPoint(point, h.Database, h.Data)

	if err != nil {

		switch err {
		case query.ErrInvalidPoint:
			statusCode = http.StatusBadRequest
			response = ErrInvalidPoint.Error()
		case rtyq.ErrDatabaseFailedToGetResults:
			statusCode = http.StatusInternalServerError
			response = ErrUnableToGetDataFromDB.Error()
		default:
			statusCode = http.StatusInternalServerError
			response = ErrUnknown.Error()
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(response))
		return
	}

	response = rtyq.FeaturesToString(features)

	w.Write([]byte(response))
	return
}

func responseTile(w http.ResponseWriter, r *http.Request, h *Handler) {

	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var response string

	z := chi.URLParam(r, "z")
	x := chi.URLParam(r, "x")
	y := chi.URLParam(r, "y")

	tile := fmt.Sprintf("%s/%s/%s", z, x, y)

	features, err := query.GetFeaturesFromTile(tile, h.ZoomLimit, h.Database, h.Data)

	if err != nil {

		switch err {
		case query.ErrInvalidTile:
			statusCode = http.StatusBadRequest
			response = ErrInvalidTile.Error()
		case query.ErrTileZoomLimitExceeded:
			statusCode = http.StatusBadRequest
			response = ErrTileZoomLimitExceeded.Error()
		case rtyq.ErrDatabaseFailedToGetResults:
			statusCode = http.StatusInternalServerError
			response = ErrUnableToGetDataFromDB.Error()
		default:
			statusCode = http.StatusInternalServerError
			response = ErrUnknown.Error()
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(response))
		return
	}

	response = rtyq.FeaturesToString(features)

	w.Write([]byte(response))
	return
}

func responseID(w http.ResponseWriter, r *http.Request, h *Handler) {

	w.Header().Set("Content-Type", "application/json")

	var statusCode int
	var response string

	id := chi.URLParam(r, "id")

	features, err := query.GetFeaturesFromID(id, h.Data)

	if err != nil {

		switch err {
		case query.ErrInvalidTile:
			statusCode = http.StatusBadRequest
			response = ErrInvalidTile.Error()
		case query.ErrTileZoomLimitExceeded:
			statusCode = http.StatusBadRequest
			response = ErrTileZoomLimitExceeded.Error()
		case rtyq.ErrDatabaseFailedToGetResults:
			statusCode = http.StatusInternalServerError
			response = ErrUnableToGetDataFromDB.Error()
		default:
			statusCode = http.StatusInternalServerError
			response = ErrUnknown.Error()
		}

		w.WriteHeader(statusCode)
		w.Write([]byte(response))
		return
	}

	response = rtyq.FeaturesToString(features)

	w.Write([]byte(response))
	return
}
