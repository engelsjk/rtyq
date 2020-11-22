package server

import (
	"fmt"
	"net/http"

	"github.com/engelsjk/rtyq/data"
	"github.com/go-chi/chi"
)

const (
	routeVarLayer = "layer"
	routeVarPoint = "point"
	routeVarTileX = "x"
	routeVarTileY = "y"
	routeVarTileZ = "z"
	routeVarID    = "id"
)

func initRouter() *chi.Mux {
	router := chi.NewRouter()
	return router
}

func addRoutes(router *chi.Mux) {
	addRoute(router, "/", handleRoot)
	addRoute(router, "/{layer}", handleLayer)
	addRoute(router, "/{layer}/point/{point}", handlePoint)
	addRoute(router, "/{layer}/tile/{z}/{x}/{y}", handleTile)
	addRoute(router, "/{layer}/id/{id}", handleID)
}

func addRoute(router *chi.Mux, path string, handler func(http.ResponseWriter, *http.Request) *serverError) {
	router.Handle(path, serverHandler(handler))
}

////////////////////////////////////////////////////////////////////////

func handleRoot(w http.ResponseWriter, r *http.Request) *serverError {
	return nil
}

func handleLayer(w http.ResponseWriter, r *http.Request) *serverError {
	layer := getRequestVar(routeVarLayer, r)
	_ = layer
	return nil
	// return writeJSON(w, ContentTypeJSON, fs)
}

func handlePoint(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	point := getRequestVar(routeVarPoint, r)

	features, err := data.QueryHandler.Point(layer, point)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func handleTile(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	tileX := getRequestVar(routeVarTileX, r)
	tileY := getRequestVar(routeVarTileY, r)
	tileZ := getRequestVar(routeVarTileZ, r)

	features, err := data.QueryHandler.Tile(layer, tileX, tileY, tileZ)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func handleID(w http.ResponseWriter, r *http.Request) *serverError {

	layer := getRequestVar(routeVarLayer, r)
	id := getRequestVar(routeVarID, r)

	features, err := data.QueryHandler.ID(layer, id)
	if err != nil {
		return errorQueryToServer(err)
	}

	return writeJSON(w, ContentTypeJSON, features)
}

func errorQueryToServer(err error) *serverError {
	switch err {
	case data.ErrQueryNoLayer:
		return serverErrorNotFound(err, err.Error())
	case data.ErrQueryInvalidPoint:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryInvalidPoint:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryExceededTileZoomLimit:
		return serverErrorBadRequest(err, err.Error())
	case data.ErrQueryRequest:
		return serverErrorInternal(err, err.Error())
	default:
		err := fmt.Errorf("unknown error")
		return serverErrorInternal(err, err.Error())
	}
}
