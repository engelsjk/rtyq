package api

import (
	"fmt"
	"net/http"

	"github.com/engelsjk/rtyq"
	"github.com/engelsjk/rtyq/query"
	"github.com/go-chi/chi"
)

// Handler is a structure that contains pointers
// to Data and Database structures, and includes a
// zoom limit integer
type Handler struct {
	Data      rtyq.Data
	Database  rtyq.DB
	ZoomLimit int
}

var (
	ErrUnknown               error = fmt.Errorf(`{"error": "unknown"}`)
	ErrBadQueryParam         error = fmt.Errorf(`{"error": "invalid query parameter"}`)
	ErrUnknownQueryType      error = fmt.Errorf(`{"error": "unknown query type"}`)
	ErrInvalidPoint          error = fmt.Errorf(`{"error": "invalid lon,lat point"}`)
	ErrInvalidTile           error = fmt.Errorf(`{"error": "invalid z/x/y tile"}`)
	ErrTileZoomLimitExceeded error = fmt.Errorf(`{"error": "tile request exceeded zoom limit"}`)
	ErrUnableToGetDataFromDB error = fmt.Errorf(`{"error": "unable to get data from db"}`)
)

// HandleLayer parses an API query by type and runs a response function
// to write the query response
func (h Handler) HandleLayer(w http.ResponseWriter, r *http.Request, queryType string) {

	switch queryType {
	case "point":
		responsePoint(w, r, h)
	case "tile":
		responseTile(w, r, h)
	case "id":
		responseID(w, r, h)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrUnknownQueryType.Error()))
		return
	}
}

func responsePoint(w http.ResponseWriter, r *http.Request, h Handler) {

	point := chi.URLParam(r, "point")

	features, err := query.GetFeaturesFromPoint(point, h.Database, h.Data) // this needs to return bytes

	if err != nil {
		responseError(w, err)
		return
	}

	response := query.FeaturesToResponse(features)

	w.Write(response)
	return
}

func responseTile(w http.ResponseWriter, r *http.Request, h Handler) {

	z := chi.URLParam(r, "z")
	x := chi.URLParam(r, "x")
	y := chi.URLParam(r, "y")

	features, err := query.GetFeaturesFromTile(z, x, y, h.ZoomLimit, h.Database, h.Data)

	if err != nil {
		responseError(w, err)
		return
	}

	response := query.FeaturesToResponse(features)

	w.Write([]byte(response))
	return
}

func responseID(w http.ResponseWriter, r *http.Request, h Handler) {

	id := chi.URLParam(r, "id")

	features, err := query.GetFeaturesFromID(id, h.Data)

	if err != nil {
		responseError(w, err)
		return
	}

	response := query.FeaturesToResponse(features)

	w.Write([]byte(response))
	return
}

func responseError(w http.ResponseWriter, err error) {

	var statusCode int
	var response []byte

	switch err {
	case query.ErrInvalidTile:
		statusCode = http.StatusBadRequest
		response = []byte(ErrInvalidTile.Error())
	case query.ErrTileZoomLimitExceeded:
		statusCode = http.StatusBadRequest
		response = []byte(ErrTileZoomLimitExceeded.Error())
	case rtyq.ErrDatabaseFailedToGetResults:
		statusCode = http.StatusInternalServerError
		response = []byte(ErrUnableToGetDataFromDB.Error())
	default:
		statusCode = http.StatusInternalServerError
		response = []byte(ErrUnknown.Error())
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(response))
}
